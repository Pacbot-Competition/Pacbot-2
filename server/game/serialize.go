package game

// Serialize the current number of ticks (2 bytes)
func (gs *gameState) serCurrTicks(outputBuf []byte, startIdx int) int {

	// Serialize the field, and return the new start index
	return serUint16(gs.currTicks, outputBuf, startIdx)
}

// Serialize the update ticks (1 byte)
func (gs *gameState) serUpdatePeriod(outputBuf []byte, startIdx int) int {

	// Serialize the field, and return the new start index
	return serUint8(gs.updatePeriod, outputBuf, startIdx)
}

// Serialize the game mode (1 byte)
func (gs *gameState) serGameMode(outputBuf []byte, startIdx int) int {

	// Serialize the field, and return the new start index
	return serUint8(gs.gameMode, outputBuf, startIdx)
}

// Serialize the pellets (4 * mazeRows bytes)
func (gs *gameState) serPellets(outputBuf []byte, startIdx int) int {

	// Loop over each row
	for row := 0; row < int(mazeRows); row++ {

		// Serialize the row from uint32 to 4 bytes
		startIdx = serUint32(gs.pellets[row], outputBuf, startIdx)
	}

	// Return the starting index of the next field
	return startIdx
}

// Serialize the location of Pacman (2 bytes)
func (gs *gameState) serPacman(outputBuf []byte, startIdx int) int {

	// Serialize the pacman state (with locking to be thread-safe)
	startIdx = serLocation(gs.pacmanLoc, outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}

// Serialize the location of the fruit (2 bytes)
func (gs *gameState) serFruit(outputBuf []byte, startIdx int) int {

	if gs.fruitExists {
		startIdx = serLocation(gs.fruitLoc, outputBuf, startIdx)
	} else {
		startIdx = serLocation(nullLoc, outputBuf, startIdx)
	}

	// Return the starting index of the next field
	return startIdx
}

// Serialize a ghost's information (3 bytes) - TODO: implement spawn offset
func (gs *gameState) serGhost(color int8, outputBuf []byte, startIdx int) int {

	// Add a flag at the 7th (highest) bit to indicate spawning
	var spawnFlag uint8 = 0
	g := gs.ghosts[color]
	if g.spawning {
		spawnFlag = 0x80
	}

	// Serialize the location information, followed by the fright cycles and spawn flag info
	startIdx = serLocation(g.loc, outputBuf, startIdx)
	startIdx = serUint8(g.frightCycles|spawnFlag, outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}

// Serialize a ghost's information (3 bytes)
func (gs *gameState) serGhosts(outputBuf []byte, startIdx int) int {

	// Serialize the ghost states, in the order (red -> pink -> cyan -> orange)
	startIdx = gs.serGhost(red, outputBuf, startIdx)
	startIdx = gs.serGhost(pink, outputBuf, startIdx)
	startIdx = gs.serGhost(cyan, outputBuf, startIdx)
	startIdx = gs.serGhost(orange, outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}
