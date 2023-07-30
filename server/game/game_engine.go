package game

import (
	"fmt"
	"sync"
	"time"
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
	ticker   HighResTicker // serves as the game clock
}

// Create a new game engine, casting input and output channels to be uni-directional
func NewGameEngine(_outputCh chan<- []byte, _inpuCh <-chan []byte, clockRate int32) *GameEngine {
	return &GameEngine{
		quitCh:   make(chan struct{}),
		outputCh: _outputCh,
		inputCh:  _inpuCh,
		hasQuit:  false,
		ticker:   *NewHighResTicker(clockRate),
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

	// Start the game clock (called as a new go-routine)
	go ge.ticker.Start()

	for {

		/* STEP 1: Serialize and send the current game state */

		// Check if the write will be blocked
		b := len(ge.outputCh) == cap(ge.outputCh)
		start := time.Now()
		ge.outputCh <- SerializePellets(Pellets)

		// If the write was blocked for too long (> 1ms), send a warning to the terminal
		if b {
			wait := time.Since(start)
			if wait > time.Millisecond {
				fmt.Printf("\033[35mWARN: The game engine output channel was full (%s)\033[0m\n", wait)
			}
		}

		/* STEP 2: Read the input channel and update the game state accordingly */
		select {

		// If we get a message from the web broker, handle it
		case msg := <-ge.inputCh:
			fmt.Printf("\033[2m\033[36m| Response: %s`\033[0m\n", string(msg))

		// If we get a quit signal, quit this broker
		case <-ge.quitCh:
			return

		// If the web broker response channel hits full capacity, send a warning to the terminal
		default:
			if len(ge.inputCh) == cap(ge.inputCh) {
				fmt.Println("\033[35mWARN: Game engine input channel full\033[0m")
			}
		}

		/* STEP 3: Update the game state for the next tick */
		Pellets[0] += 1 // Test reactivity of Svelte frontend

		/* STEP 4: Wait for the ticker to complete the current frame */
		<-ge.ticker.ReadyCh
	}
}
