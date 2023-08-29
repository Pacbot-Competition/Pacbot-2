package game

import (
	"log"
)

/******************************** Ghost Resets ********************************/

// Respawn the ghost
func (g *ghostState) reset() {

	// Lock the motion mutex to synchronize with other motion routines
	g.muMotion.Lock()
	defer g.muMotion.Unlock()

	// Mark this operation as done once we return
	defer g.game.wgGhosts.Done()

	// Set the ghost to be eaten and spawning
	g.setSpawning(true)
	g.setTrappedCycles(ghostTrappedCycles[g.color])

	// Set the current ghost to be at an empty location
	g.loc.copyFrom(emptyLoc)

	/*
		Set the current location of the ghost to be its spawn point
		(or pink's spawn location, in the case of red, so it spawns in the box)
	*/
	g.nextLoc.copyFrom(ghostSpawnLocs[g.color])
}

/****************************** Ghost Respawning ******************************/

// Respawn the ghost
func (g *ghostState) respawn() {

	// Lock the motion mutex to synchronize with other motion routines
	g.muMotion.Lock()
	defer g.muMotion.Unlock()

	// Mark this operation as done once we return
	defer g.game.wgGhosts.Done()

	// Set the ghost to be eaten and spawning
	g.setSpawning(true)
	g.setEaten(true)

	// Set the current ghost to be at an empty location
	g.loc.copyFrom(emptyLoc)

	/*
		Set the current location of the ghost to be its spawn point
		(or pink's spawn location, in the case of red, so it spawns in the box)
	*/
	if g.color == red {
		g.nextLoc.copyFrom(ghostSpawnLocs[pink])
	} else {
		g.nextLoc.copyFrom(ghostSpawnLocs[g.color])
	}
}

/******************** Ghost Updates (before serialization) ********************/

// Update the ghost's position
func (g *ghostState) update() {

	/*
		If the ghost is at the red spawn point and not moving downwards,
		we can mark it as done spawning
	*/
	if g.loc.collidesWith(ghostSpawnLocs[red]) && g.loc.getDir() != down {
		g.setSpawning(false)
	}

	// Set the ghost to be no longer eaten, if applicable
	if g.isEaten() {
		g.setEaten(false)
		g.setFrightCycles(0)
	}

	// Decrement the ghost's fright cycles count if necessary
	if g.isFrightened() {
		g.decFrightCycles()
	}

	// Copy the next location into the current location
	g.loc.copyFrom(g.nextLoc)
}

/******************** Ghost Planning (after serialization) ********************/

// Plan the ghost's next move
func (g *ghostState) plan() {

	// Lock the motion mutex to synchronize with other motion routines
	g.muMotion.Lock()
	defer g.muMotion.Unlock()

	// Mark the plan as done
	defer g.game.wgGhosts.Done()

	// If the location is empty (i.e. after a reset/respawn), don't plan
	if g.loc.isEmpty() {
		return
	}

	// Determine the next position based on the current direction
	g.nextLoc.advanceFrom(g.loc)

	// If the ghost is trapped, reverse the current direction and return
	if g.isTrapped() {
		g.nextLoc.reverseDir()
		g.decTrappedCycles()
		return
	}

	// Keep local copies of the fright cycles and spawning variables
	frightCycles := g.getFrightCycles()
	spawning := g.isSpawning()

	// Decide on a target for this ghost, depending on the game mode
	var targetRow, targetCol int8

	// Capture the last unpaused current game mode (could be the current mode)
	mode := g.game.getLastUnpausedMode()

	/*
		If the ghost is spawning in the ghost house, choose red's spawn
		location as the target to encourage it to leave the ghost house

		Otherwise: pick chase or scatter targets, depending on the mode
	*/
	if spawning && !g.loc.collidesWith(ghostSpawnLocs[red]) &&
		!g.nextLoc.collidesWith(ghostSpawnLocs[red]) {
		targetRow, targetCol = ghostSpawnLocs[red].getCoords()
	} else if mode == chase { // Chase mode targets
		targetRow, targetCol = g.game.getChaseTarget(g.color)
	} else if mode == scatter { // Scatter mode targets
		targetRow, targetCol = g.scatterTarget.getCoords()
	}

	/*
		Determine whether each of the four neighboring moves to the next
		location is valid, and count how many are good
	*/
	numValidMoves := 0
	var moveValid [4]bool
	var moveDistSq [4]int
	for dir := uint8(0); dir < 4; dir++ {

		// Get the neighboring cell in that location
		row, col := g.nextLoc.getNeighborCoords(dir)

		// Calculate the distance from the target to the move location
		moveDistSq[dir] = g.game.distSq(row, col, targetRow, targetCol)

		// Determine if that move is valid
		moveValid[dir] = !g.game.wallAt(row, col)

		// Considerations when the ghost is spawning
		if spawning {

			// Determine if the move would be within the ghost house
			if g.game.ghostSpawnAt(row, col) {
				moveValid[dir] = true
			}

			/*
				Determine if the move would help the ghost escape the ghost house,
				and make it a valid one if so
			*/
			if row == ghostHouseExitRow && col == ghostHouseExitCol {
				moveValid[dir] = true
			}
		}

		// If this move would make the ghost reverse, skip it
		if dir == g.nextLoc.getReversedDir() {
			moveValid[dir] = false
		}

		// Increment the valid moves counter if necessary
		if moveValid[dir] {
			numValidMoves++
		}
	}

	// Debug statement, in case a ghost somehow is surrounded by all walls
	if numValidMoves == 0 {
		row, col := g.nextLoc.getCoords()
		log.Printf("\033[35mWARN: %s has nowhere to go "+
			"(row = %d, col = %d, spawning = %t)\n\033[0m",
			ghostNames[g.color], row, col, spawning)
		return
	}

	/*
		 	If the ghost will still frightened one tick later, immediately choose
			a random valid direction and return
	*/
	if frightCycles > 1 {

		// Generate a random index out of the valid moves
		randomNum := g.game.rng.Intn(numValidMoves)

		// Loop over all directions
		for dir, count := uint8(0), 0; dir < 4; dir++ {

			// Skip any invalid moves
			if !moveValid[dir] {
				continue
			}

			// If we have reached the correct move, update the direction and return
			if count == randomNum {
				g.nextLoc.updateDir(dir)
				return
			}

			// Update the count of valid moves so far
			count++
		}
	}

	// Otherwise, choose the best direction to reach the target
	bestDir := up
	bestDist := 0xffffffff // Some arbitrarily high number
	for dir := uint8(0); dir < 4; dir++ {

		// Skip any invalid moves
		if !moveValid[dir] {
			continue
		}

		// Compare this direction to the best so far
		if moveDistSq[dir] < bestDist {
			bestDir = dir
			bestDist = moveDistSq[dir]
		}
	}

	// Once we have picked the best direction, update it
	g.nextLoc.updateDir(bestDir)
}
