package game

import "sync"

// Directions:                         U   L   D   R  None
var rowDirections [5]int8 = [...]int8{-1, -0, +1, +0, +0}
var colDirections [5]int8 = [...]int8{-0, -1, +0, +1, +0}

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

// Get the row direction (-1, 0, or 1)
func (loc *locationState) dRow() int8 {
	return rowDirections[loc.dir]
}

// Get the col direction (-1, 0, or 1)
func (loc *locationState) dCol() int8 {
	return colDirections[loc.dir]
}

// Determine if another location state matches with the given location
func (loc *locationState) at(loc2 *locationState) bool {

	// (Read) lock the states (to prevent other reads or rights during the update)
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

// Copy all the attributes from another location state into the given state
func (loc *locationState) update(loc2 *locationState) {

	// (Write) lock the state (to prevent other reads or rights during the update)
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

func (loc *locationState) stepNow() {

	// (Write) lock the state (to prevent other reads or rights during the update)
	loc.Lock()
	defer loc.Unlock()

	// Add the deltas to the coordinates
	loc.row += loc.dRow()
	loc.col += loc.dCol()
}

func (loc *locationState) reverseDir() {

	// (Write) the state (to prevent other reads or rights during the update)
	loc.Lock()
	defer loc.Unlock()

	// Bitwise trick to switch between up and down, or left and right
	if loc.dir < 4 {
		loc.dir ^= 2
	}
}
