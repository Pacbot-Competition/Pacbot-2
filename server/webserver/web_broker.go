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
and messages from the game engine - its responsibility is to forward byte
messages from the game engine to the browsers and vice versa
*/
type WebBroker struct {
	quitCh      chan struct{}
	broadcastCh <-chan []byte
	responseCh  chan<- []byte
	hasQuit     bool
}

// Create a new web broker, casting input and output channels to be uni-directional
func NewWebBroker(_broadcastCh <-chan []byte, _responseCh chan<- []byte) *WebBroker {
	return &WebBroker{
		quitCh:      make(chan struct{}),
		broadcastCh: _broadcastCh,
		responseCh:  _responseCh,
		hasQuit:     false,
	}
}

// Quit by closing all web sessions, in case the loop ends
func (wb *WebBroker) quit() {

	// Log that all websocket connections are closed upon broker exit, then close them individually
	fmt.Println("\033[35m\033[1mLOG:  Web broker exit: killing all websocket connections")
	fmt.Println("      No new connections allowed\033[0m")
	newConnectionsAllowed = false // Doesn't need a mutex since we can't ever make it true again
	muOWS.Lock()
	{
		// Decrement the number of active web broker loops
		muAWB.Lock()
		{
			activeWebBrokerLoops--
		}
		muAWB.Unlock()

		// Individually quit each of the open web sessions
		for ws := range openWebSessions {
			ws.quitCh <- struct{}{}
		}
	}
	muOWS.Unlock()

	// Log that the web broker has quit (if this message doesn't get sent, we are blocked by some mutex)
	fmt.Println("\033[35mLOG:  Web broker successfully quit\033[0m")
}

// Quit function exported to other packages
func (wb *WebBroker) Quit() {
	wb.quitCh <- struct{}{}
	wb.hasQuit = true
}

// Has Quit function exported to other packages
func (wb *WebBroker) HasQuit() bool {
	return wb.hasQuit
}

// Start the web-broker - should be launched as a go-routine
func (wb *WebBroker) RunLoop() {

	// Quit if we ever run into an error or the program ends
	defer func() {
		wb.quit()
	}()

	// Update the number of active web broker groups
	// (Lock the mutex to prevent races in case of a multiple broker issue)
	var _activeWebBrokerLoops int
	muAWB.Lock()
	{
		activeWebBrokerLoops++
		_activeWebBrokerLoops = activeWebBrokerLoops
	}
	muAWB.Unlock()

	// If there was already a web broker, kill this one and throw an error
	if _activeWebBrokerLoops > 1 {
		fmt.Println("\033[35m\033[1mERR:  Cannot simultaneously dispatch more than one web broker loop. Quitting...\033[0m")
		return
	}

	// Copy (by reference) the response channel of the package to match the broker's
	responseCh = wb.responseCh

	// "While" loop, keep running until we quit the web broker
	for {
		select {

		// If we get a message, broadcast it to all web sessions
		case msg := <-wb.broadcastCh:
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
		case <-wb.quitCh:
			return

		// If the web broker broadcast channel hits full capacity, send a warning to the terminal
		default:
			if len(wb.broadcastCh) == cap(wb.broadcastCh) {
				fmt.Println("\033[35mWARN: Web broker broadcast channel full\033[0m")
			}
		}
	}
}
