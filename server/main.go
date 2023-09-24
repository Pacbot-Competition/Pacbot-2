package main

import (
	"fmt"
	"log"
	"net/http"
	"pacbot_server/game"
	"pacbot_server/webserver"
	"sync"
)

const (
	PORT = 3002
)

func main() {

	// Disable logging timestamps
	log.SetFlags(0)

	// Get the configuration info (config.go)
	conf := GetConfig()

	// Use this configuration info to set up server subunits
	webserver.ConfigOneClientPerIP(conf.OneClientPerIP)
	webserver.ConfigTrustedClientIPs(conf.TrustedClientIPs)

	// Make channels for communication between web broker and game engine
	webBroadcastCh := make(chan []byte, 100)
	webResponseCh := make(chan []byte, 10)

	// A wait group for quitting synchronously (allowing go-routines to complete)
	var wgQuit sync.WaitGroup

	// Websocket setup (package webserver)
	wb := webserver.NewWebBroker(webBroadcastCh, webResponseCh, &wgQuit)
	go wb.RunLoop() // Run the web broker loop asynchronously
	fmt.Printf("[Start] Server running on :%v\n", PORT)
	http.HandleFunc("/", webserver.WebSocketHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocketPort), nil)

	// Game engine setup (package game)
	ge := game.NewGameEngine(webBroadcastCh, webResponseCh, &wgQuit, conf.GameFPS)
	go ge.RunLoop() // Run the game engine loop asynchronously

	// Set the enable for game command logging to be false by default
	game.SetCommandLogEnable(false)

	// Keep the game engine alive until a user types 'q'
	var input string
	for {
		fmt.Scanf("%s", &input) // Blocking I/O to keep the program alive
		if input == "q" {       // Quit signal
			break
		} else if input == "p" { // Pause signal
			webResponseCh <- []byte("p")
		} else if input == "P" { // Play signal
			webResponseCh <- []byte("P")
		}
	}

	// Quit the web server and game engine once complete
	wb.Quit()
	ge.Quit()

	// Synchronize to allow all processes to end safely
	wgQuit.Wait()
}
