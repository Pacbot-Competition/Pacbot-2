package game

// Serialize all the information of the game state
func (gs *gameState) serFull(outputBuf []byte, startIdx int) int {

	// Packet header - contains the necessary information to render the ticker
	startIdx = gs.serCurrTicks(outputBuf, startIdx)
	startIdx = gs.serUpdatePeriod(outputBuf, startIdx)
	startIdx = gs.serGameMode(outputBuf, startIdx)

	// Ghosts - serializes the ghost states to the buffer
	startIdx = gs.serGhosts(outputBuf, startIdx)

	// Pacman - serializes the pacman location to the buffer
	startIdx = gs.serPacman(outputBuf, startIdx)

	// Fruit - serializes the fruit location (null if fruit doesn't exist)
	startIdx = gs.serFruit(outputBuf, startIdx)

	// Pellets - serializes the pellets to the buffer
	startIdx = gs.serPellets(outputBuf, startIdx)

	// Return the starting index of the next field
	return startIdx
}
