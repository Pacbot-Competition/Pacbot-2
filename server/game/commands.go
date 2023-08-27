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

	// If there's no message, return
	if len(msg) < 1 {
		return
	}

	// Decide the command type based on the first byte
	switch msg[0] {
	case 'p':
		gs.pause()
		fmt.Printf("\033[32m\033[2mGame  paused (t = %d)\033[0m\n", gs.getCurrTicks())
	case 'P':
		gs.play()
		fmt.Printf("\033[32mGame resumed (t = %d)\033[0m\n", gs.getCurrTicks())
	case 'w':
		gs.movePacmanDir(up)
	case 'a':
		gs.movePacmanDir(left)
	case 's':
		gs.movePacmanDir(down)
	case 'd':
		gs.movePacmanDir(right)
	}
}
