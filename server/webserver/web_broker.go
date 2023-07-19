package webserver

import (
	"fmt"
	"sync"
	"time"
)

// Keep track of number of active web brokers (should only be one)
var activeWebBrokerLoops = 0

// Keep track of whether new connections are allowed
var newConnectionsAllowed = true

// Mutex to protect numActiveWebBrokerLoops
var muAWB sync.Mutex

/*
A web-broker object, to act as an intermediary between web sessions
and messages from the game engine
*/
type WebBroker struct {
	QuitCh      chan struct{}
	BroadcastCh chan []byte
	ResponseCh  chan []byte
}

// Create a new web broker
func NewWebBroker() *WebBroker {
	return &WebBroker{
		QuitCh:      make(chan struct{}),
		BroadcastCh: make(chan []byte, 10),
		ResponseCh:  make(chan []byte, 10),
	}
}

// Quit by closing all web sessions, in case the loop ends
func (wb *WebBroker) quit() {
	fmt.Println("\033[35m\033[1mLOG:  Web broker exit: killing all websocket connections\n      No new connections allowed\033[0m")
	newConnectionsAllowed = false // Doesn't need a mutex since we can't ever make it true again
	muOWS.Lock()
	{
		muAWB.Lock()
		{
			activeWebBrokerLoops--
		}
		muAWB.Unlock()
		for ws := range openWebSessions {
			ws.quitCh <- struct{}{}
		}
	}
	muOWS.Unlock()
	fmt.Println("\033[35mLOG:  Web broker successfully quit\033[0m")
}

// Quit function exported to other packages
func (wb *WebBroker) Quit() {
	wb.QuitCh <- struct{}{}
}

// Start the web-broker - should be launched as a go-routine
func (wb *WebBroker) RunLoop() {

	// Quit if we ever run into an error or the program ends
	defer func() {
		wb.quit()
	}()

	var _activeWebBrokerLoops int
	muAWB.Lock()
	{
		activeWebBrokerLoops++
		_activeWebBrokerLoops = activeWebBrokerLoops
	}
	muAWB.Unlock()

	if _activeWebBrokerLoops > 1 {
		fmt.Println("\033[35m\033[1mERR:  Cannot simultaneously dispatch more than one web broker loop. Quitting...\033[0m")
		return
	}

	for {
		select {

		// If we get a message, broadcast it to all web sessions
		case msg := <-wb.BroadcastCh:
			muOWS.Lock()
			{
				for ws := range openWebSessions {

					// Check if the write will be blocked
					b := len(ws.sendCh) == cap(ws.sendCh)
					start := time.Now()
					ws.sendCh <- msg

					// If the write was blocked for too long (> 1ms), send a warning to the terminal
					if b {
						wait := time.Since(start)
						if wait > time.Millisecond {
							fmt.Printf("\033[35mWARN: A web-session send channel was full (%s, client = %s)\033[0m\n", wait, getIP(ws.conn))
						}
					}
				}
			}
			muOWS.Unlock()

		// If we get a quit signal, quit this broker
		case <-wb.QuitCh:
			return

		default:
			if len(wb.BroadcastCh) == cap(wb.BroadcastCh) {
				fmt.Println("\033[35mWarning: Web broker broadcast channel full\033[0m")
			}
		}
	}
}
