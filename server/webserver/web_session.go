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

// Keep track of active websocket sessions in a set (empty-valued map)
var openWebSocketSessions = make(map[*webSession](struct{}))

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

	// If we've seen this IP address before, kick the old one and start a new one
	ip := getIP(ws.conn)
	if oldQuitCh, ok := ipQuitMap[ip]; ok && oneBrowserPerIP {
		oldQuitCh <- struct{}{}
	}

	// Connect this quit channel to the IP address
	ipQuitMap[ip] = ws.quitCh

	// Add this to the open web sessions set
	muWS.Lock()
	{
		openWebSocketSessions[ws] = struct{}{}
		fmt.Printf("\033[34m[%d -> %d] browser connected\033[0m\n", len(openWebSocketSessions)-1, len(openWebSocketSessions))
	}
	muWS.Unlock()
}

// Unregister this web session in the active connections
func (ws *webSession) unregister() {

	// Close the connection and data channels
	ws.closed = true
	ws.conn.Close()
	close(ws.readCh)
	close(ws.sendCh)

	// Lock the mutex so we can keep track of the number of open clients
	muWS.Lock()
	{
		delete(openWebSocketSessions, ws)
		fmt.Printf("\033[33m[%d -> %d] browser disconnected\033[0m\n", len(openWebSocketSessions)+1, len(openWebSocketSessions))
	}
	muWS.Unlock()
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
			closeErr := websocket.IsCloseError(err, websocket.CloseGoingAway)
			serverCloseErr := (err == websocket.ErrCloseSent)
			_, netErr := err.(*net.OpError)
			if closeErr || serverCloseErr || netErr {
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
func (ws *webSession) writeLoop() {

	// If we ever stop sending messages (due to some error), kill the connection
	defer func() { ws.quitCh <- struct{}{} }()

	// Message to send
	msg := []byte("msg")

	for {

		// Block until the next message is ready
		<-ws.sendCh

		// Try writing the message
		if err := ws.conn.WriteMessage(websocket.TextMessage, msg); err != nil {

			// Types of errors which we intentially catch and return from
			closeErr := websocket.IsCloseError(err, websocket.CloseGoingAway)
			serverCloseErr := (err == websocket.ErrCloseSent)
			_, netErr := err.(*net.OpError)
			if closeErr || serverCloseErr || netErr {
				return
			}

			// For all other unspecified errors, log them and quit
			fmt.Println("write error: ", err)
			return
		}
	}
}
