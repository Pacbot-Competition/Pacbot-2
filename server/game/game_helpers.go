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

// Determines if the game state is ready to update
func (gs *gameState) updateReady() bool {
	return (gs.currTicks%uint16(gs.updatePeriod) == 0)
}

/************************** General Helper Functions **************************/

// Helper function to frighten all the ghosts
func (gs *gameState) frightenGhosts() {
	for color := uint8(0); color < 4; color++ {
		gs.ghosts[color].frighten()
	}
}

// Helper function to increment the current score of the game
func (gs *gameState) incrementScore(change uint16) {

	// (Write) lock the current score
	gs.muScore.Lock()
	defer gs.muScore.Unlock()

	// Add the change to the score
	gs.currScore += change
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

	// (Read) lock the pellets array
	gs.muPellets.RLock()
	defer gs.muPellets.RUnlock()

	// Returns the bit of the pellet row corresponding to the column
	return getBit(gs.pellets[row], col)
}

/*
Collects a pellet if it is at a given location
Returns the number of pellets that are left
*/
func (gs *gameState) collectPellet(row int8, col int8) uint16 {
	if !gs.pelletAt(row, col) {

		// (Read) lock the number of pellets, then return
		gs.muPellets.RLock()
		defer gs.muPellets.RUnlock()
		return gs.numPellets
	}

	// If the we are in particular rows and columns, it is a super pellet
	superPellet := ((row == 3) || (row == 23)) && ((col == 1) || (col == 26))

	// Make all the ghosts frightened if a super pellet is collected
	if superPellet {
		gs.frightenGhosts()
	}

	// Update the score, depending on the pellet type
	if superPellet {
		gs.currScore += 50 // Super pellet = 50 pts
	} else {
		gs.currScore += 10 // Normal pellet = 10 pts
	}

	// (Write) lock the pellets array, then clear the pellet's bit
	gs.muPellets.Lock()
	defer gs.muPellets.Unlock()

	// Clear the pellet's bit and decrement the number of pellets
	modifyBit(&(gs.pellets[row]), col, false)
	gs.numPellets--

	return gs.numPellets
}

// Determines if a wall is at a given location
func (gs *gameState) wallAt(row int8, col int8) bool {
	if !gs.inBounds(row, col) {
		return true
	}

	// Returns the bit of the wall row corresponding to the column
	return getBit(gs.walls[row], col)
}

// Determines if the ghost house is at a given location
func (gs *gameState) ghostSpawnAt(row int8, col int8) bool {
	if !gs.inBounds(row, col) {
		return false
	}

	// Returns the bit of the wall row corresponding to the column
	return ((row >= 13) && (row <= 14)) && ((col >= 11) && (col <= 15))
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

	// (Read) lock the pacman location and red ghost location
	gs.pacmanLoc.RLock()
	gs.ghosts[red].loc.RLock()
	defer func() {
		gs.pacmanLoc.RUnlock()
		gs.ghosts[red].loc.RUnlock()
	}()

	// Shorthands to make computation simpler
	pLoc := gs.pacmanLoc       // Pacman location
	rLoc := gs.ghosts[red].loc // Red location

	// Calculate the position of the 'pivot' square (2 ahead of Pacman)
	pivotRow := pLoc.row + 2*dRow[pLoc.dir]
	pivotCol := pLoc.col + 2*dCol[pLoc.dir]

	// Return the pair of coordinates of the calculated target
	return (2*pivotRow - rLoc.row),
		(2*pivotCol - rLoc.col)
}

/*
Returns the chase location of the orange ghost
(i.e. Pacman's exact location, the same as red's target most of the time)
Though, if close enough to Pacman, it should choose its scatter target
*/
func (gs *gameState) getChaseTargetOrange() (int8, int8) {

	// (Read) lock the pacman location and orange ghost location
	gs.pacmanLoc.RLock()
	gs.ghosts[orange].loc.RLock()
	defer func() {
		gs.pacmanLoc.RUnlock()
		gs.ghosts[orange].loc.RUnlock()
	}()

	// Shorthands to make computation simpler
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
