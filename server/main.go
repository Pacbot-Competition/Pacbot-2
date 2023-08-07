package main

import (
	"fmt"
	"net/http"
	"pacbot_server/game"
	"pacbot_server/tcpserver"
	"pacbot_server/webserver"
	"sync"
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

	// A wait group for quitting synchronously (allowing go-routines to complete)
	var wgQuit sync.WaitGroup

	// Websocket setup (package webserver)
	wb := webserver.NewWebBroker(webBroadcastCh, webResponseCh, &wgQuit)
	go wb.RunLoop() // Run the web broker loop asynchronously
	http.HandleFunc("/", webserver.WebSocketHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocketPort), nil)

	// TCP setup (package tcpserver)
	ts := tcpserver.NewTcpServer(fmt.Sprintf(":%d", conf.TcpPort))
	go ts.Printer()
	go ts.TcpStart(&wgQuit) // Start the TCP server asynchronously

	// Game engine setup (package game)
	ge := game.NewGameEngine(webBroadcastCh, webResponseCh, &wgQuit, conf.GameFPS)
	go ge.RunLoop() // Run the game engine loop asynchronously

	// Keep the game engine alive until a user types 'q'
	var input string
	for {
		fmt.Scanf("%s", &input) // Blocking I/O to keep the program alive
		if input == "q" {
			break
		}
	}

	// Quit the web server and game engine once complete
	wb.Quit()
	ge.Quit()
	ts.Quit()

	// Synchronize to allow all processes to end safely
	wgQuit.Wait()
}
