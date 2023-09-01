package game

import "log"

// Enum-like declaration to hold the game mode options
const (
	paused   uint8 = 0
	scatter  uint8 = 1
	chase    uint8 = 2
	numModes uint8 = 3
)

// Names of the modes (for logging)
var modeNames [numModes]string = [...]string{
	"paused",
	"scatter",
	"chase",
}

/******************************** Current Mode ********************************/

// Helper function to get the game mode
func (gs *gameState) getMode() uint8 {

	// (Read) lock the game mode
	gs.muMode.RLock()
	defer gs.muMode.RUnlock()

	// Return the current game mode
	return gs.mode
}

// Helper function to set the game mode
func (gs *gameState) setMode(mode uint8) {

	// Read the current game mode
	currMode := gs.getMode()

	// If the game is not paused and won't be paused, log the change
	if currMode != paused && mode != paused && currMode != mode {
		log.Printf("\033[36mGAME: Mode changed (%s -> %s) (t = %d)\033[0m\n",
			modeNames[currMode], modeNames[mode], gs.getCurrTicks())
	}

	// (Write) lock the game mode
	gs.muMode.Lock()
	{
		gs.mode = mode // Update the game mode
	}
	gs.muMode.Unlock()
}

/***************************** Last Unpaused Mode *****************************/

// Helper function to get the last unpaused mode
func (gs *gameState) getLastUnpausedMode() uint8 {

	// (Read) lock the game mode
	gs.muMode.RLock()
	defer gs.muMode.RUnlock()

	// If the current mode is not paused, return it
	if gs.mode != paused {
		return gs.mode
	}

	// Return the last unpaused game mode
	return gs.lastUnpausedMode
}

// Helper function to set the game mode
func (gs *gameState) setLastUnpausedMode(mode uint8) {

	// Get the last unpaused mode
	unpausedMode := gs.getLastUnpausedMode()

	// If the game is paused and the last unpaused mode changes, log the change
	if gs.getMode() == paused && unpausedMode != mode {
		log.Printf("\036[32mGAME: Mode changed while paused (%s -> %s) "+
			"(t = %d)\033[0m\n",
			modeNames[unpausedMode], modeNames[mode], gs.getCurrTicks())
	}

	// (Write) lock the game mode
	gs.muMode.Lock()
	{
		gs.lastUnpausedMode = mode // Update the game mode
	}
	gs.muMode.Unlock()
}

/******************************** Pause / Play ********************************/

// Helper function to determine if the game is paused
func (gs *gameState) isPaused() bool {
	return gs.getMode() == paused
}

// Helper function to pause the game
func (gs *gameState) pause() {

	// If the game engine is already paused, there's no more to do
	if gs.isPaused() {
		return
	}

	// Otherwise, save the current mode
	gs.setLastUnpausedMode(gs.getMode())

	// Set the mode to paused
	gs.setMode(paused)

	// Log message to alert the user
	log.Printf("\033[32m\033[2mGAME: Paused  (t = %d)\033[0m\n",
		gs.getCurrTicks())
}

// Helper function to play the game
func (gs *gameState) play() {

	// If the game engine is already playing or can't play, return
	if !gs.isPaused() || gs.getLives() == 0 || gs.getCurrTicks() == 0xffff {
		return
	}

	// Otherwise, set the current mode to the last unpaused mode
	gs.setMode(gs.getLastUnpausedMode())

	// Log message to alert the user
	log.Printf("\033[32mGAME: Resumed (t = %d)\033[0m\n",
		gs.getCurrTicks())
}

/*************************** Pausing on Next Update ***************************/

// Helper function to return whether the game should pause after next update
func (gs *gameState) getPauseOnUpdate() bool {

	// (Read) lock the game mode
	gs.muMode.RLock()
	defer gs.muMode.RUnlock()

	// Return whether the pause on update flag
	return gs.pauseOnUpdate
}

// Helper function to pause the game after the next update
func (gs *gameState) setPauseOnUpdate(flag bool) {

	// (Write) lock the game mode
	gs.muMode.Lock()
	{
		gs.pauseOnUpdate = flag // Set a flag to pause at the next update
	}
	gs.muMode.Unlock()
}

/********************************* Mode Steps *********************************/

// Helper function to get the number of steps until the mode changes
func (gs *gameState) getModeSteps() uint8 {

	// (Read) lock the mode steps
	gs.muModeSteps.RLock()
	defer gs.muModeSteps.RUnlock()

	// Return the mode steps
	return gs.modeSteps
}

// Helper function to set the number of steps until the mode changes
func (gs *gameState) setModeSteps(steps uint8) {

	// (Write) lock the mode steps
	gs.muModeSteps.Lock()
	{
		gs.modeSteps = steps // Set the mode steps
	}
	gs.muModeSteps.Unlock()
}

// Helper function to decrement the number of steps until the mode changes
func (gs *gameState) decrementModeSteps() {

	// (Write) lock the mode steps
	gs.muModeSteps.Lock()
	{
		if gs.modeSteps != 0 {
			gs.modeSteps-- // Decrease the mode steps
		}
	}
	gs.muModeSteps.Unlock()
}
