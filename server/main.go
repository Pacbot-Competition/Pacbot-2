package main

import (
	"fmt"
	"log"
	"net/http"
	"pacbot_server/game"
	"pacbot_server/tcpserver"
	"pacbot_server/webserver"
	"time"

	"github.com/loov/hrtime"
)

func main() {

	// Get the configuration info (config.go)
	conf := GetConfig()

	// Use this configuration info to set up server subunits
	webserver.ConfigOneBrowserPerIP(conf.OneBrowserPerIP)
	webserver.ConfigTrustedBrowserIPs(conf.TrustedBrowserIPs)

	// Make channels for communication between web broker and game engine
	webBroadcastCh := make(chan []byte, 100)
	webResponseCh := make(chan []byte, 10)

	// Websocket setup (package webserver)
	wb := webserver.NewWebBroker(webBroadcastCh, webResponseCh)
	go wb.RunLoop() // Run the web broker loop asynchronously
	http.HandleFunc("/", webserver.WebSocketHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocketPort), nil)

	// TCP setup (package tcpserver)
	server := tcpserver.NewTcpServer(fmt.Sprintf(":%d", conf.TcpPort))
	go server.Printer()

	// Game engine setup (package game)
	ge := game.NewGameEngine(webBroadcastCh, webResponseCh, conf.GameFPS)
	go ge.RunLoop() // Run the game engine loop asynchronously

	/*
		Demo for time-keeping abilities: after 60s, all websockets will be killed through
		quitting the broker, freezing the time received on the web client
	*/
	go func(wb *webserver.WebBroker) {
		start := hrtime.Now()
		time.Sleep(10 * time.Second)
		wb.Quit()
		fmt.Println("slow:", hrtime.Since(start))
		wb.Quit()
	}(wb)

	// Log TCP errors
	log.Fatal(server.TcpStart())
}
