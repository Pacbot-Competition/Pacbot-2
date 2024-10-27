package game

import (
	"log"
)

/***************************** Interpret Commands *****************************/

// Convert byte messages from clients into commands to the game state
func (ge *GameEngine) interpretCommand(msg []byte) {

	// Log the command if necessary
	if getCommandLogEnable() {
		if len(msg) > 1 {
			log.Printf("\033[2m\033[35mCOMM: %c %v\033[0m", msg[0], msg[1:])
		} else {
			log.Printf("\033[2m\033[35mCOMM: %c\033[0m", msg[0])
		}
	}

	// Decide the command type based on the first byte
	switch msg[0] {

	// Pause command
	case 'p':
		ge.state.pause()
	// Play command
	case 'P':
		ge.state.play()
	case 'R': 
		ge.reset()
	// Move up (decrease row index)
	case 'w':
		ge.state.movePacmanDir(up)

	// Move left (decrease column index)
	case 'a':
		ge.state.movePacmanDir(left)

	// Move down (increase row index)
	case 's':
		ge.state.movePacmanDir(down)

	// Move right (increase column index)
	case 'd':
		ge.state.movePacmanDir(right)
	
	// Absolute position (from tracking)
	case 'x':
		ge.state.movePacmanAbsolute(int8(msg[1]), int8(msg[2]))
	}
}
