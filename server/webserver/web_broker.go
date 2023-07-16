package webserver

import "fmt"

// Keep track of active websocket sessions in a set (empty-valued map)
var openWebSessions = make(map[*webSession](struct{}))

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
		QuitCh:      make(chan struct{}, 1),
		BroadcastCh: make(chan []byte, 2),
		ResponseCh:  make(chan []byte, 10),
	}
}

// Quit by closing all web sessions, in case the loop ends
func (wb *WebBroker) quit() {
	fmt.Println("\033[31mWeb broker failure: killing all websocket connections\033[30m")
	muWS.Lock()
	{
		for ws := range openWebSessions {
			ws.quitCh <- struct{}{}
		}
	}
	muWS.Unlock()
}

// Start the web-broker - should be launched as a go-routine
func (wb *WebBroker) RunLoop() {

	// Quit if we ever run into an error or the program ends
	defer func() {
		wb.quit()
	}()

	for {
		select {

		// If we get a message, broadcast it to all web sessions
		case msg := <-wb.BroadcastCh:
			muWS.Lock()
			{
				for ws := range openWebSessions {
					if len(ws.sendCh) == cap(ws.sendCh) {
						fmt.Println("\033[31mChannel full\033[30m")
					}
					ws.sendCh <- msg
				}
			}
			muWS.Unlock()

		// If we get a quit signal, quit this broker
		case <-wb.QuitCh:
			return
		}
	}
}
