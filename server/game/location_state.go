package game

import "sync"

// Directions:                U   L   D   R  None
var dRow [5]int8 = [...]int8{-1, -0, +1, +0, +0}
var dCol [5]int8 = [...]int8{-0, -1, +0, +1, +0}

// Enum-like declaration to hold the direction indices from above
const (
	up    = 0
	left  = 1
	down  = 2
	right = 3
	none  = 4
)

/*
An object to keep track of the position and direction of an agent
*/
type locationState struct {
	row int8 // Row
	col int8 // Col
	dir int8 // Index of the direction, within the direction arrays
	sync.RWMutex
}

// Create a new location state with given position and direction values
func newLocationState(_row int8, _col int8, _dir int8) *locationState {
	return &locationState{
		row: _row,
		col: _col,
		dir: _dir,
	}
}

// Create a new location state as a copy-by-value of an existing one
func newLocationStateCopy(_loc *locationState) *locationState {

	// Just to be thread safe, (read) lock the other state
	_loc.RLock()
	defer _loc.RUnlock()

	// Copy over the variables into a new location state
	return &locationState{
		row: _loc.row,
		col: _loc.col,
		dir: _loc.dir,
	}
}

/******************************** Read Location *******************************/

// Determine if another location state matches with the given location
func (loc *locationState) collidesWith(loc2 *locationState) bool {

	// (Read) lock the states (to prevent other reads or writes)
	loc.RLock()
	loc2.RLock()
	defer func() {
		loc.RUnlock()
		loc2.RUnlock()
	}()

	// Return if both coordinates match
	return ((loc.row == loc2.row) && (loc.col == loc2.col))
}

/******************************* Update Location ******************************/

// Copy all the variables from another location state into the given location
func (loc *locationState) copyFrom(loc2 *locationState) {

	// Lock the states for thread safety
	loc.Lock()
	loc2.RLock()
	defer func() {
		loc.Unlock()
		loc2.RUnlock()
	}()

	// Update the values
	loc.row = loc2.row
	loc.col = loc2.col
	loc.dir = loc2.dir
}

/*
Set the given location to be one time step after another location,
and copy the current direction
*/
func (loc *locationState) advanceFrom(loc2 *locationState) {

	// Lock the states for thread safety
	loc.Lock()
	loc2.RLock()
	defer func() {
		loc.Unlock()
		loc2.RUnlock()
	}()

	// Update the values
	loc.row += dRow[loc2.dir]
	loc.col += dCol[loc2.dir]

	// Keep the same direction by default
	loc.dir = loc2.dir
}

// Copy all the variables from another location state into the given location
func (loc *locationState) updateDir(dir int8) {

	// Lock the state for thread safety
	loc.Lock()
	defer loc.Unlock()

	// Update the values
	loc.dir = dir
}

func (loc *locationState) reverseDir() {

	// (Write) lock the state (to prevent other reads or writes)
	loc.Lock()
	defer loc.Unlock()

	// Bitwise trick to switch between up and down, or left and right
	if loc.dir < 4 {
		loc.dir ^= 2
	}
}

func (loc *locationState) getReversedDir() int8 {

	// (Read) the state (to prevent writes)
	loc.RLock()
	defer loc.RUnlock()

	// Bitwise trick to switch between up and down, or left and right
	if loc.dir < 4 {
		return loc.dir ^ 2
	}
	return loc.dir
}

// Create a new location state as the neighbor of an existing one
func (loc *locationState) getNeighborCoords(dir int8) (int8, int8) {

	// Lock the states for thread safety
	loc.RLock()
	defer loc.RUnlock()

	// Add the deltas to the coordinates and return the pair
	return (loc.row + dRow[dir]),
		(loc.col + dCol[dir])
}
