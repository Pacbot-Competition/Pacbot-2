package game

import (
	"fmt"
	"sync"
)

// Keep track of number of active game engines (should only be one)
var activeGameEngines = 0

// Mutex to protect numActiveGameEngines
var muAGE sync.Mutex

/*
A game engine object, to act as an intermediary between the web broker
and the internal game state - its responsibility is to read responses from
clients and routinely send serialized copies of the game state to them
*/
type GameEngine struct {
	quitCh   chan struct{}
	outputCh chan<- []byte
	inputCh  <-chan []byte
	hasQuit  bool
}

// Create a new game engine, casting input and output channels to be uni-directional
func NewGameEngine(_broadcastCh chan<- []byte, _responseCh <-chan []byte) *GameEngine {
	return &GameEngine{
		quitCh:   make(chan struct{}),
		outputCh: _broadcastCh,
		inputCh:  _responseCh,
		hasQuit:  false,
	}
}

// Quit by closing the game engine, in case the loop ends
func (ge *GameEngine) quit() {

}

// Quit function exported to other packages
func (ge *GameEngine) Quit() {
	ge.quitCh <- struct{}{}
	ge.hasQuit = true
}

// Start the game engine - should be launched as a go-routine
func (ge *GameEngine) RunLoop() {

	// Quit if we ever run into an error or the program ends
	defer func() {
		ge.quit()
	}()

	// Update the number of active game engines
	// (Lock the mutex to prevent races in case of a multiple engine issue)
	var _activeGameEngines int
	muAGE.Lock()
	{
		activeGameEngines++
		_activeGameEngines = activeGameEngines
	}
	muAGE.Unlock()

	// If there was already a game engine, kill this one and throw an error
	if _activeGameEngines > 1 {
		fmt.Println("\033[35m\033[1mERR:  Cannot simultaneously dispatch more than one game engine. Quitting...\033[0m")
		return
	}

	for {
		select {

		// If we get a message from the web broker, handle it
		case msg := <-ge.inputCh:
			fmt.Printf("\033[2m\033[36m| Browser: %s`\033[0m\n", string(msg))

		// If we get a quit signal, quit this broker
		case <-ge.quitCh:
			return

		// If the web broker response channel hits full capacity, send a warning to the terminal
		default:
			if len(ge.inputCh) == cap(ge.inputCh) {
				fmt.Println("\033[35mWARN: Game engine input channel full\033[0m")
			}
		}
	}
}
