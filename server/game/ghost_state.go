package game

// Enum-like declaration to hold the ghost colors
const (
	red    = 0
	pink   = 1
	cyan   = 2
	orange = 3
)

/*
An object to keep track of the location and attributes of a ghost
*/
type ghostState struct {
	loc          *locationState
	pacmanLoc    *locationState // TODO: Implement targeting
	scatterLoc   *locationState // TODO: Implement scatter locations
	color        int8
	frightCycles uint8
	spawning     bool
}

// Create a new ghost state with given location and color values
func newGhostState(_color int8, _pacmanLoc *locationState) *ghostState {
	g := ghostState{
		pacmanLoc:    _pacmanLoc,
		color:        _color,
		frightCycles: 0,
		spawning:     true,
	}
	switch _color {
	case red:
		g.loc = newLocationStateCopy(initLocRed)
	case pink:
		g.loc = newLocationStateCopy(initLocPink)
	case cyan:
		g.loc = newLocationStateCopy(initLocCyan)
	case orange:
		g.loc = newLocationStateCopy(initLocOrange)
	}
	return &g
}
