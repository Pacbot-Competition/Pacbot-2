package main

import (
	"fmt"
	"log"
	"net/http"
	"pacbot_server/webserver"
)

func main() {

	// Get the configuration info (config.go)
	conf := GetConfig()
	OneBrowserPerIP := conf.OneBrowserPerIP

	// Websocket stuff (webserver)
	webserver.ConfigOneBrowserPerIP(OneBrowserPerIP)
	http.HandleFunc("/", webserver.WebSocketHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocketPort), nil)

	// TCP stuff (tcp_server.go)
	server := NewTcpServer(fmt.Sprintf(":%d", conf.TcpPort))
	go server.Printer()

	// Log TCP errors
	log.Fatal(server.tcpStart())
}
