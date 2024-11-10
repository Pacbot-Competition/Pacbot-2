package game

import (
	"log"
)

/***************************** Interpret Commands *****************************/

// Convert byte messages from clients into commands to the game state
func (gs *gameState) interpretCommand(msg []byte) {

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
		gs.pause()

	// Play command
	case 'P':
		gs.play()

	// Move up (decrease row index)
	case 'w':
		gs.movePacmanDir(up)

	// Move left (decrease column index)
	case 'a':
		gs.movePacmanDir(left)

	// Move down (increase row index)
	case 's':
		gs.movePacmanDir(down)

	// Move right (increase column index)
	case 'd':
		gs.movePacmanDir(right)
	
	// Absolute position (from tracking)
	case 'x':
		gs.movePacmanAbsolute(int8(msg[1]), int8(msg[2]))
	}
}
