package game

import (
	"fmt"
)

/***************************** Interpret Commands *****************************/

// Convert byte messages from clients into commands to the game state
func (gs *gameState) interpretCommand(msg []byte) {

	// Log the command if necessary
	if getCommandLogEnable() {
		fmt.Printf("\033[2m\033[36m| Response: %s`\033[0m\n", string(msg))
	}

	// Decide the command type based on the first byte
	switch msg[0] {

	// Pause command
	case 'p':
		gs.pause()
		fmt.Printf("\033[32m\033[2mGame paused  (t = %d)\033[0m\n",
			gs.getCurrTicks())

	// Play command
	case 'P':
		gs.play()
		fmt.Printf("\033[32mGame resumed (t = %d)\033[0m\n",
			gs.getCurrTicks())

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
	}
}
