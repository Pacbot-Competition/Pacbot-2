package game

/***************************** Bitwise Operations *****************************/

/*
Get a bit within an unsigned integer (treating the integers
in pellets and walls as bit arrays)
*/
func getBit[T uint8 | uint16 | uint32](num T, bitIdx int8) bool {

	/*
		Uses bitwise operation magic (not really, look up how the >> and &
		operators work if you're interested)
	*/
	return bool(((num >> bitIdx) & 1) == 1)
}

/*
Get a bit within an unsigned integer (treating the integers in pellets
and walls as bit arrays)
*/
func modifyBit[T uint8 | uint16 | uint32](num *T, bitIdx int8, bitVal bool) {

	// If the bit is true, we should set the bit, otherwise we clear it
	if bitVal {
		*num |= (1 << bitIdx)
	} else {
		*num &= (^(1 << bitIdx))
	}
}

/****************************** Timing Functions ******************************/

/*
Determines if the game state is ready to update (i.e. we reached the start
of an update cycle, excluding the first cycle
*/
func (gs *gameState) updateReady() bool {
	return (gs.currTicks%uint16(gs.updatePeriod) == 0) && (gs.currTicks > 0)
}

/**************************** Positional Functions ****************************/

// Determines if a position is within the bounds of the maze
func (gs *gameState) inBounds(row int8, col int8) bool {
	return ((row >= 0 && row < mazeRows) && (col >= 0 && row < mazeCols))
}

// Determines if a pellet is at a given location
func (gs *gameState) pelletAt(row int8, col int8) bool {
	if !gs.inBounds(row, col) {
		return false
	}

	// Returns the bit of the pellet row corresponding to the column
	return getBit(gs.pellets[row], col)
}

// Determines if a super pellet is at a location (based on row or column)
func (gs *gameState) superPelletAt(row int8, col int8) bool {
	if !gs.pelletAt(row, col) {
		return false
	}

	// If the we are in particular rows and columns, the condition is met
	return ((row == 3) || (row == 23)) && ((col == 1) || (col == 26))
}

// Collects a pellet if it is at a given location
func (gs *gameState) collectPellet(row int8, col int8) {
	if !gs.inBounds(row, col) || !gs.pelletAt(row, col) {
		return
	}

	/* TODO: Update score */

	// Clears the bit of the pellet row
	modifyBit(&(gs.pellets[row]), col, false)
}

// Determines if a wall is at a given location
func (gs *gameState) wallAt(row int8, col int8) bool {
	if !gs.inBounds(row, col) {
		return false
	}

	// Returns the bit of the wall row corresponding to the column
	return getBit(gs.walls[row], col)
}
