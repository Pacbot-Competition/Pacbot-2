package game

// Determines if a position is within the bounds of the maze
func (gs *gameState) inBounds(row int8, col int8) bool {
	return ((row >= 0 && row < mazeRows) && (col >= 0 && row < mazeCols))
}

// Determines if a pellet is at a given location
func (gs *gameState) pelletAt(row int8, col int8) bool {
	if !gs.inBounds(row, col) {
		return false
	}

	// Returns the bit of the pellet row, by shifting right and checking the lowest bit
	return ((gs.pellets[row] >> col) & 1) == 1
}

// Determines if a pellet is at a given location
func (gs *gameState) collectPellet(row int8, col int8) {
	if !gs.inBounds(row, col) {
		return
	}

	/* TODO: Update score */

	// Clears the bit of the pellet row, by AND-ing with the complement of the bit
	gs.pellets[row] &= (^(1 << col))
}

// Determines if a wall is at a given location
func (gs *gameState) wallAt(row int8, col int8) bool {
	if !gs.inBounds(row, col) {
		return false
	}
	return ((gs.walls[row] >> col) & 1) == 1
}
