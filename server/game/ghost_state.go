package game

import "fmt"

// Enum-like declaration to hold the ghost colors
const (
	red    = 0
	pink   = 1
	cyan   = 2
	orange = 3
)

// Scatter targets for the ghosts - should remain constant
var ghostScatterLocs [4]*locationState = [...]*locationState{
	newLocationState(-3, 25, none), // red
	newLocationState(-3, 2, none),  // pink
	newLocationState(31, 0, none),  // cyan
	newLocationState(31, 27, none), // orange
}

// The number of cycles that the ghosts stay in the trapped state for
var ghostTrappedCycles [4]uint8 = [...]uint8{
	0,  // red
	1,  // pink
	8,  // cyan
	16, // orange
}

/*
An object to keep track of the location and attributes of a ghost
*/
type ghostState struct {
	loc           *locationState // Current location
	nextLoc       *locationState // Planned location (for next update)
	scatterLoc    *locationState // Position of Ifixed) scatter target
	game          *gameState     // The game state tied to the ghost
	color         int8
	trappedCycles uint8
	frightCycles  uint8
	spawning      bool // Flag set when spawning
	wasEaten      bool // Flag set when eaten and returning to ghost house
}

// Create a new ghost state with given location and color values
func newGhostState(_gameState *gameState, _color int8) *ghostState {
	return &ghostState{
		loc:           newLocationStateCopy(ghostSpawnLocs[_color]),
		nextLoc:       newLocationStateCopy(ghostSpawnLocs[_color]),
		scatterLoc:    newLocationStateCopy(ghostScatterLocs[_color]),
		game:          _gameState,
		color:         _color,
		trappedCycles: ghostTrappedCycles[_color],
		frightCycles:  10,
		spawning:      true,
		wasEaten:      false,
	}
}

// Update the ghost's position: copy the location from the next location state
func (g *ghostState) update() {
	if g.loc.collidesWith(ghostSpawnLocs[red]) {
		g.spawning = false
	}
	g.loc.copyFrom(g.nextLoc)
}

func (g *ghostState) plan() {

	// Determine the next position based on the current direction
	g.nextLoc.advanceFrom(g.loc)

	// If the ghost is trapped, reverse the current direction and return
	if g.trappedCycles > 0 {
		g.nextLoc.reverseDir()
		g.trappedCycles--
		return
	}

	/*
		Determine whether each of the four neighboring moves to the next
		location is valid, and count how many are good
	*/
	numValidMoves := 0
	var moveValid [4]bool
	//var moveDistSq [4]uint32
	for dir := int8(0); dir < 4; dir++ {

		// If this move would make the ghost reverse, skip it
		if dir == g.nextLoc.getReverseDir() {
			continue
		}

		// Get the neighboring cell in that location
		row, col := g.nextLoc.getNeighborCoords(dir)

		// Determine if that move is valid
		moveValid[dir] = !g.game.wallAt(row, col)

		/*
			Determine if the move would help the ghost escape the ghost house,
			and make it a valid one if so
		*/
		if g.spawning && row == ghostHouseExitRow && col == ghostHouseExitCol {
			moveValid[dir] = true
		}

		// Increment the valid moves counter if necessary
		if moveValid[dir] {
			numValidMoves++
		}
	}

	// If frightened, immediately choose a random direction and return
	if g.frightCycles > 0 {
		if numValidMoves == 0 {
			fmt.Printf("\033[35mWARN: Ghost #%d has nowhere to go\n\033[0m", g.color)
			return
		}

		// Generate a random index out of the valid moves
		randomNum := g.game.rng.Intn(numValidMoves)

		// Loop over all directions
		for dir, count := int8(0), 0; dir < 4; dir++ {

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
}
