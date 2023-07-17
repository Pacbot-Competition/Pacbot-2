package main

import (
	"fmt"
	"log"
	"net/http"
	"pacbot_server/webserver"
	"time"

	"github.com/loov/hrtime"
)

func main() {

	// Get the configuration info (config.go)
	conf := GetConfig()
	OneBrowserPerIP := conf.OneBrowserPerIP

	// Websocket stuff (webserver)
	webserver.ConfigOneBrowserPerIP(OneBrowserPerIP)
	wb := webserver.NewWebBroker()
	go wb.RunLoop()
	http.HandleFunc("/", webserver.WebSocketHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", conf.WebSocketPort), nil)

	// TCP stuff (tcp_server.go)
	server := NewTcpServer(fmt.Sprintf(":%d", conf.TcpPort))
	go server.Printer()

	// High-resolution ticker (for keeping the frame rate roughly constant)
	hrt := NewHighResTicker(100)

	/*
		Demo for high-resolution ticker - at 60fps, it updates the web broker's
		send channel with the time (in seconds) and frames elapsed for that second
	*/
	go func(wb *webserver.WebBroker) {
		go hrt.Start()
	L:
		for idx := 0; idx < 5000; idx++ {
			select {
			case wb.BroadcastCh <- []byte(fmt.Sprintf("%d ~ %d", idx/100, idx%100)):
			case <-wb.QuitCh:
				break L
			}

			<-hrt.ReadyCh
		}
		fmt.Println("fast:", hrt.Lifetime())
	}(wb)

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
	log.Fatal(server.tcpStart())
}
