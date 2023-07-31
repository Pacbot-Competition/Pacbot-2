package game

// Determines if a position is within the bounds of the maze
func (gs *gameState) inBounds(row int, col int) bool {
	return ((row >= 0 && row < mazeRows) && (col >= 0 && row < mazeCols))
}

// Determines if a pellet is at a given location
func (gs *gameState) pelletAt(row int, col int) bool {
	if !gs.inBounds(row, col) {
		return false
	}
	return ((gs.pellets[row] >> col) & 1) == 1
}

// Determines if a wall is at a given location
func (gs *gameState) wallAt(row int, col int) bool {
	if !gs.inBounds(row, col) {
		return false
	}
	return ((gs.walls[row] >> col) & 1) == 1
}
