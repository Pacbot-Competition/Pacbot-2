package game

import "fmt"

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

	// Get the current ticks value
	currTicks := gs.getCurrTicks()

	// Get the update period (uint16 to match the type of current ticks)
	updatePeriod := uint16(gs.getUpdatePeriod())

	// Update if the update period divides the current ticks
	return currTicks%updatePeriod == 0
}

/************************** General Helper Functions **************************/

// Helper function to frighten all the ghosts
func (gs *gameState) frightenGhosts() {
	for _, ghost := range gs.ghosts {
		ghost.frighten()
	}
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
		gs.incrementScore(superPelletPoints)
	} else {
		gs.incrementScore(pelletPoints)
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

/************************** Motion (Pacman Location) **************************/

// Move Pacman one space in a given direction
func (gs *gameState) movePacmanDir(dir uint8) {

	// Shorthand to make computation simpler
	pLoc := gs.pacmanLoc // Pacman location

	// Calculate the next row and column
	nextRow, nextCol := pLoc.getNeighborCoords(dir)

	// Update Pacman's direction
	pLoc.updateDir(dir)

	// Check if there is a wall at the anticipated location, and return if so
	if gs.wallAt(nextRow, nextCol) {
		return
	}

	// Move Pacman the anticipated spot
	pLoc.moveToCoords(nextRow, nextCol)
	pelletsLeft := gs.collectPellet(nextRow, nextCol)

	// Spawn fruit if 70 or 170 pellets are eaten
	if pelletsLeft == initPelletCount-70 {
		fmt.Println("Fruit should spawn")
	} else if pelletsLeft == initPelletCount-170 {
		fmt.Println("Fruit should spawn")
	}
}

/************************ Ghost Targeting (Chase Mode) ************************/

/*
Returns the chase location of the red ghost
(i.e. Pacman's exact location)
*/
func (gs *gameState) getChaseTargetRed() (int8, int8) {

	// Return Pacman's current location
	return gs.pacmanLoc.getCoords()
}

/*
Returns the chase location of the pink ghost
(i.e. 4 spaces ahead of Pacman's location)
*/
func (gs *gameState) getChaseTargetPink() (int8, int8) {

	// Return the red pink's target (4 spaces ahead of Pacman)
	return gs.pacmanLoc.getAheadCoords(4)
}

/*
Returns the chase location of the cyan ghost
(i.e. The red ghost's location, reflected about 2 spaces ahead of Pacman)
*/
func (gs *gameState) getChaseTargetCyan() (int8, int8) {

	// Get the 'pivot' square, 2 steps ahead of Pacman
	pivotRow, pivotCol := gs.pacmanLoc.getAheadCoords(2)

	// Get the current location of the red ghost
	redRow, redCol := gs.ghosts[red].loc.getCoords()

	// Return the pair of coordinates of the calculated target
	return (2*pivotRow - redRow),
		(2*pivotCol - redCol)
}

/*
Returns the chase location of the orange ghost
(i.e. Pacman's exact location, the same as red's target most of the time)
Though, if close enough to Pacman, it should choose its scatter target
*/
func (gs *gameState) getChaseTargetOrange() (int8, int8) {

	// Get Pacman's current location
	pacmanRow, pacmanCol := gs.pacmanLoc.getCoords()

	// Get the orange ghost's current location
	orangeRow, orangeCol := gs.ghosts[orange].loc.getCoords()

	// Calculate the squared distance to Pacman's location
	distSq := gs.distSq(orangeRow, orangeCol, pacmanRow, pacmanCol)

	// If Pacman is far enough, return Pacman's location
	if distSq >= 64 {
		return (pacmanRow),
			(pacmanCol)
	}

	// Otherwise, return the scatter location of orange
	return gs.ghosts[orange].scatterTarget.getCoords()
}
