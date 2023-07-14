package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Keep track of the number of open websocket connections
var openWebClients int = 0

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
func webSocketReceiveHandler(conn *websocket.Conn, webSocketQuitChannel chan int) {

	for {
		// Read a message (discard the type since we don't need it)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// If the browser gets disconnected, quit (which will close the connection)
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				webSocketQuitChannel <- 0
				return
			}
			// For all other errors, log them
			fmt.Println("read error: ", err)
			return
		}
		// Print the message we received
		fmt.Printf("\033[2m\033[36m| Browser: %s\033[0m\n", string(msg))
	}
}

// Sending websocket data (binary)
func webSocketSendHandler(conn *websocket.Conn, webSocketWriteChannel chan int) {

	// Message to send
	msg := []byte(".")

	for {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			if err == websocket.ErrCloseSent {
				return
			}
			fmt.Println("write error: ", err)
			return
		}
		// Block until the next message is ready
		<-webSocketWriteChannel
	}
}

/*
This handler makes sure that once we connect to the websocket,
all communication goes smoothly.
*/
func webSocketHandler(w http.ResponseWriter, r *http.Request) {

	// Print a status message
	openWebClients++
	fmt.Printf("\033[34m[%d -> %d] browser connected\033[0m\n", openWebClients-1, openWebClients)

	// Upgrades the connection, and quits if it didn't work out.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the connection at the end of the function, or if
	// something goes wrong.
	defer conn.Close()

	// Quit channel to postpone closing the websocket
	webSocketQuitChan := make(chan int, 1)

	// Write channel to postpone writing to the websocket
	webSocketWriteChan := make(chan int, 1)

	// Goroutine: handles reads from the websocket
	go webSocketReceiveHandler(conn, webSocketQuitChan)

	// Goroutine: handles writes to the websocket
	go webSocketSendHandler(conn, webSocketWriteChan)

	// Wait until the quit channel empties before stopping the program
	<-webSocketQuitChan

	// Return after getting a quit signal
	openWebClients--
	fmt.Printf("\033[33m[%d -> %d] browser disconnected\033[0m\n", openWebClients+1, openWebClients)
}
