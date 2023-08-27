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

// Wait group to safely close all open browsers when quitting
var wgQuit *sync.WaitGroup

/*
A web-broker object, to act as an intermediary between web sessions
and messages from the game engine - its responsibility is to forward byte
messages from the game engine to the browsers and vice versa
*/
type WebBroker struct {
	quitCh      chan struct{}
	broadcastCh <-chan []byte
	responseCh  chan<- []byte
}

// Create a new web broker, casting input and output channels to be uni-directional
func NewWebBroker(_broadcastCh <-chan []byte, _responseCh chan<- []byte, _wgQuit *sync.WaitGroup) *WebBroker {
	wb := WebBroker{
		quitCh:      make(chan struct{}),
		broadcastCh: _broadcastCh,
		responseCh:  _responseCh,
	}
	wgQuit = _wgQuit
	return &wb
}

// Quit by closing all web sessions, in case the loop ends
func (wb *WebBroker) quit() {

	// Log that all websocket connections are closed upon broker exit, then close them individually
	fmt.Println("\033[35m\033[1mLOG:  Web broker exit: killing all websocket connections")
	fmt.Println("      No new connections allowed\033[0m")
	newConnectionsAllowed = false // Doesn't need a mutex since we can't ever make it true again
	muOWS.RLock()
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
	muOWS.RUnlock()

	// Log that the web broker has quit (if this message doesn't get sent, we are blocked by some mutex)
	fmt.Println("\033[35mLOG:  Web server successfully quit\033[0m")
}

// Quit function exported to other packages
func (wb *WebBroker) Quit() {
	wb.quitCh <- struct{}{}
}

// Start the web-broker - should be launched as a go-routine
func (wb *WebBroker) RunLoop() {

	// Quit if we ever run into an error or the program ends
	defer func() {
		wb.quit()
	}()

	/*
		Update the number of active web broker groups
		(Lock the mutex to prevent races in case of a multiple broker issue)
	*/
	var _activeWebBrokerLoops int
	muAWB.Lock()
	{
		activeWebBrokerLoops++
		_activeWebBrokerLoops = activeWebBrokerLoops
	}
	muAWB.Unlock()

	// If there was already a web broker, kill this one and throw an error
	if _activeWebBrokerLoops > 1 {
		fmt.Println("\033[35m\033[1mERR:  Cannot simultaneously dispatch more " +
			"than one web broker loop. Quitting...\033[0m")
		return
	}

	// Copy (by reference) the response channel to match the broker's
	responseCh = wb.responseCh

	// "While" loop, keep running until we quit the web broker
	for {
		select {

		// If we get a message, broadcast it to all web sessions
		case msg := <-wb.broadcastCh:
			muOWS.RLock()
			{
				for ws := range openWebSessions {

					// Check if the write will be blocked due to a full send channel
					b := len(ws.sendCh) == cap(ws.sendCh)
					start := time.Now()
					ws.sendCh <- msg

					/*
						If the write was blocked for too long (> 1ms),
						send a warning to the terminal

						What this means: a web session channel was full,
						blocking this write
					*/
					if b {
						wait := time.Since(start)
						if wait > time.Millisecond {
							fmt.Printf("\033[35mWARN: A web-session send channel was full"+
								" (%s, client = %s)\033[0m\n", wait, getIP(ws.conn))
						}
					}
				}
			}
			muOWS.RUnlock()

		// If we get a quit signal, quit this broker
		case <-wb.quitCh:
			return

		/*
			If the web broker broadcast channel hits full capacity, send a warning to the terminal
			What this means: either the game engine is sending too much input, or the web broker loop can't keep up
		*/
		default:
			if len(wb.broadcastCh) == cap(wb.broadcastCh) {
				fmt.Println("\033[35mWARN: Web broker broadcast channel full\033[0m")
			}
		}
	}
}
