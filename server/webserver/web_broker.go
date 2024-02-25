package webserver

import (
	"log"
	"sync"
)

// Wait group to safely close all open clients when quitting
var wgQuit *sync.WaitGroup

/*
A web-broker object, to act as an intermediary between web sessions
and messages from the game engine - its responsibility is to forward byte
messages from the game engine to the clients and vice versa
*/
type WebBroker struct {
	quitCh      chan struct{}
	broadcastCh <-chan []byte
	responseCh  chan<- []byte
}

// Create a new web broker, casting input and output channels to be uni-directional
func NewWebBroker(_broadcastCh <-chan []byte, _responseCh chan<- []byte, _wgQuit *sync.WaitGroup) *WebBroker {
	wb := WebBroker{
		quitCh:      make(chan struct{}, 0),
		broadcastCh: _broadcastCh,
		responseCh:  _responseCh,
	}
	wgQuit = _wgQuit
	return &wb
}

// Quit by closing all web sessions, in case the loop ends
func (wb *WebBroker) quit() {

	// Log that all websocket connections are closed upon broker exit, then close them individually
	log.Println("\033[35m\033[1mLOG:  Web broker exit: killing all websocket connections")
	muOWS.RLock()
	{
		// Individually quit each of the open web sessions
		for ws := range openWebSessions {
			ws.quit()
		}
	}
	muOWS.RUnlock()

	// Log that the web broker has quit (if this message doesn't get sent, we are blocked by some mutex)
	log.Println("\033[35mLOG:  Web broker successfully quit\033[0m")

	wgQuit.Done()
}

// Quit function exported to other packages
func (wb *WebBroker) Quit() {
	close(wb.quitCh)
}

// Start the web-broker - should be launched as a go-routine
func (wb *WebBroker) RunLoop() {
	// Make sure we wait for web broker to complete before exit
	wgQuit.Add(1)

	// Quit if we ever run into an error or the program ends
	defer wb.quit()

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

					// Issue update to client if they are keeping up
					select {
					case ws.sendCh <- msg:
						// Don't wait, we won't hold everything up for a slow client
					default:
						/*
							What this means: a web session channel was full,
							preventing this write
						*/
						log.Printf("\033[35mWARN: A web-session send channel was full"+
							" (client = %s)\033[0m\n", getIP(ws.conn))
					}
				}
			}
			muOWS.RUnlock()

		// If we get a quit signal, quit this broker
		case <-wb.quitCh:
			return
		}
	}
}
