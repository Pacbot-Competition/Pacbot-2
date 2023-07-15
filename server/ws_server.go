package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// Keep track of the number of open websocket connections
var openWebClients int = 0

// Make a map to keep track of web socket connections
// IP address (string) -> quit channel
var ipQuitMap = make(map[string](chan struct{}))

// Protect the above two values from race conditions with a
// mutex --> it will lock the resources until they are written to
var mu sync.Mutex

/*
Websockets are the way that the server will communicate with
the browser clients over the LAN. When connecting to a web-
socket, we connect normally over HTTP and then upgrade the
connection upon agreement between the server and client.
*/
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Receiving websocket data (binary)
func webSocketReceiveHandler(conn *websocket.Conn, webSocketQuitChan chan<- struct{}) {

	defer func() { webSocketQuitChan <- struct{}{} }()

	for {

		// Read a message (discard the type since we don't need it)
		_, msg, err := conn.ReadMessage()
		if err != nil {

			// If the browser gets disconnected, quit (which will close the connection)
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				return
			}

			// If the connection forcefully ends, quit
			if _, ok := err.(*net.OpError); ok {
				return
			}

			// For all other errors, log them
			fmt.Println("read error: ", err)
			return
		}

		// Print the message we received
		fmt.Printf("\033[2m\033[36m| Browser: %s`\033[0m\n", string(msg))
	}
}

// Sending websocket data (binary)
func webSocketSendHandler(conn *websocket.Conn, webSocketWriteChan <-chan []byte) {

	// Message to send
	msg := []byte(".")

	for {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {

			// If we closed the channel, quit
			if err == websocket.ErrCloseSent {
				return
			}

			// If the connection forcefully ends, quit
			if _, ok := err.(*net.OpError); ok {
				return
			}

			// For all other errors, log them
			fmt.Println("write error: ", err)
			return
		}

		// Block until the next message is ready
		<-webSocketWriteChan
	}
}

/*
Get the IP address by taking the part of the remote address
(such as 192.168.1.1:3000) before the last colon --> we use
this to limit the number of web sockets on an IP address to 1
*/
func getIP(conn *websocket.Conn) string {
	addr := conn.RemoteAddr().String()
	sepIdx := strings.LastIndex(addr, ":")
	return addr[:sepIdx]
}

/*
This handler makes sure that once we connect to the websocket,
all communication goes smoothly.
*/
func webSocketHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrades the connection, and quits if it didn't work out.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// If we've seen this IP address before, kick the old one and start a new one
	ip := getIP(conn)
	if oldQuitChannel, ok := ipQuitMap[ip]; ok {
		mu.Lock()
		oldQuitChannel <- struct{}{}
		mu.Unlock()
	}

	// Print a status message
	mu.Lock()
	openWebClients++
	fmt.Printf("\033[34m[%d -> %d] browser connected\033[0m\n", openWebClients-1, openWebClients)
	mu.Unlock()

	// Quit channel to postpone closing the websocket -> once this added to, we quit immediately
	webSocketQuitChan := make(chan struct{}, 1)

	// Write channel to postpone writing to the websocket
	webSocketWriteChan := make(chan []byte, 10)

	// Close the connection at the end of the function, or if
	// something goes wrong.
	defer func() {
		conn.Close()
		close(webSocketWriteChan)

		// Lock the mutex so we can keep track of the number of open clients
		mu.Lock()
		openWebClients--
		fmt.Printf("\033[33m[%d -> %d] browser disconnected\033[0m\n", openWebClients+1, openWebClients)
		mu.Unlock()
	}()

	// Store this as the latest quit channel for this IP
	ipQuitMap[ip] = webSocketQuitChan

	// Goroutine: handles reads from the websocket
	go webSocketReceiveHandler(conn, webSocketQuitChan)

	// Goroutine: handles writes to the websocket
	go webSocketSendHandler(conn, webSocketWriteChan)

	// Wait until the quit channel has items before stopping the program
	<-webSocketQuitChan
}
