package game

import (
	"math/rand"
	"sync"
	"time"
)

// Enum-like declaration to hold the game mode options
const (
	paused  uint8 = 0
	scatter uint8 = 1
	chase   uint8 = 2
)

/*
A game state object, to hold the internal game state and provide
helper methods that can be accessed by the game engine
*/
type gameState struct {

	/* Message header - 4 bytes */

	/*
		Current ticks elapsed

		NOTE: at 24 ticks/sec, this will have suffer an integer
		overflow after about 45 minutes, so don't run it continuously
		for too long (indefinite pausing is fine though, as it doesn't
		increment the current tick amount)
	*/
	currTicks uint16
	muTicks   sync.RWMutex // Associated mutex

	updatePeriod uint8        // Ticks / update
	muPeriod     sync.RWMutex // Associated mutex

	mode   uint8        // Game mode, encoded using an enum (TODO)
	muMode sync.RWMutex // Associated mutex

	/* Game information - 4 bytes */

	currScore uint16       // Current score
	muScore   sync.RWMutex // Associated mutex

	currLevel uint8        // Current level (by default, starts at 1)
	muLevel   sync.RWMutex // Associated mutex

	currLives uint8        // Current lives (by default, starts at 3)
	muLives   sync.RWMutex // Associated mutex

	/* Pacman location - 2 bytes */

	pacmanLoc *locationState

	/* Fruit location - 2 bytes */

	fruitExists bool
	fruitLoc    *locationState
	muFruit     sync.RWMutex // Associated mutex (for fruitExists)

	/* Ghosts - 4 * 3 = 12 bytes */

	ghosts []*ghostState

	/* Pellet State - 31 * 4 = 124 bytes */

	// Pellets encoded within an array, with each uint32 acting as a bit array
	pellets    [mazeRows]uint32
	numPellets uint16       // Number of pellets
	muPellets  sync.RWMutex // Associated mutex

	/* Auxiliary (non-serialized) state information */

	// Wall state
	walls [mazeRows]uint32

	// A random number generator for making frightened ghost decisions
	rng *rand.Rand
}

// Create a new game state with default values
func newGameState() *gameState {

	// New game state object
	gs := gameState{

		// Message header
		currTicks:    0,
		updatePeriod: initUpdatePeriod,
		mode:         initMode,

		// Game info
		currScore: 0,
		currLevel: initLevel,
		currLives: initLives,

		// Fruit
		fruitExists: false,

		// Ghosts
		ghosts: make([]*ghostState, 4),

		// RNG source
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),

		// Pellet count at the start
		numPellets: initPelletCount,
	}

	// Declare the initial locations of Pacman and the fruit
	gs.pacmanLoc = newLocationStateCopy(pacmanSpawnLoc)
	gs.fruitLoc = newLocationStateCopy(fruitSpawnLoc)

	// Initialize the ghosts
	for color := uint8(0); color < 4; color++ {
		gs.ghosts[color] = newGhostState(&gs, color)
	}

	// Copy over maze bit arrays
	copy(gs.pellets[:], initPellets[:])
	copy(gs.walls[:], initWalls[:])

	// Return the new game state
	return &gs
}
