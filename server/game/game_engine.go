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
	quitCh      chan struct{}
	webOutputCh chan<- []byte
	webInputCh  <-chan []byte
	hasQuit     bool
	state       *gameState
	ticker      *time.Ticker   // serves as the game clock
	hrticker    *HighResTicker // serves as the game clock (unused, but available)
}

// Create a new game engine, casting input and output channels to be uni-directional
func NewGameEngine(_webOutputCh chan<- []byte, _webInputCh <-chan []byte, clockRate int32) *GameEngine {
	_tickTime := 1000000 * time.Microsecond / time.Duration(clockRate) // Time between ticks
	ge := GameEngine{
		quitCh:      make(chan struct{}),
		webOutputCh: _webOutputCh,
		webInputCh:  _webInputCh,
		hasQuit:     false,
		state:       newGameState(),
		hrticker:    NewHighResTicker(_tickTime),
		ticker:      time.NewTicker(_tickTime),
	}
	return &ge
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

	// Start the high-res game clock
	// ge.hrticker.Start()

	// Output buffer to store the serialized output
	outputBuf := make([]byte, 256)

	for {

		// Test: update game state on the fly
		if ge.state.updateReady() {

			// Loop over the individual ghosts
			for _, ghost := range ge.state.ghosts {
				ghost.update()
			}
		}

		/* STEP 1: Serialize the current game state to the output buffer */
		idx := 0

		// Send the full state
		idx = ge.state.serFull(outputBuf, idx)

		// If we're ready for an update, plan the next move
		if ge.state.updateReady() {
			for _, ghost := range ge.state.ghosts {
				ghost.plan()
			}
		}

		/* STEP 2: Write the serialized game state to the output channel */

		// Check if the write will be blocked, and try to write the serialized state
		b := len(ge.webOutputCh) == cap(ge.webOutputCh)
		start := time.Now()
		ge.webOutputCh <- outputBuf[:idx]

		// If the write was blocked for too long (> 1ms), send a warning to the terminal
		if b {
			wait := time.Since(start)
			if wait > time.Millisecond {
				fmt.Printf("\033[35mWARN: The game engine output channel was full (%s)\033[0m\n", wait)
			}
		}

		/* STEP 3: Read the input channel and update the game state accordingly */
		select {

		// If we get a message from the web broker, handle it
		case msg := <-ge.webInputCh:
			fmt.Printf("\033[2m\033[36m| Response: %s`\033[0m\n", string(msg))

		// If we get a quit signal, quit this broker
		case <-ge.quitCh:
			return

		/*
			If the web input channel hits full capacity, send a warning to the terminal
			What this means: either the browsers are sending too much input, or the game loop can't keep up
		*/
		default:
			if len(ge.webInputCh) == cap(ge.webInputCh) {
				fmt.Println("\033[35mWARN: Game engine input channel full\033[0m")
			}
		}

		/* STEP 4: Update the game state for the next tick */

		// Increment the number of ticks
		ge.state.currTicks++

		/* STEP 5: Wait for the ticker to complete the current frame */
		<-ge.ticker.C
	}
}
