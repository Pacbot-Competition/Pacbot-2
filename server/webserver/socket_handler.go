package webserver

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Keep track of active websocket sessions in a set (empty-valued map)
var openWebSessions = make(map[*webSession](struct{}))

/*
Protects the "openWebSessions" map from race conditions with a
mutex - it will lock the resources until they are successfully written to
*/
var muOWS sync.RWMutex

/*
Websockets are the way that the server will communicate with
the clients over the LAN. When connecting to a web-socket,
we connect normally over HTTP and then upgrade the
connection upon agreement between the server and client.
*/
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all clients to connect
	},
}

/*
This handler makes sure that once we connect to the websocket,
all communication goes smoothly.
*/
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrades the connection, and quits if it didn't work out.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket err:", err)
		return
	}

	// Create a websocket session object
	ws := newWebSession(conn)
	
	// Ensure we wait for clients to finish
	wgQuit.Add(1)
	defer wgQuit.Done()

	// Register the websocket session
	ws.register()

	/*
	  Close the connection at the end of the function, or if
	  something goes wrong.
	*/
	defer ws.unregister()

	ws.loop()
}
