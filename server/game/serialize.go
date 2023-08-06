package game

/*
IMPORTANT NOTE: All serializations are encoded in big-endian form
(most significant byte, MSB, first)
*/

/**************************** Integer Serialization ***************************/

/*
Get the byte at a particular index (0 = least significant byte,
1 = second least, etc.)
*/
func getByte[T uint8 | uint16 | uint32](num T, byteIdx int) byte {

	/*
		Uses bitwise operation magic (not really, look up how the >> and &
		operators work if you're interested)
	*/
	return byte((num >> (8 * byteIdx)) & 0xff)
}

// Serialize an individual byte (this should be simple, just getByte call)
func serUint8(num uint8, outputBuf []byte, startIdx int) int {
	outputBuf[startIdx] = getByte(num, 0)
	return startIdx + 1
}

// Serialize a uint16 (two getByte calls)
func serUint16(num uint16, outputBuf []byte, startIdx int) int {

	// Loop over each of the 2 bytes within the row (MSB first)
	for byteIdx := 1; byteIdx >= 0; byteIdx-- {

		// Serialize the byte
		outputBuf[startIdx] = getByte(num, byteIdx)

		// Add 1 to the start index, to prepare for serializing the next byte
		startIdx++
	}
	return startIdx
}

// Serialize a uint32 (four getByte calls)
func serUint32(num uint32, outputBuf []byte, startIdx int) int {

	// Loop over each of the 4 bytes within the row (MSB first)
	for byteIdx := 3; byteIdx >= 0; byteIdx-- {

		// Serialize the byte
		outputBuf[startIdx] = getByte(num, byteIdx)

		// Add 1 to the start index, to prepare for serializing the next byte
		startIdx++
	}

	// Return the starting index of the next field
	return startIdx
}

/***************************** Field Serialization ****************************/

// Serialize a location (no getByte calls, serialized manually)
func serLocation(loc *locationState, outputBuf []byte, startIdx int) int {

	// Lock the location state (so no writes can be done simultaneously)
	loc.RLock()
	defer loc.RUnlock()

	// Cover each coordinate of the location, one at a time
	outputBuf[startIdx+0] = byte((dRow[loc.dir] << 6) | loc.row)
	outputBuf[startIdx+1] = byte((dCol[loc.dir] << 6) | loc.col)

	// Return the starting index of the next field
	return startIdx + 2
}

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
	return serUint8(gs.mode, outputBuf, startIdx)
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
		startIdx = serLocation(emptyLoc, outputBuf, startIdx)
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

/***************************** State Serialization ****************************/

// Serialize all the information of the game state
func (gs *gameState) serFull(outputBuf []byte, startIdx int) int {

	// Packet header - contains the necessary information to render the ticker
	startIdx = gs.serCurrTicks(outputBuf, startIdx)
	startIdx = gs.serUpdatePeriod(outputBuf, startIdx)
	startIdx = gs.serGameMode(outputBuf, startIdx)

	// Ghosts - serializes the ghost states to the buffer
	startIdx = gs.serGhosts(outputBuf, startIdx)

	// Pacman - serializes the pacman location to the buffer
	startIdx = gs.serPacman(outputBuf, startIdx)

	// Fruit - serializes the fruit location (null if fruit doesn't exist)
	startIdx = gs.serFruit(outputBuf, startIdx)

	// Pellets - serializes the pellets to the buffer
	startIdx = gs.serPellets(outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}
