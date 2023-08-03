package game

/*
A game state object, to hold the internal game state and provide
helper methods that can be accessed by the game engine
*/
type gameState struct {

	/* Message header - 4 bytes */

	currTicks    uint16 // Current ticks elapsed
	updatePeriod uint8  // Ticks / update
	gameMode     uint8  // Game mode, encoded using an enum (TODO)

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
}

// Create a new game state with default values
func newGameState() *gameState {

	// New game state object
	gs := gameState{

		// Message header
		currTicks:    0,
		updatePeriod: initUpdatePeriod,
		gameMode:     0,

		// Game info
		currScore: 0,
		currLevel: initLevel,
		currLives: initLives,

		// Fruit
		fruitExists: false,

		// Ghosts
		ghosts: make([]*ghostState, 4),
	}

	gs.pacmanLoc = newLocationStateCopy(initLocPacman)
	gs.fruitLoc = newLocationStateCopy(initLocFruit)

	gs.ghosts[red] = newGhostState(red, gs.pacmanLoc)
	gs.ghosts[pink] = newGhostState(pink, gs.pacmanLoc)
	gs.ghosts[cyan] = newGhostState(cyan, gs.pacmanLoc)
	gs.ghosts[orange] = newGhostState(orange, gs.pacmanLoc)

	// Copy over maze bit arrays
	copy(gs.pellets[:], initPellets[:])
	copy(gs.walls[:], initWalls[:])
	return &gs
}
