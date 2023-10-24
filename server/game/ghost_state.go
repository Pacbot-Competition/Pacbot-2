package game

import (
	"sync"
)

// Enum-like declaration to hold the ghost colors
const (
	red       uint8 = 0
	pink      uint8 = 1
	cyan      uint8 = 2
	orange    uint8 = 3
	numColors uint8 = 4
)

/*
The number of "active" ghosts (the others are invisible and don't affect
the progression of the game)
*/
var numActiveGhosts uint8 = 4

// Configure the number of active ghosts
func ConfigNumActiveGhosts(_numActiveGhosts uint8) {
	numActiveGhosts = _numActiveGhosts
}

// Names of the ghosts (not the nicknames, just the colors, for debugging)
var ghostNames [numColors]string = [...]string{
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
	scatterTarget *locationState // Position of (fixed) scatter target
	game          *gameState     // The game state tied to the ghost
	color         uint8
	trappedSteps  uint8
	frightSteps   uint8
	spawning      bool         // Flag set when spawning
	eaten         bool         // Flag set when eaten and returning to ghost house
	muState       sync.RWMutex // Mutex to lock general state parameters
}

// Create a new ghost state with given location and color values
func newGhostState(_gameState *gameState, _color uint8) *ghostState {

	// Ghost state object
	g := ghostState{
		loc:           newLocationStateCopy(emptyLoc),
		nextLoc:       newLocationStateCopy(ghostSpawnLocs[_color]),
		scatterTarget: newLocationStateCopy(ghostScatterTargets[_color]),
		game:          _gameState,
		color:         _color,
		trappedSteps:  ghostTrappedSteps[_color],
		frightSteps:   0,
		spawning:      true,
		eaten:         false,
	}

	// If the color is greater than the number of active ghosts, hide this ghost
	if _color >= numActiveGhosts {
		g.nextLoc = newLocationStateCopy(emptyLoc)
	}

	// Return the ghost state
	return &g
}

/*************************** Ghost Frightened State ***************************/

// Set the fright steps of a ghost
func (g *ghostState) setFrightSteps(steps uint8) {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.frightSteps = steps
	}
	g.muState.Unlock()
}

// Decrement the fright steps of a ghost
func (g *ghostState) decFrightSteps() {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.frightSteps--
	}
	g.muState.Unlock()
}

// Get the fright steps of a ghost
func (g *ghostState) getFrightSteps() uint8 {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return the current fright steps
	return g.frightSteps
}

// Check if a ghost is frightened
func (g *ghostState) isFrightened() bool {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return whether there is at least one fright step left
	return g.frightSteps > 0
}

/****************************** Ghost Trap State ******************************/

// Set the trapped steps of a ghost
func (g *ghostState) setTrappedSteps(steps uint8) {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.trappedSteps = steps
	}
	g.muState.Unlock()
}

// Decrement the trapped steps of a ghost
func (g *ghostState) decTrappedSteps() {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.trappedSteps--
	}
	g.muState.Unlock()
}

// Check if a ghost is trapped
func (g *ghostState) isTrapped() bool {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return whether there is at least one fright step left
	return g.trappedSteps > 0
}

/**************************** Ghost Spawning State ****************************/

// Set the ghost spawning flag
func (g *ghostState) setSpawning(spawning bool) {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.spawning = spawning
	}
	g.muState.Unlock()
}

// Check if a ghost is spawning
func (g *ghostState) isSpawning() bool {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return the current ghost spawning flag
	return g.spawning
}

/****************************** Ghost Eaten Flag ******************************/

// Set the ghost eaten flag
func (g *ghostState) setEaten(eaten bool) {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.eaten = eaten
	}
	g.muState.Unlock()
}

// Check if a ghost is eaten
func (g *ghostState) isEaten() bool {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return the current ghost eaten flag
	return g.eaten
}
