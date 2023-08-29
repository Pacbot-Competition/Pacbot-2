package game

import "sync"

// Directions:                U   L   D   R  None
var dRow [5]int8 = [...]int8{-1, -0, +1, +0, +0}
var dCol [5]int8 = [...]int8{-0, -1, +0, +1, +0}

// Enum-like declaration to hold the direction indices from above
const (
	up      uint8 = 0
	left    uint8 = 1
	down    uint8 = 2
	right   uint8 = 3
	none    uint8 = 4
	numDirs uint8 = 5
)

// Names of the directions (forr debugging)
var dirNames [numDirs]string = [...]string{
	"up",
	"left",
	"down",
	"right",
	"none",
}

/*
An object to keep track of the position and direction of an agent
*/
type locationState struct {
	row int8  // Row
	col int8  // Col
	dir uint8 // Index of the direction, within the direction arrays
	sync.RWMutex
}

// Create a new location state with given position and direction values
func newLocationState(_row int8, _col int8, _dir uint8) *locationState {
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

	// If any of the rows or columns is at least 32, they don't collide
	if loc.row >= 32 || loc.col >= 32 || loc2.row >= 32 || loc2.col >= 32 {
		return false
	}

	// Return if both coordinates match
	return ((loc.row == loc2.row) && (loc.col == loc2.col))
}

// Determine if a given location state matches with the empty location
func (loc *locationState) isEmpty() bool {

	// (Read) lock the states (to prevent other reads or writes)
	loc.RLock()
	defer func() {
		loc.RUnlock()
	}()

	// Return if both coordinates match
	return ((loc.row == emptyLoc.row) && (loc.col == emptyLoc.col))
}

// Return a direction corresponding to an existing location
func (loc *locationState) getDir() uint8 {

	// Lock the states for thread safety
	loc.RLock()
	defer loc.RUnlock()

	// Return the direction
	return loc.dir
}

func (loc *locationState) getReversedDir() uint8 {

	// (Read) the state (to prevent writes)
	loc.RLock()
	defer loc.RUnlock()

	// Bitwise trick to switch between up and down, or left and right
	if loc.dir < 4 {
		return loc.dir ^ 2
	}
	return loc.dir
}

// Return a set of coordinates corresponding to an existing location
func (loc *locationState) getCoords() (int8, int8) {

	// Lock the states for thread safety
	loc.RLock()
	defer loc.RUnlock()

	// Return the pair of coordinates
	return (loc.row),
		(loc.col)
}

// Create a new set of coordinates as the neighbor of an existing location
func (loc *locationState) getNeighborCoords(dir uint8) (int8, int8) {

	// Lock the states for thread safety
	loc.RLock()
	defer loc.RUnlock()

	// Add the deltas to the coordinates and return the pair
	return (loc.row + dRow[dir]),
		(loc.col + dCol[dir])
}

/*
Return a set of coordinates a few steps ahead (in the direction it is facing)
of a given location state
*/
func (loc *locationState) getAheadCoords(spaces int8) (int8, int8) {

	// Lock the states for thread safety
	loc.RLock()
	defer loc.RUnlock()

	// Add the deltas to the coordinates and return the pair
	return (loc.row + dRow[loc.dir]*spaces),
		(loc.col + dCol[loc.dir]*spaces)
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
func (loc *locationState) updateDir(dir uint8) {

	// Lock the state for thread safety
	loc.Lock()
	defer loc.Unlock()

	// Update the values
	loc.dir = dir
}

// Reverse the direction of a given location state
func (loc *locationState) reverseDir() {

	// Lock the state for thread safety
	loc.Lock()
	defer loc.Unlock()

	// Bitwise trick to switch between up and down, or left and right
	if loc.dir < 4 {
		loc.dir ^= 2
	}
}

// Move a given location state to specified coordinates
func (loc *locationState) moveToCoords(row int8, col int8) {

	// Lock the state for thread safety
	loc.Lock()
	defer loc.Unlock()

	// Update the values
	loc.row = row
	loc.col = col
}
