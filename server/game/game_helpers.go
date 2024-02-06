package game

import (
	"log"
	"slices"
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

	// Collect fruit, if applicable
	if gs.fruitExists() && gs.pacmanLoc.collidesWith(gs.fruitLoc) {
		gs.setFruitSteps(0)
		gs.incrementScore(fruitPoints)
	}

	// If there's no pellet, return
	if !gs.pelletAt(row, col) {
		return
	}

	// If we can clear the pellet's bit, decrease the number of pellets
	modifyBit(&(gs.pellets[row]), col, false)
	gs.decrementNumPellets()

	// If the we are in particular rows and columns, it is a super pellet
	superPellet := ((row == 3) || (row == 23)) && ((col == 1) || (col == 26))

	// Make all the ghosts frightened if a super pellet is collected
	if superPellet {
		gs.frightenAllGhosts()
	}

	// Update the score, depending on the pellet type
	if superPellet {
		gs.incrementScore(superPelletPoints)
	} else {
		gs.incrementScore(pelletPoints)
	}

	// Act depending on the number of pellets left over
	numPellets := gs.getNumPellets()

	// Spawn fruit, if applicable
	if (numPellets == fruitThreshold1) && !gs.fruitExists() {
		gs.setFruitSteps(fruitDuration)
	} else if (numPellets == fruitThreshold2) && !gs.fruitExists() {
		gs.setFruitSteps(fruitDuration)
	}

	// Other pellet-related events
	if numPellets == angerThreshold1 { // Ghosts get angry (speeding up)
		gs.setUpdatePeriod(uint8(max(1, int(gs.getUpdatePeriod())-2)))
		gs.setMode(chase)
		gs.setModeSteps(modeDurations[chase])
	} else if numPellets == angerThreshold2 { // Ghosts get angrier
		gs.setUpdatePeriod(uint8(max(1, int(gs.getUpdatePeriod())-2)))
		gs.setMode(chase)
		gs.setModeSteps(modeDurations[chase])
	} else if numPellets == 0 {
		gs.levelReset()
		gs.incrementLevel()
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

	// Keep track of how many ghosts need to respawn
	numGhostRespawns := 0

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
				numGhostRespawns++
			} else {
				gs.deathReset()
				return
			}
		}
	}

	// If no ghosts need to respawn, there's no more work to do
	if numGhostRespawns == 0 {
		return
	}

	// Lock the motion mutex to synchronize with other ghost update routines
	gs.respawnGhosts(numGhostRespawns, ghostRespawnFlag)
}

/***************************** Event-Based Resets *****************************/

// Reset the board (while leaving pellets alone) after Pacman dies
func (gs *gameState) deathReset() {

	// Acquire the Pacman control lock, to prevent other Pacman movement
	gs.muPacman.Lock()
	defer gs.muPacman.Unlock()

	// Set the game to be paused at the next update
	gs.setPauseOnUpdate(true)

	// Set Pacman to be in an empty state
	gs.pacmanLoc.copyFrom(emptyLoc)

	// Decrease the number of lives Pacman has left
	gs.decrementLives()

	/*
		If the mode is not the initial mode and the ghosts aren't angry,
		change the mode back to the initial mode
	*/
	if gs.getNumPellets() > angerThreshold1 {
		gs.setMode(initMode)
		gs.setModeSteps(modeDurations[initMode])
	}

	// Set the fruit steps back to 0
	gs.setFruitSteps(0)

	// Reset all the ghosts to their original locations
	gs.resetAllGhosts()
}

// Reset the board (including pellets) after Pacman clears a level
func (gs *gameState) levelReset() {

	// Set the game to be paused at the next update
	gs.setPauseOnUpdate(true)

	// Set Pacman to be in an empty state
	gs.pacmanLoc.copyFrom(emptyLoc)

	// If the mode is not the initial mode, change it
	gs.setMode(initMode)
	gs.setModeSteps(modeDurations[initMode])

	// Reset the level penalty
	gs.setLevelSteps(levelDuration)

	// Set the fruit steps back to 0
	gs.setFruitSteps(0)

	// Reset all the ghosts to their original locations
	gs.resetAllGhosts()

	// Reset the pellet bit array and count
	gs.resetPellets()
}

/************************** Motion (Pacman Location) **************************/

