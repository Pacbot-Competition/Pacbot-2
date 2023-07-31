package game

// Serialize the current number of ticks (2 bytes)
func (gs *gameState) serCurrTicks(outputBuf []byte, startIdx int) int {

	// Serialize the field, and return the new start index
	return serUint16(gs.currTicks, outputBuf, startIdx)
}

// Serialize the update ticks (1 byte)
func (gs *gameState) serUpdateTicks(outputBuf []byte, startIdx int) int {

	// Serialize the field, and return the new start index
	return serUint8(gs.updateTicks, outputBuf, startIdx)
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
