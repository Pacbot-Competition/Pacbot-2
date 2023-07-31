package game

/*
A game state object, to hold the internal game state and provide
helper methods that can be accessed by the game engine
*/
type gameState struct {
	currTicks   uint16
	updateTicks uint8
	gameMode    uint8
	pellets     [mazeRows]uint32
	walls       [mazeRows]uint32
}

// Create a new game state with default values
func newGameState() *gameState {
	gs := gameState{
		currTicks:   0,
		updateTicks: initUpdateTicks,
		gameMode:    0,
	}
	copy(gs.pellets[:], initPellets[:])
	copy(gs.walls[:], initWalls[:])
	return &gs
}