// Move Pacman one space in a given direction
func (gs *gameState) movePacmanDir(dir uint8) {

	// Acquire the Pacman control lock, to prevent other Pacman movement
	gs.muPacman.Lock()
	defer func() {

		// Unlock when we return
		gs.muPacman.Unlock()

		// Check collisions with all the ghosts
		gs.checkCollisions()
	}()

	// Ignore the command if the game is paused
	if gs.isPaused() || gs.getPauseOnUpdate() {
		return
	}

	// Shorthand to make computation simpler
	pLoc := gs.pacmanLoc

	// Calculate the next row and column
	nextRow, nextCol := pLoc.getNeighborCoords(dir)

	// Update Pacman's direction
	pLoc.updateDir(dir)

	// Check if there is a wall at the anticipated location, and return if so
	if gs.wallAt(nextRow, nextCol) {
		return
	}

	// Move Pacman the anticipated spot
	pLoc.updateCoords(nextRow, nextCol)
	gs.collectPellet(nextRow, nextCol)
}

// Move pacman to destination along shortest path (CV update)
func (gs *gameState) movePacmanAbsolute(newRow, newCol int8) {
	// Lock control over Pacman object and release at end
	gs.muPacman.Lock()
	defer gs.muPacman.Unlock()


	// Don't update position if we're paused
	if gs.isPaused() || gs.getPauseOnUpdate() {
		return
	}
	
	// Reject invalid coords
	if gs.wallAt(newRow, newCol) {
		return
	}

	pLoc := gs.pacmanLoc

	// Reject same coords
	if pLoc.row == newRow && pLoc.col == newCol {
		return
	}

	// Find likely path
	path := gs.findLikelyPath(newRow, newCol)

	// This really shouldn't happen but somehow the pathfinding has failed
	if path == nil {
		log.Println("\033[31mERR: Failed to find correct path\033[0m")
		return
	}

	// The new position is far from the old one, let's not traverse the path
	if len(path) > 5 {
		log.Println("\033[33mWARN: Interpolated path too long! "+
			"Tracking performance is likely degraded\033[0m")

		pLoc.updateCoords(newRow, newCol)
		gs.collectPellet(newRow, newCol)

		// Check collisions with ghosts
		gs.checkCollisions()
		return
	}

	// Move Pacman along the detected route
	for i := range path {
		nextPos := path[i]
		r, c := nextPos.r, nextPos.c
		pLoc.updateCoords(r, c)
		gs.collectPellet(r, c)

		// Check collisions with ghosts
		gs.checkCollisions()
	}
}

type pos struct {r, c int8}

func (p pos) getAdjacent() [4]pos {
	return [...]pos{
		{p.r + 1, p.c},
		{p.r, p.c + 1},
		{p.r - 1, p.c},
		{p.r, p.c - 1},
	}
}

// Find likely/shortest path to new coords
// precondition: lock pacman pos
func (gs *gameState) findLikelyPath(newRow, newCol int8) []pos {
	// Begin breadth-first search
	start := pos{gs.pacmanLoc.row, gs.pacmanLoc.col}
	queue := []pos{start}
	parent := make(map[pos]pos)
	parent[start] = pos{-1, -1}

	// Define target position
	target := pos{newRow, newCol}

	// Keep searching until we have exhausted all options or found it
	search_loop: for len(queue) != 0 {
		// Peek top, and remove
		curr := queue[0]
		queue = queue[1:]
		
		// Find adjacencies/neighbors of current cell
		neighbors := curr.getAdjacent()
		for i := range neighbors {
			adj := neighbors[i]

			// Already searched this one, continue
			if _, ok := parent[adj]; ok {
				continue
			}

			// Skip walls
			if gs.wallAt(adj.r, adj.c) {
				continue
			}

			// We can validly travel to the destination cell from curr
			parent[adj] = curr

			// We found the correct node
			if adj == target {
				break search_loop
			}

			// Search this later
			queue = append(queue, adj)
		}
	}

	// Backtrack the path
	path := []pos{target}
	for last, ok := target, true; ok && last != start; last, ok = parent[last] {
		path = append(path, last)
	}

	// Something has gone horribly wrong, last node doesnt backtrack to start
	if parent[path[len(path) - 1]] != start {
		return nil
	} 
	slices.Reverse(path)
	return path
}

