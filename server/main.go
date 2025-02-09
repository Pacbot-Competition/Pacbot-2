package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"pacbot_server/game"
	"pacbot_server/webserver"
	"sync"
	"time"
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
	webResponseCh := make(chan []byte, 100)
	tcpSendCh := make(chan []byte, 2)

	// Set up the TCP server
	tcp := webserver.NewTcpServer(fmt.Sprintf(":%d", conf.TcpPort), tcpSendCh)
	go tcp.TcpStart()
	go tcp.Printer()
	log.Printf("\033[35mLOG:  Tcp server running on %s:%d\033[0m\n", conf.ServerIP, conf.TcpPort)

	// A wait group for quitting synchronously (allowing go-routines to complete)
	var wgQuit sync.WaitGroup

	// Websocket setup (package webserver)
	server := http.Server{Addr: fmt.Sprintf(":%d", conf.WebSocketPort)}
	log.Printf("\033[35mLOG:  Web server running on %s:%d\033[0m\n", conf.ServerIP, conf.WebSocketPort)
	wb := webserver.NewWebBroker(webBroadcastCh, tcpSendCh, webResponseCh, &wgQuit)
	go wb.RunLoop() // Run the web broker loop asynchronously
	http.HandleFunc("/", webserver.WebSocketHandler)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %e", err)
		}
		log.Println("\033[35mLOG:  HTTP server successfully quit\033[0m")
	}()

	// Game engine setup (package game)
	game.ConfigNumActiveGhosts(min(conf.NumActiveGhosts, 4))
	ge := game.NewGameEngine(webBroadcastCh, webResponseCh, &wgQuit, conf.GameFPS)
	go ge.RunLoop() // Run the game engine loop asynchronously

	// Set the enable for game command logging to be false by default
	game.SetCommandLogEnable(false)

	// Keep the game engine alive until a user types 'q'
	var input string
	fmt.Println("Ready")
	for {
		fmt.Scanf("%s\n", &input) // Blocking I/O to keep the program alive
		if input == "q" {         // Quit signal
			break
		} else {
			webResponseCh <- []byte(input)
		}
	}

	// Shutdown HTTP server to prevent new and finish old connections
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownRelease()
	server.Shutdown(shutdownCtx)

	// Quit the web server and game engine once complete
	wb.Quit()
	ge.Quit()

	// Synchronize to allow all processes to end safely
	wgQuit.Wait()
}
