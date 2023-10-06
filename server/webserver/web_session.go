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
var ipQuitMap = make(map[string](chan struct{}))

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
	quitCh chan struct{}
	sendCh chan []byte
	readEn bool // read enabled (allowed by IP whitelist)
	readOk bool // read ready (rate limiting)
	conn   *websocket.Conn
	sync.Mutex
}

// Create a new web session object
func newWebSession(conn *websocket.Conn) *webSession {
	return &webSession{
		quitCh: make(chan struct{}, 2),
		sendCh: make(chan []byte, 10),
		readEn: true,
		readOk: false,
		conn:   conn,
	}
}

// Register this web session in the active connections
func (ws *webSession) register() {

	// Increment the web broker quit wait group counter
	wgQuit.Add(1)

	// If we're not allowing new connections, kick the connection
	if !newConnectionsAllowed {
		ws.quitCh <- struct{}{}
		return
	}

	// If we've seen this IP address before, kick the old one and start a new one
	ip := getIP(ws.conn)
	if oldQuitCh, ok := ipQuitMap[ip]; ok && oneClientPerIP {
		oldQuitCh <- struct{}{}
	}

	/*
		Determine if we trust this new connection, by checking against configured
		trusted connections
	*/
	_, trusted := trustedClientIPs[ip]
	if !trusted {
		ws.readEn = false
	}

	// Connect this quit channel to the IP address
	ipQuitMap[ip] = ws.quitCh

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

	// Close the connection and data channels
	ws.readEn = false // Prevent this session from reading anymore
	ws.conn.Close()
	close(ws.sendCh)

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
		if newConnectionsAllowed || (len(openWebSessions) > 0) {
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

	// Decrement the web broker quit wait group counter
	wgQuit.Done()
}

/*
Run a read-loop go-routine to keep track of the incoming messages
for a given session
*/
func (ws *webSession) readLoop() {

	// If we ever stop receiving messages (due to some error), kill the connection
	defer func() { ws.quitCh <- struct{}{} }()

	// Local variable, to keep track of whether a message is discarded
	ignored := false

	// "While" loop, keep reading until the connection closes
	for {

		// Block indefinitely if we should not be reading from this client
		if !ws.readEn {
			select {}
		}

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
			_, netErr := err.(*net.OpError)
			if clientCloseErr || serverCloseErr || netErr {
				return
			}

			closeErr, ok := err.(*websocket.CloseError)
			if ok && closeErr.Code == websocket.CloseNormalClosure {
				log.Println("caught!")
			}

			// For all other unspecified errors, log them and quit
			log.Println("read error:", err)
			return
		}

		// Skip this message if it is empty or we should not read from this client
		if len(msg) == 0 {
			continue
		}

		// Relay the message immediately if it was a pause or play command
		if msg[0] == 'p' || msg[0] == 'P' {
			responseCh <- msg
			continue
		}

		// Determine whether this message should be ignored
		ws.Lock()
		{
			// Ignore the message if rate limiting applies
			ignored = !ws.readOk

			// Rate limiting - discard all incoming messages until we send a message
			ws.readOk = false
		}
		ws.Unlock()

		// If the message shouldn't be ignored, relay it, otherwise warn the user
		if !ignored {
			responseCh <- msg
		} else {
			log.Println("\033[35mWARN: An incoming message was ignored " +
				"(reason: rate limiting)\033[0m")
		}
	}
}

// Sending websocket data (binary)
func (ws *webSession) sendLoop() {

	// If we ever stop sending messages (due to some error), kill the connection
	defer func() { ws.quitCh <- struct{}{} }()

	// "While" loop, keep sending until the connection closes
	for {

		// Block until the next message is ready
		msg := <-ws.sendCh

		// Try writing the message
		if err := ws.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {

			// Types of errors which we intentionally catch and return from
			clientCloseErr := websocket.IsCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			)
			serverCloseErr := (err == websocket.ErrCloseSent)
			_, netErr := err.(*net.OpError)
			if clientCloseErr || serverCloseErr || netErr {
				return
			}

			// For all other unspecified errors, log them and quit
			log.Println("write error:", err)
			return
		}

		// Update the read ready state, to allow another read after this write
		ws.Lock()
		{
			ws.readOk = true
		}
		ws.Unlock()
	}
}
