package game

import (
	"math/rand"
	"time"
)

// Enum-like declaration to hold the game mode options
const (
	paused  = 0
	scatter = 1
	chase   = 2
)

/*
A game state object, to hold the internal game state and provide
helper methods that can be accessed by the game engine
*/
type gameState struct {

	/* Message header - 4 bytes */

	// Current ticks elapsed - WARN: at 24 ticks/sec, this will have
	// an integer overflow after about 45 minutes, so don't run it
	// continuously for too long
	currTicks uint16

	updatePeriod uint8 // Ticks / update
	mode         uint8 // Game mode, encoded using an enum (TODO)

	/* Game information - 4 bytes */

	currScore uint16 // Current score
	currLevel uint8  // Current level number (by default, starts at 1)
	currLives uint8  // Current lives        (by default, starts at 3)

	/* Pacman location - 2 bytes */

	pacmanLoc *locationState

	/* Fruit location - 2 bytes */

	fruitExists bool
	fruitLoc    *locationState

	/* Ghosts - 4 * 3 = 12 bytes */

	ghosts []*ghostState

	/* Pellet State - 31 * 4 = 124 bytes */

	// Pellets encoded within an array, with each uint32 acting as a bit array
	pellets [mazeRows]uint32

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
	}

	gs.pacmanLoc = newLocationStateCopy(pacmanSpawnLoc)
	gs.fruitLoc = newLocationStateCopy(fruitSpawnLoc)

	gs.ghosts[red] = newGhostState(&gs, red)
	gs.ghosts[pink] = newGhostState(&gs, pink)
	gs.ghosts[cyan] = newGhostState(&gs, cyan)
	gs.ghosts[orange] = newGhostState(&gs, orange)

	// Copy over maze bit arrays
	copy(gs.pellets[:], initPellets[:])
	copy(gs.walls[:], initWalls[:])
	return &gs
}