// Move Pacman back to its spawn point, if necessary
func (gs *gameState) tryRespawnPacman() {
	// Acquire the Pacman control lock, to prevent other Pacman movement
	gs.muPacman.Lock()
	defer gs.muPacman.Unlock()

	// Set Pacman to be in its original state
	if gs.pacmanLoc.isEmpty() && gs.getLives() > 0 {
		gs.pacmanLoc.copyFrom(pacmanSpawnLoc)
	}
}

/******************************* Ghost Movement *******************************/

// Frighten all ghosts at once
func (gs *gameState) frightenAllGhosts() {

	// Acquire the ghost control lock, to prevent other ghost movement decisions
	gs.muGhosts.Lock()
	defer gs.muGhosts.Unlock()

	// Reset the ghost respawn combo back to 0
	gs.ghostCombo = 0

	// Loop over all the ghosts
	for _, ghost := range gs.ghosts {

		/*
			To frighten a ghost, set its fright steps to a specified value
			and trap it for one step (to force the direction to reverse)
		*/
		ghost.setFrightSteps(ghostFrightSteps)
		if !ghost.isTrapped() {
			ghost.setTrappedSteps(1)
		}
	}
}

// Reverse all ghosts at once (similar to frightenAllGhosts)
func (gs *gameState) reverseAllGhosts() {

	// Loop over all the ghosts
	for _, ghost := range gs.ghosts {

		/*
			To change the direction a ghost, trap it for one step
			(to force the direction to reverse)
		*/
		if !ghost.isTrapped() {
			ghost.setTrappedSteps(1)
		}
	}
}

// Reset all ghosts at once
func (gs *gameState) resetAllGhosts() {

	// Acquire the ghost control lock, to prevent other ghost movement
	gs.muGhosts.Lock()
	defer gs.muGhosts.Unlock()

	// Reset the ghost respawn combo back to 0
	gs.ghostCombo = 0

	// Add relevant ghosts to a wait group
	gs.wgGhosts.Add(int(numColors))

	// Reset each of the ghosts
	for _, ghost := range gs.ghosts {
		go ghost.reset()
	}

	// Wait for the resets to finish
	gs.wgGhosts.Wait()

	// If no lives are left, set all ghosts to stare at the player, menacingly
	if gs.getLives() == 0 {
		for _, ghost := range gs.ghosts {
			if ghost.color != orange {
				ghost.nextLoc.updateDir(none)
			} else { // Orange does like making eye contact, unfortunately
				ghost.nextLoc.updateDir(left)
			}
		}
	}
}

// Respawn some ghosts, according to a flag
func (gs *gameState) respawnGhosts(
	numGhostRespawns int, ghostRespawnFlag uint8) {

	// Acquire the ghost control lock, to prevent other ghost movement
	gs.muGhosts.Lock()
	defer gs.muGhosts.Unlock()

	// Add relevant ghosts to a wait group
	gs.wgGhosts.Add(numGhostRespawns)

	// Loop over the ghost colors again, to decide which should respawn
	for _, ghost := range gs.ghosts {

		// If the ghost should respawn, do so and increase the score and combo
		if getBit(ghostRespawnFlag, ghost.color) {

			// Respawn the ghost
			ghost.respawn()

			// Add points corresponding to the current combo length
			gs.incrementScore(comboMultiplier << uint16(gs.ghostCombo))

			// Increment the ghost respawn combo
			gs.ghostCombo++
		}
	}

	// Wait for the respawns to finish
	gs.wgGhosts.Wait()
}

// Update all ghosts at once
func (gs *gameState) updateAllGhosts() {

	// Acquire the ghost control lock, to prevent other ghost movement
	gs.muGhosts.Lock()
	defer gs.muGhosts.Unlock()

	// Add relevant ghosts to a wait group
	gs.wgGhosts.Add(int(numColors))

	// Loop over the individual ghosts
	for _, ghost := range gs.ghosts {
		go ghost.update()
	}

	// Wait for the respawns to finish
	gs.wgGhosts.Wait()
}

// A game state function to plan all ghosts at once
func (gs *gameState) planAllGhosts() {

	// Acquire the ghost control lock, to prevent other ghost movement
	gs.muGhosts.Lock()
	defer gs.muGhosts.Unlock()

	// Add pending ghost plans
	gs.wgGhosts.Add(int(numColors))

	// Plan each ghost's next move concurrently
	for _, ghost := range gs.ghosts {
		go ghost.plan()
	}

	// Wait until all pending ghost plans are complete
	gs.wgGhosts.Wait()
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
