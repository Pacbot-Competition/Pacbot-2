package game

import (
	"fmt"
	"sync"
)

// Enum-like declaration to hold the ghost colors
const (
	red    uint8 = 0
	pink   uint8 = 1
	cyan   uint8 = 2
	orange uint8 = 3
)

// Names of the ghosts (not the nicknames, just the colors, for debugging)
var ghostNames [4]string = [...]string{
	"red",
	"pink",
	"cyan",
	"orange",
}

/*
An object to keep track of the location and attributes of a ghost
*/
type ghostState struct {
	loc           *locationState // Current location
	nextLoc       *locationState // Planned location (for next update)
	scatterTarget *locationState // Position of Ifixed) scatter target
	game          *gameState     // The game state tied to the ghost
	color         uint8
	trappedCycles uint8
	frightCycles  uint8
	spawning      bool         // Flag set when spawning
	wasEaten      bool         // Flag set when eaten and returning to ghost house
	muState       sync.RWMutex // Mutex to lock general state parameters
}

// Create a new ghost state with given location and color values
func newGhostState(_gameState *gameState, _color uint8) *ghostState {
	return &ghostState{
		loc:           newLocationStateCopy(emptyLoc),
		nextLoc:       newLocationStateCopy(ghostSpawnLocs[_color]),
		scatterTarget: newLocationStateCopy(ghostScatterTargets[_color]),
		game:          _gameState,
		color:         _color,
		trappedCycles: ghostTrappedCycles[_color],
		frightCycles:  0,
		spawning:      true,
		wasEaten:      false,
	}
}

// Frighten the ghost
func (g *ghostState) frighten() {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.frightCycles = ghostFrightCycles
	}
	g.muState.Unlock()
}

// Update the ghost's position: copy the location from the next location state
func (g *ghostState) update() {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		// If we were just at the red spawn point, the ghost is done spawning
		if g.loc.collidesWith(ghostSpawnLocs[red]) {
			g.spawning = false
		}

		// Decrement the ghost's fright cycles count if necessary
		if g.frightCycles > 0 {
			g.frightCycles--
		}
	}
	g.muState.Unlock()

	// Copy the next location into the current location
	g.loc.copyFrom(g.nextLoc)
}

func (g *ghostState) plan(wg *sync.WaitGroup) {

	// Return that this go-routine has completed
	defer wg.Done()

	// Determine the next position based on the current direction
	g.nextLoc.advanceFrom(g.loc)

	// If the ghost is trapped, reverse the current direction and return
	if g.trappedCycles > 0 {
		g.nextLoc.reverseDir()
		g.trappedCycles-- // Decrement the trapped cycles counter (no lock needed)
		return
	}

	// Keep local copies of the fright cycles and spawning variables
	var spawning bool
	var frightCycles uint8
	g.muState.Lock()
	{
		spawning = g.spawning   // Copy the spawning flag
		if g.frightCycles > 0 { // Decrement the fright cycles
			g.frightCycles--
		}
		frightCycles = g.frightCycles // Copy the fright cycles
	}
	g.muState.Unlock()

	// Decide on a target for this ghost, depending on the game mode
	var targetRow, targetCol int8

	/*
		If the ghost is spawning in the ghost house, choose red's spawn
		location as the target to encourage it to leave the ghost house
	*/
	if spawning && !g.loc.collidesWith(ghostSpawnLocs[red]) &&
		!g.nextLoc.collidesWith(ghostSpawnLocs[red]) {
		targetRow, targetCol = ghostSpawnLocs[red].row, ghostSpawnLocs[red].col
	} else if g.game.mode == chase { // Chase mode targets
		switch g.color {
		case red:
			targetRow, targetCol = g.game.getChaseTargetRed()
		case pink:
			targetRow, targetCol = g.game.getChaseTargetPink()
		case cyan:
			targetRow, targetCol = g.game.getChaseTargetCyan()
		case orange:
			targetRow, targetCol = g.game.getChaseTargetOrange()
		}
	} else if g.game.mode == scatter { // Scatter mode targets
		targetRow, targetCol = g.scatterTarget.row, g.scatterTarget.col
	}

	/*
		Determine whether each of the four neighboring moves to the next
		location is valid, and count how many are good
	*/
	numValidMoves := 0
	var moveValid [4]bool
	var moveDistSq [4]int
	for dir := uint8(0); dir < 4; dir++ {

		// If this move would make the ghost reverse, skip it
		if dir == g.nextLoc.getReversedDir() {
			continue
		}

		// Get the neighboring cell in that location
		row, col := g.nextLoc.getNeighborCoords(dir)

		// Calculate the distance from the target to the move location
		moveDistSq[dir] = g.game.distSq(row, col, targetRow, targetCol)

		// Determine if that move is valid
		moveValid[dir] = !g.game.wallAt(row, col)

		/*
			Determine if the move would help the ghost escape the ghost house,
			and make it a valid one if so
		*/
		if spawning && row == ghostHouseExitRow && col == ghostHouseExitCol {
			moveValid[dir] = true
		}

		// Increment the valid moves counter if necessary
		if moveValid[dir] {
			numValidMoves++
		}
	}

	// Debug statement, in case a ghost somehow is surrounded by all walls
	if numValidMoves == 0 {
		fmt.Printf("\033[35mWARN: Ghost #%d (%s) has nowhere to go\n\033[0m", g.color, ghostNames[g.color])
		return
	}

	// If frightened, immediately choose a random direction and return
	if frightCycles > 0 {

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
