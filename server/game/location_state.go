package game

import "sync"

// Directions:                         U   L   D   R  None
var rowDirections [5]int8 = [...]int8{-1, -0, +1, +0, +0}
var colDirections [5]int8 = [...]int8{-0, -1, +0, +1, +0}

// Enum-like declaration to hold the directions
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
	row          int8 // Row
	col          int8 // Col
	dir          int8 // Index of the direction, within the direction arrays
	sync.RWMutex      // Mainly useful for Pacman, to make sure races don't happen
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
func newLocationStateCopy(_existingLocationState *locationState) *locationState {

	// Just to be safe, lock the existing state so its variables can be copied safely
	_existingLocationState.RLock()
	defer _existingLocationState.RUnlock()

	// Copy over the variables into a new location state
	return &locationState{
		row: _existingLocationState.row,
		col: _existingLocationState.col,
		dir: _existingLocationState.dir,
	}
}

// Get the row direction (-1, 0, or 1)
func (loc *locationState) dRow() int8 {
	return rowDirections[loc.dir]
}

// Get the col direction (-1, 0, or 1)
func (loc *locationState) dCol() int8 {
	return colDirections[loc.dir]
}

func (loc *locationState) stepNow() {

	// Lock the state (to prevent other reads or rights during the update)
	loc.Lock()
	defer loc.Unlock()

	loc.row += loc.dRow()
	loc.col += loc.dCol()
}

// TODO: Add changeDir and updatePos functions
