package webserver

import (
	"log"
	"net"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// One client per IP restriction
var oneClientPerIP bool = false

// Keep track of trusted client IPs in a set (empty-valued map)
var trustedClientIPs = make(map[string](struct{}))

// Set the one client per IP restriction based on a configuration
func ConfigOneClientPerIP(_oneClientPerIP bool) {
	oneClientPerIP = _oneClientPerIP
}

// Set the one client per IP restriction based on a configuration
func ConfigTrustedClientIPs(_trustedClientIPs []string) {
	for _, ip := range _trustedClientIPs {
		trustedClientIPs[ip] = struct{}{}
	}
}

// Store the responses from trusted clients in a (send-only) channel
var responseCh chan<- []byte

/*
Map to keep track of websocket client IPs; if only
one client connection is allowed per IP, kick the oldest
*/
var ipSessionMap = make(map[string]*webSession)
var muISM sync.Mutex

/*
Get the IP address by taking the part of the remote address
{ip}:{port} before the last colon separating the address and port
*/
func getIP(conn *websocket.Conn) string {
	addr := conn.RemoteAddr().String()
	sepIdx := strings.LastIndex(addr, ":")
	return addr[:sepIdx]
}

// Web session object, for keeping track of individual websocket sessions
type webSession struct {
	sendCh chan []byte
	readEn bool // read enabled (allowed by IP whitelist)
	conn   *websocket.Conn
	sync.Mutex
}

// Create a new web session object
func newWebSession(conn *websocket.Conn) *webSession {
	return &webSession{
		sendCh: make(chan []byte, 10),
		readEn: true,
		conn:   conn,
	}
}

// Register this web session in the active connections
func (ws *webSession) register() {
	muISM.Lock()
	// If we've seen this IP address before, kick the old one and start a new one
	ip := getIP(ws.conn)
	if oldSession, ok := ipSessionMap[ip]; ok && oneClientPerIP {
		oldSession.quit()
	}

	// Connect this quit channel to the IP address
	ipSessionMap[ip] = ws
	muISM.Unlock()

	/*
		Determine if we trust this new connection, by checking against configured
		trusted connections
	*/
	_, trusted := trustedClientIPs[ip]
	if !trusted {
		ws.readEn = false
	}

	// Lock the mutex so we can keep track of the number of open clients
	muOWS.Lock()
	{
		// Add this web session to the web sessions set
		openWebSessions[ws] = struct{}{}
		if trusted {
			log.Printf("\033[34m[%d -> %d] trusted client connected (%s)\033[0m\n",
				len(openWebSessions)-1, len(openWebSessions), ip)
		} else {
			log.Printf("\033[34m[%d -> %d] client connected (%s)\033[0m\n",
				len(openWebSessions)-1, len(openWebSessions), ip)
		}
	}
	muOWS.Unlock()
}

// Unregister this web session in the active connections
func (ws *webSession) unregister() {
	// Record the IP address of the disconnecting client
	ip := getIP(ws.conn)
	_, trusted := trustedClientIPs[ip]

	/*
		Lock the mutex so that other channels will not read the open web
		sessions map until this is complete
	*/
	muOWS.Lock()
	{
		// Print information regarding the disconnect
		if len(openWebSessions) > 0 {
			if trusted {
				log.Printf("\033[33m[%d -> %d] trusted client disconnected (%s)\033[0m\n",
					len(openWebSessions), len(openWebSessions)-1, ip)
			} else {
				log.Printf("\033[33m[%d -> %d] client disconnected (%s)\033[0m\n",
					len(openWebSessions), len(openWebSessions)-1, ip)
			}
		} else {
			log.Printf("\033[33m[X -> X] client(s) blocked\033[0m\n")
		}

		// Remove this websession from the open web sessions set
		delete(openWebSessions, ws)
	}
	muOWS.Unlock()

	// We aren't active anymore, don't need to remember us in IP session map
	muISM.Lock()
	if ipSessionMap[ip] == ws {
		delete(ipSessionMap, ip)
	}
	muISM.Unlock()

	// Wait until deregister to prevent write to closed channel
	close(ws.sendCh)
}

// Close the websocket client (causes loop to unblock)
func (ws *webSession) quit() {
	ws.conn.Close()
	// Wake the send loop, if it needs to be reminded to exit
	// Any message will cause readLoop to exit as the socket is closed
	select {
	case ws.sendCh <- nil:
	default:
	}
}

// Runs all loops to service the connection and blocks until complete
func (ws *webSession) loop() {
	if !ws.readEn {
		defer ws.quit()
		ws.sendLoop()
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer ws.quit()
		defer wg.Done()
		ws.readLoop()
	}()
	go func() {
		defer ws.quit()
		defer wg.Done()
		ws.sendLoop()
	}()
	wg.Wait()
}

/*
Run a read-loop go-routine to keep track of the incoming messages
for a given session
*/
func (ws *webSession) readLoop() {
	// "While" loop, keep reading until the connection closes
	for {
		// Read a message (discard the type since we don't need it)
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {

			// Types of errors which we intentionally catch and return from
			clientCloseErr := websocket.IsCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			)
			serverCloseErr := (err == websocket.ErrCloseSent)
			abnormalCloseErr := websocket.IsUnexpectedCloseError(err)
			_, netErr := err.(*net.OpError)
			if clientCloseErr || serverCloseErr || netErr || abnormalCloseErr {
				return
			}

			closeErr, ok := err.(*websocket.CloseError)
			if ok && closeErr.Code == websocket.CloseNormalClosure {
				return
			}

			// For all other unspecified errors, log them and quit
			log.Println("read error:", err)
			return
		}

		// Skip this message if it is empty
		if len(msg) == 0 {
			continue
		}

		responseCh <- msg
		if cap(responseCh) == len(responseCh) {
			log.Println("\033[35mWARN: Incoming messages " +
				"full, server not keeping up \033[0m")
		}
	}
}

// Sending websocket data (binary)
func (ws *webSession) sendLoop() {
	// "While" loop, keep sending until the connection closes
	for {

		// Block until the next message is ready
		msg := <-ws.sendCh

		// nil means we are told to exit
		if msg == nil {
			return
		}

		// Try writing the message
		if err := ws.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {

			// Types of errors which we intentionally catch and return from
			clientCloseErr := websocket.IsCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			)
			serverCloseErr := (err == websocket.ErrCloseSent)
			abnormalCloseErr := websocket.IsUnexpectedCloseError(err)
			_, netErr := err.(*net.OpError)
			if clientCloseErr || serverCloseErr || netErr || abnormalCloseErr {
				return
			}

			// For all other unspecified errors, log them and quit
			log.Println("write error:", err)
			return
		}
	}
}
