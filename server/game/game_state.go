package game

type gameState struct {
	currTicks   uint16
	updateTicks uint8
	gameMode    uint8
	pellets     [mazeRows]uint32
	walls       [mazeRows]uint32
}

func newGameState() *gameState {
	gs := gameState{
		currTicks:   0,
		updateTicks: 0,
		gameMode:    0,
	}
	copy(gs.pellets[:], startingPellets[:])
	copy(gs.walls[:], walls[:])
	return &gs
}

// Serializes in big-endian form (most significant byte first)
func (gs *gameState) serPellets(outputBuf []byte, startIdx int) int {
	for row := 0; row < mazeRows; row++ {
		for byte_num := 0; byte_num < 4; byte_num++ {
			outputBuf[(row*4+byte_num)+startIdx] = byte((gs.pellets[row] >> (8 * (3 - byte_num))) & 0xff)
		}
	}
	return startIdx + (4 * mazeRows)
}
