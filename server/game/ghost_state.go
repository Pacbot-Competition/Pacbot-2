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
	loc          *locationState // Current location
	nextLoc      *locationState // Planned location (for next update)
	pacmanLoc    *locationState // Position of pacman
	scatterLoc   *locationState // Position of scatter target
	gameMode     *uint8
	color        int8
	frightCycles uint8
	spawning     bool // A flag which is set when spawning
	trapped      bool // A flag which is set when trapped in the ghost house
	eaten        bool // A flag which is set when eaten and returning to house
}

// Create a new ghost state with given location and color values
func newGhostState(_color int8, _pacmanLoc *locationState, _gameMode *uint8) *ghostState {
	g := ghostState{
		pacmanLoc:    _pacmanLoc,
		gameMode:     _gameMode,
		color:        _color,
		frightCycles: 0,
		spawning:     true,
		eaten:        false,
	}
	switch _color {
	case red:
		g.loc = newLocationStateCopy(initLocRed)
		g.scatterLoc = newLocationStateCopy(initScatterTargetRed)
		g.trapped = false
	case pink:
		g.loc = newLocationStateCopy(initLocPink)
		g.scatterLoc = newLocationStateCopy(initScatterTargetPink)
		g.trapped = true
	case cyan:
		g.loc = newLocationStateCopy(initLocCyan)
		g.scatterLoc = newLocationStateCopy(initScatterTargetCyan)
		g.trapped = true
	case orange:
		g.loc = newLocationStateCopy(initLocOrange)
		g.scatterLoc = newLocationStateCopy(initScatterTargetOrange)
		g.trapped = true
	}
	return &g
}
