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
	trappedCycles uint8
	frightCycles  uint8
	spawning      bool         // Flag set when spawning
	eaten         bool         // Flag set when eaten and returning to ghost house
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
		eaten:         false,
	}
}

/*************************** Ghost Frightened State ***************************/

// Set the fright cycles of a ghost
func (g *ghostState) setFrightCycles(cycles uint8) {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.frightCycles = cycles
	}
	g.muState.Unlock()
}

// Decrement the fright cycles of a ghost
func (g *ghostState) decFrightCycles() {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.frightCycles--
	}
	g.muState.Unlock()
}

// Get the fright cycles of a ghost
func (g *ghostState) getFrightCycles() uint8 {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return the current fright cycles
	return g.frightCycles
}

// Check if a ghost is frightened
func (g *ghostState) isFrightened() bool {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return whether there is at least one fright cycle left
	return g.frightCycles > 0
}

/****************************** Ghost Trap State ******************************/

// Set the trapped cycles of a ghost
func (g *ghostState) setTrappedCycles(cycles uint8) {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.trappedCycles = cycles
	}
	g.muState.Unlock()
}

// Decrement the trapped cycles of a ghost
func (g *ghostState) decTrappedCycles() {

	// (Write) lock the ghost state
	g.muState.Lock()
	{
		g.trappedCycles--
	}
	g.muState.Unlock()
}

// Check if a ghost is trapped
func (g *ghostState) isTrapped() bool {

	// (Read) lock the ghost state
	g.muState.RLock()
	defer g.muState.RUnlock()

	// Return whether there is at least one fright cycle left
	return g.trappedCycles > 0
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
