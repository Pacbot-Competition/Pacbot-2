package webserver

import (
	"fmt"
	"net"
	"strings"

	"github.com/gorilla/websocket"
)

// One browser per IP restriction
var oneBrowserPerIP bool = false

// Set the one browser per IP restriction based on a configuration
func ConfigOneBrowserPerIP(OneBrowserPerIP bool) {
	oneBrowserPerIP = OneBrowserPerIP
}

/*
Map to keep track of websocket client IPs; if only
one browser connection is allowed per IP, kick the oldest
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
	readCh chan []byte
	sendCh chan []byte
	closed bool
	conn   *websocket.Conn
}

// Create a new web session object
func newWebSession(conn *websocket.Conn) *webSession {
	return &webSession{
		quitCh: make(chan struct{}, 2),
		readCh: make(chan []byte, 10),
		sendCh: make(chan []byte, 10),
		closed: false,
		conn:   conn,
	}
}

// Register this web session in the active connections
func (ws *webSession) register() {

	// If we're not allowing new connections, kick the connection
	if !newConnectionsAllowed {
		ws.quitCh <- struct{}{}
		return
	}

	// If we've seen this IP address before, kick the old one and start a new one
	ip := getIP(ws.conn)
	if oldQuitCh, ok := ipQuitMap[ip]; ok && oneBrowserPerIP {
		oldQuitCh <- struct{}{}
	}

	// Connect this quit channel to the IP address
	ipQuitMap[ip] = ws.quitCh

	// Lock the mutex so we can keep track of the number of open clients
	muOWS.Lock()
	{
		openWebSessions[ws] = struct{}{}
		fmt.Printf("\033[34m[%d -> %d] browser connected\033[0m\n", len(openWebSessions)-1, len(openWebSessions))
	}
	muOWS.Unlock()
}

// Unregister this web session in the active connections
func (ws *webSession) unregister() {

	// Close the connection and data channels
	ws.closed = true
	ws.conn.Close()
	close(ws.readCh)
	close(ws.sendCh)

	// Lock the mutex so that other channels will not read the open web sessions map until this is complete
	muOWS.Lock()
	{
		// Print information regarding the disconnect
		if newConnectionsAllowed || (len(openWebSessions) > 0) {
			fmt.Printf("\033[33m[%d -> %d] browser disconnected\033[0m\n", len(openWebSessions), len(openWebSessions)-1)
		} else {
			fmt.Printf("\033[33m[X -> X] browser(s) blocked\033[0m\n")
		}

		delete(openWebSessions, ws)
	}
	muOWS.Unlock()
}

// Run a read-loop go-routine to keep track of the incoming messages for a session
func (ws *webSession) readLoop() {

	// If we ever stop receiving messages (due to some error), kill the connection
	defer func() { ws.quitCh <- struct{}{} }()

	for {

		// Read a message (discard the type since we don't need it)
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {

			// Types of errors which we intentially catch and return from
			clientCloseErr := websocket.IsCloseError(err, websocket.CloseGoingAway)
			serverCloseErr := (err == websocket.ErrCloseSent)
			_, netErr := err.(*net.OpError)
			if clientCloseErr || serverCloseErr || netErr {
				return
			}

			// For all other unspecified errors, log them and quit
			fmt.Println("read error: ", err)
			return
		}

		// Save the message received into the read channel
		if !ws.closed {
			ws.readCh <- msg
		}

		// Print the message we received
		fmt.Printf("\033[2m\033[36m| Browser: %s`\033[0m\n", string(<-ws.readCh))
	}

}

// Sending websocket data (binary)
func (ws *webSession) sendLoop() {

	// If we ever stop sending messages (due to some error), kill the connection
	defer func() { ws.quitCh <- struct{}{} }()

	for {

		// Block until the next message is ready
		msg := <-ws.sendCh

		// Try writing the message
		if err := ws.conn.WriteMessage(websocket.TextMessage, msg); err != nil {

			// Types of errors which we intentially catch and return from
			clientCloseErr := websocket.IsCloseError(err, websocket.CloseGoingAway)
			serverCloseErr := (err == websocket.ErrCloseSent)
			_, netErr := err.(*net.OpError)
			if clientCloseErr || serverCloseErr || netErr {
				return
			}

			// For all other unspecified errors, log them and quit
			fmt.Println("write error: ", err)
			return
		}
	}
}
