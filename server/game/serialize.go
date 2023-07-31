package game

/*
	NOTE: All serializations are in big-endian form
				(most significant byte, MSB, first)
*/

// Get the byte at a particular index (0 = least significant, 1 = second least, etc.)
func getByte[T uint8 | uint16 | uint32](num T, byteIdx int) byte {
	return byte((num >> (8 * byteIdx)) & 0xff)
}

// Serialize the current number of ticks (2 bytes)
func (gs *gameState) serCurrTicks(outputBuf []byte, startIdx int) int {

	// Serialize each of the bytes separately (MSB first)
	outputBuf[startIdx] = getByte(gs.currTicks, 1)
	outputBuf[startIdx+1] = getByte(gs.currTicks, 0)

	// Return the starting index of the next field
	return startIdx + 2
}

// Serialize the update ticks (1 byte)
func (gs *gameState) serUpdateTicks(outputBuf []byte, startIdx int) int {

	// Serialize as a single byte
	outputBuf[startIdx+1] = getByte(gs.updateTicks, 0)

	// Return the starting index of the next field
	return startIdx + 1
}

// Serialize the game mode (1 byte)
func (gs *gameState) serGameMode(outputBuf []byte, startIdx int) int {

	// Serialize each of the bytes separately
	outputBuf[startIdx+1] = getByte(gs.gameMode, 0)

	// Return the starting index of the next field
	return startIdx + 1
}

// Serialize the pellets (4 * mazeRows bytes)
func (gs *gameState) serPellets(outputBuf []byte, startIdx int) int {

	// Loop over each row
	for row := 0; row < mazeRows; row++ {

		// Loop over each of the 4 bytes within the row (MSB first)
		for byteIdx := 3; byteIdx >= 0; byteIdx-- {

			// Serialize the byte
			outputBuf[startIdx] = getByte(gs.pellets[row], byteIdx)

			// Add 1 to the start index, to prepare for serializing the next byte
			startIdx++
		}
	}

	// Return the starting index of the next field
	return startIdx
}
