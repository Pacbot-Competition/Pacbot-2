package main

import (
	"log"
	"net/http"
)

func main() {

	// Websocket stuff
	http.HandleFunc("/", webSocketHandler)
	go http.ListenAndServe(":3002", nil)

	// TCP stuff
	server := NewTcpServer(":3001")
	go server.Printer()

	// Log TCP errors
	log.Fatal(server.tcpStart())
}
