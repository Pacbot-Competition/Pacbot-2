package game

import (
	"log"
)

/***************************** Bitwise Operations *****************************/

/*
Get a bit within an unsigned integer (treating the integers
in pellets and walls as bit arrays)
*/
func getBit[N uint8 | uint16 | uint32, I int8 | uint8](
	num N, bitIdx I) bool {

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
func modifyBit[N uint8 | uint16 | uint32, I int8 | uint8](
	num *N, bitIdx I, bitVal bool) {

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

	// Loop over all the ghosts
	for _, ghost := range gs.ghosts {

		/*
			To frighten a ghost, set its fright cycles to a specified value
			and trap it for one cycle (to force the direction to reverse)
		*/
		ghost.setFrightCycles(ghostFrightCycles)
		if !ghost.isTrapped() {
			ghost.setTrappedCycles(1)
		}
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
func (gs *gameState) collectPellet(row int8, col int8) {
	if !gs.pelletAt(row, col) {
		return
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
	{
		// Clear the pellet's bit and decrement the number of pellets
		modifyBit(&(gs.pellets[row]), col, false)
		gs.numPellets--
	}
	gs.muPellets.Unlock()

	// Act depending on the number of pellets left over
	gs.muPellets.RLock()
	defer gs.muPellets.RUnlock()

	// Spawn fruit, if applicable
	gs.muFruit.Lock()
	{
		if gs.numPellets == fruitThreshold1 && !gs.fruitSpawned1 {
			log.Println("Fruit 1 should spawn")
			gs.fruitSpawned1 = true
		} else if gs.numPellets == fruitThreshold2 && !gs.fruitSpawned2 {
			log.Println("Fruit 2 should spawn")
			gs.fruitSpawned2 = true
		}
	}
	gs.muFruit.Unlock()

	// Other pellet-related events
	if gs.numPellets == angerThreshold1 { // Ghosts get angry (speeding up)
		gs.setUpdatePeriod(gs.getUpdatePeriod() - 2)
	} else if gs.numPellets == angerThreshold2 { // Ghosts get angrier
		gs.setUpdatePeriod(gs.getUpdatePeriod() - 2)
	} else if gs.numPellets == 0 {
		log.Println("Pacman won!")
	}
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

/***************************** Collision Handling *****************************/

// Check collisions between Pacman and all the ghosts
func (gs *gameState) checkCollisions() {

	// Flag to decide which ghosts should respawn
	var ghostRespawnFlag uint8 = 0

	// Loop over all the ghosts
	for _, ghost := range gs.ghosts {

		// Check each collision individually
		if gs.pacmanLoc.collidesWith(ghost.loc) {

			// If the ghost was already eaten, skip it
			if ghost.isEaten() {
				continue
			}

			// If the ghost is frightened, Pacman eats it, otherwise Pacman dies
			if ghost.isFrightened() {
				modifyBit(&ghostRespawnFlag, ghost.color, true)
			} else {
				gs.deathReset()
				return
			}
		}
	}

	// Loop over the ghost colors again, to decide which should respawn
	for _, ghost := range gs.ghosts {

		// If the ghost should respawn, call its respawn function
		if getBit(ghostRespawnFlag, ghost.color) {
			ghost.respawn()
		}
	}
}

// Reset the board (while leaving pellets alone) after Pacman dies
func (gs *gameState) deathReset() {
	log.Println("Pacman Lost!")

	// Add 4 ghosts to the ghost plans wait group, to halt updates
	gs.wgGhosts.Add(4)

	// Reset each of the 4 ghosts
	for _, ghost := range gs.ghosts {
		go ghost.reset()
	}

	// Pause the game
	gs.setPauseOnUpdate(true)
}

/************************** Motion (Pacman Location) **************************/

// Move Pacman one space in a given direction
func (gs *gameState) movePacmanDir(dir uint8) {

	// Shorthand to make computation simpler
	pLoc := gs.pacmanLoc

	// Calculate the next row and column
	nextRow, nextCol := pLoc.getNeighborCoords(dir)

	// Update Pacman's direction
	pLoc.updateDir(dir)

	// Check collisions with all the ghosts
	gs.checkCollisions()

	// Check if there is a wall at the anticipated location, and return if so
	if gs.wallAt(nextRow, nextCol) {
		return
	}

	// Move Pacman the anticipated spot
	pLoc.moveToCoords(nextRow, nextCol)
	gs.collectPellet(nextRow, nextCol)
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

	// If Pacman is far enough from the ghost, return Pacman's location
	if gs.distSq(orangeRow, orangeCol, pacmanRow, pacmanCol) >= 64 {
		return (pacmanRow),
			(pacmanCol)
	}

	// Otherwise, return the scatter location of orange
	return gs.ghosts[orange].scatterTarget.getCoords()
}

// Returns the chase location of an arbitrary ghost color
func (gs *gameState) getChaseTarget(color uint8) (int8, int8) {
	switch color {
	case red:
		return gs.getChaseTargetRed()
	case pink:
		return gs.getChaseTargetPink()
	case cyan:
		return gs.getChaseTargetCyan()
	case orange:
		return gs.getChaseTargetOrange()
	}
	return emptyLoc.getCoords()
}
