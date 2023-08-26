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
	state       *gameState
	ticker      *time.Ticker    // serves as the game clock
	wgQuit      *sync.WaitGroup // wait group to make sure it quits safely
}

// Create a new game engine, casting input and output channels to be uni-directional
func NewGameEngine(_webOutputCh chan<- []byte, _webInputCh <-chan []byte, _wgQuit *sync.WaitGroup, clockRate int32) *GameEngine {
	_tickTime := 1000000 * time.Microsecond / time.Duration(clockRate) // Time between ticks
	ge := GameEngine{
		quitCh:      make(chan struct{}),
		webOutputCh: _webOutputCh,
		webInputCh:  _webInputCh,
		state:       newGameState(),
		ticker:      time.NewTicker(_tickTime),
		wgQuit:      _wgQuit,
	}
	return &ge
}

// Quit by closing the game engine, in case the loop ends
func (ge *GameEngine) quit() {

	// Log that the game engine successfully quit
	fmt.Println("\033[35mLOG:  Game engine successfully quit\033[0m")

	// Decrement the quit wait group counter
	ge.wgQuit.Done()
}

// Quit function exported to other packages
func (ge *GameEngine) Quit() {
	ge.quitCh <- struct{}{}
}

// Start the game engine - should be launched as a go-routine
func (ge *GameEngine) RunLoop() {

	// Quit if we ever run into an error or the program ends
	defer func() {
		ge.quit()
	}()

	// Increment the quit wait group counter
	ge.wgQuit.Add(1)

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

	// Output buffer to store the serialized output
	outputBuf := make([]byte, 256)

	// Create a wait group for synchronizing ghost plans
	var wgPlans sync.WaitGroup

	// Create a variable for the first update
	firstUpdate := false

	for {

		/* STEP 1: Update the ghost positions if necessary */

		// If the game state is ready to update, update the ghost positions
		if ge.state.updateReady() || !firstUpdate {

			// Wait until all pending ghost plans are complete
			wgPlans.Wait()

			// Loop over the individual ghosts
			for _, ghost := range ge.state.ghosts {
				ghost.update()
			}
		}

		/* STEP 2: Serialize the current game state to the output buffer */

		// Starting index for serialization
		idx := 0

		// Send the full state
		idx = ge.state.serFull(outputBuf, idx)

		/* STEP 3: Start planning the next ghost moves if an update just happened */

		// If we're ready for an update, plan the next ghost moves asynchronously
		if ge.state.updateReady() || !firstUpdate {

			// Add pending ghost plans
			wgPlans.Add(int(numColors))

			// Plan each ghost's next move concurrently
			for _, ghost := range ge.state.ghosts {
				go ghost.plan(&wgPlans)
			}
		}

		/* STEP 4: Write the serialized game state to the output channel */

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

		/* STEP 5: Read the input channel and update the game state accordingly */
		select {

		// If we get a message from the web broker, handle it
		case msg := <-ge.webInputCh:
			ge.state.interpretCommand(msg)

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

		/* STEP 6: Update the game state for the next tick */

		// Increment the number of ticks
		if !ge.state.isPaused() {
			ge.state.currTicks++
		}

		// Set the first update to be done
		firstUpdate = true

		/* STEP 5: Wait for the ticker to complete the current frame */
		<-ge.ticker.C
	}
}
