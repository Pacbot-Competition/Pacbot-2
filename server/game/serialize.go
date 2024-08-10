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

	// (Read) lock the location state
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

	// Serialize and return the starting index of the next field
	return serUint16(gs.getCurrTicks(), outputBuf, startIdx)
}

// Serialize the update period (1 byte)
func (gs *gameState) serUpdatePeriod(outputBuf []byte, startIdx int) int {

	// Serialize and return the starting index of the next field
	return serUint8(gs.getUpdatePeriod(), outputBuf, startIdx)
}

// Serialize the game mode (1 byte)
func (gs *gameState) serGameMode(outputBuf []byte, startIdx int) int {

	// Serialize and return the starting index of the next field
	return serUint8(gs.getMode(), outputBuf, startIdx)
}

/*
Serialize the number of steps until the mode changes, in addition to the
duration of the mode in steps (2 bytes)
*/
func (gs *gameState) serModeSteps(outputBuf []byte, startIdx int) int {

	// Serialize the number of mode steps
	startIdx = serUint8(gs.getModeSteps(), outputBuf, startIdx)

	// Serialize the duration of this (last unpaused) mode
	modeDuration := modeDurations[gs.getLastUnpausedMode()]
	startIdx = serUint8(modeDuration, outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}

/*
Serialize number of steps (update periods) before a speedup penalty starts (1 byte)
*/
func (gs *gameState) serLevelSteps(outputBuf []byte, startIdx int) int {

	// Serialize the level steps and return the starting index of the next field
	return serUint16(gs.getLevelSteps(), outputBuf, startIdx)
}

// Serialize the current score (2 bytes)
func (gs *gameState) serCurrScore(outputBuf []byte, startIdx int) int {

	// Serialize and return the starting index of the next field
	return serUint16(gs.getScore(), outputBuf, startIdx)
}

// Serialize the current level (1 byte)
func (gs *gameState) serCurrLevel(outputBuf []byte, startIdx int) int {

	// Serialize and return the starting index of the next field
	return serUint8(gs.getLevel(), outputBuf, startIdx)
}

// Serialize the current lives (1 byte)
func (gs *gameState) serCurrLives(outputBuf []byte, startIdx int) int {

	// Serialize and return the starting index of the next field
	return serUint8(gs.getLives(), outputBuf, startIdx)
}

// Serialize the ghost combo (1 byte)
func (gs *gameState) serGhostCombo(outputBuf []byte, startIdx int) int {

	// Serialize and return the starting index of the next field
	return serUint8(gs.ghostCombo, outputBuf, startIdx)
}

// Serialize the pellets (4 * mazeRows bytes)
func (gs *gameState) serPellets(outputBuf []byte, startIdx int) int {

	// (Read) lock the pellets array
	gs.muPellets.RLock()
	defer gs.muPellets.RUnlock()

	// Loop over each row
	for row := int8(0); row < mazeRows; row++ {

		// Serialize each row from uint32 to 4 bytes
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

	// (Read) lock the fruit state
	gs.muFruit.RLock()
	{
		if gs.fruitExists() { // Serialize the fruit's location if it exists
			startIdx = serLocation(gs.fruitLoc, outputBuf, startIdx)
		} else { // Otherwise, give an empty (0x00 0x00) location
			startIdx = serLocation(emptyLoc, outputBuf, startIdx)
		}
	}
	gs.muFruit.RUnlock()

	// Serialize the number of steps the fruit has been spawned
	fruitSteps := gs.getFruitSteps()
	startIdx = serUint8(fruitSteps, outputBuf, startIdx)

	// Serialize the duration of the fruit
	fruitDuration := fruitDuration
	startIdx = serUint8(fruitDuration, outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}

// Serialize a ghost's information (3 bytes)
func (gs *gameState) serGhost(color uint8, outputBuf []byte, startIdx int) int {

	// Retrieve this ghost's struct
	g := gs.ghosts[color]

	// Serialize the location information first
	startIdx = serLocation(g.loc, outputBuf, startIdx)

	// Lock the ghost's other state variables
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Add a flag at the 7th (highest) bit to indicate spawning
	var spawnFlag uint8 = 0
	if g.spawning {
		spawnFlag = 0b10000000
	}

	// Serialize the fright steps and spawn flag info next
	startIdx = serUint8(g.frightSteps|spawnFlag, outputBuf, startIdx)

	// Add a flag at the 7th (highest) bit to indicate eaten
	var eatenFlag uint8 = 0
	if g.eaten {
		eatenFlag = 0b10000000
	}

	// Serialize the trapped steps and eaten flag info next
	startIdx = serUint8(g.trappedSteps|eatenFlag, outputBuf, startIdx)

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
	startIdx = gs.serModeSteps(outputBuf, startIdx)
	startIdx = gs.serLevelSteps(outputBuf, startIdx)

	// General game state information
	startIdx = gs.serCurrScore(outputBuf, startIdx)
	startIdx = gs.serCurrLevel(outputBuf, startIdx)
	startIdx = gs.serCurrLives(outputBuf, startIdx)
	startIdx = gs.serGhostCombo(outputBuf, startIdx)

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
