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
Determines if the game state is ready to update
*/
func (gs *gameState) updateReady() bool {
	return (gs.currTicks%uint16(gs.updatePeriod) == 0)
}

/**************************** Positional Functions ****************************/

// Determines if a position is within the bounds of the maze
func (gs *gameState) inBounds(row int8, col int8) bool {
	return ((row >= 0 && row < mazeRows) && (col >= 0 && col < mazeCols))
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
		return true
	}

	// Returns the bit of the wall row corresponding to the column
	return getBit(gs.walls[row], col)
}

// Calculates the squared Euclidean distance between two points
func (gs *gameState) distSq(row1, col1, row2, col2 int8) int {
	dx := int(row2 - row1)
	dy := int(col2 - col1)
	return dx*dx + dy*dy
}

/************************ Ghost Targeting (Chase Mode) ************************/

/*
Returns the chase location of the red ghost
(i.e. Pacman's exact location)
*/
func (gs *gameState) getChaseTargetRed() (int8, int8) {

	// (Read) lock the pacman location
	gs.pacmanLoc.RLock()
	defer gs.pacmanLoc.RUnlock()

	// Return the pacman location, as a row and column
	return (gs.pacmanLoc.row),
		(pacmanSpawnLoc.col)
}

/*
Returns the chase location of the pink ghost
(i.e. 4 spaces ahead of Pacman's location)
*/
func (gs *gameState) getChaseTargetPink() (int8, int8) {

	// (Read) lock the pacman location
	gs.pacmanLoc.RLock()
	defer gs.pacmanLoc.RUnlock()

	// Return the red ghost's target
	return (gs.pacmanLoc.row + 4*dRow[gs.pacmanLoc.dir]),
		(gs.pacmanLoc.col + 4*dCol[gs.pacmanLoc.dir])
}

/*
Returns the chase location of the cyan ghost
(i.e. The red ghost's location, reflected about 2 spaces ahead of Pacman)
*/
func (gs *gameState) getChaseTargetCyan() (int8, int8) {

	// (Read) lock the pacman location
	gs.pacmanLoc.RLock()
	defer gs.pacmanLoc.RUnlock()

	// Calculate the position of the 'pivot' square (2 ahead of Pacman)
	pivotRow := gs.pacmanLoc.row + 2*dRow[gs.pacmanLoc.dir]
	pivotCol := gs.pacmanLoc.col + 2*dCol[gs.pacmanLoc.dir]

	// (Read) lock the red ghost's location
	gs.ghosts[red].loc.RLock()
	defer gs.ghosts[red].loc.RUnlock()

	// Return the pair of coordinates of the calculated target
	return (2*pivotRow - gs.ghosts[red].loc.row),
		(2*pivotCol - gs.ghosts[red].loc.col)
}

/*
Returns the chase location of the orange ghost
(i.e. Pacman's exact location, the same as red's target most of the time)
Though, if close enough to Pacman, it should choose its scatter target
*/
func (gs *gameState) getChaseTargetOrange() (int8, int8) {

	// (Read) lock the pacman location
	gs.pacmanLoc.RLock()
	gs.ghosts[orange].loc.RLock()
	defer func() {
		gs.pacmanLoc.RUnlock()
		gs.ghosts[orange].loc.RUnlock()
	}()

	// Shorthands to make return values simpler easier
	pLoc := gs.pacmanLoc                        // Pacman location
	oLoc := gs.ghosts[orange].loc               // Orange location
	oScatter := gs.ghosts[orange].scatterTarget // Orange scatter target

	distSq := gs.distSq(oLoc.row, oLoc.col, pLoc.row, pLoc.col)

	// If Pacman is far enough, return Pacman's location
	if distSq >= 64 {
		return (pLoc.row),
			(pLoc.col)
	}

	// Otherwise, return the scatter location of orange
	return (oScatter.row),
		(oScatter.col)
}
