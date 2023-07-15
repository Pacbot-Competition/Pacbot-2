package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// Get the configuration info (config.go)
	conf := GetConfig()

	// Websocket stuff (ws_server.go)
	http.HandleFunc("/", webSocketHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocketPort), nil)

	// TCP stuff (tcp_server.go)
	server := NewTcpServer(fmt.Sprintf(":%d", conf.TcpPort))
	go server.Printer()

	// Log TCP errors
	log.Fatal(server.tcpStart())
}
