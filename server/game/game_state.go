package game

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

/*
	NOTE: at 24 ticks/sec, currTicks will experience an integer overflow after
	about 45 minutes, so don't run it continuously for too long (indefinite
	pausing is fine though, as it doesn't increment the current tick amount)
*/

/*
A game state object, to hold the internal game state and provide
helper methods that can be accessed by the game engine
*/
type gameState struct {

	/* Message header - 4 bytes */

	currTicks uint16       // Current ticks (see note above)
	muTicks   sync.RWMutex // Associated mutex

	updatePeriod uint8        // Ticks / update
	muPeriod     sync.RWMutex // Associated mutex

	lastUnpausedMode uint8        // Last unpaused mode (for pausing purposes)
	mode             uint8        // Game mode
	muMode           sync.RWMutex // Associated mutex
	pauseOnUpdate    bool         // Should pause when an update is ready

	// The number of steps (update periods) before the mode changes
	modeSteps   uint8
	muModeSteps sync.RWMutex // Associated mutex

	// The number of steps (update periods) before a speedup penalty starts
	levelSteps   uint16
	muLevelSteps sync.RWMutex // Associated mutex

	/* Game information - 4 bytes */

	currScore uint16       // Current score
	muScore   sync.RWMutex // Associated mutex

	currLevel uint8        // Current level (by default, starts at 1)
	muLevel   sync.RWMutex // Associated mutex

	currLives uint8        // Current lives (by default, starts at 3)
	muLives   sync.RWMutex // Associated mutex

	/* Pacman location - 2 bytes */

	pacmanLoc *locationState

	// A mutex for synchronizing updates to Pacman
	muPacman sync.Mutex

	/* Fruit location - 2 bytes */

	fruitLoc *locationState

	// The number of steps (update periods) before fruit disappears
	fruitSteps uint8
	muFruit    sync.RWMutex // Associated mutex

	/* Ghosts - 4 * 3 = 12 bytes */

	ghosts []*ghostState

	// A mutex for synchronizing simultaneous updates of all ghosts
	muGhosts sync.Mutex

	// A wait group for synchronizing updates of multiple ghosts
	wgGhosts *sync.WaitGroup

	// A variable to keep track of the current ghost combo
	ghostCombo uint8

	/* Pellet State - 31 * 4 = 124 bytes */

	// Pellets encoded within an array, with each uint32 acting as a bit array
	pellets    [mazeRows]uint32
	numPellets uint16       // Number of pellets
	muPellets  sync.RWMutex // Associated mutex

	/* Auxiliary (non-serialized) state information */

	// Wall state
	walls [mazeRows]uint32

	// A random number generator for making frightened ghost decisions
	rng *rand.Rand
}

// Create a new game state with default values
func newGameState() *gameState {

	// New game state object
	gs := gameState{

		// Message header
		currTicks:    0,
		updatePeriod: initUpdatePeriod,
		mode:         paused,

		// Additional header-related info
		lastUnpausedMode: initMode,
		pauseOnUpdate:    false,
		modeSteps:        modeDurations[initMode],
		levelSteps:       levelDuration,

		// Game info
		currScore: 0,
		currLevel: initLevel,
		currLives: initLives,

		// Fruit
		fruitSteps: 0,

		// Ghosts
		ghosts:     make([]*ghostState, numColors),
		wgGhosts:   &sync.WaitGroup{},
		ghostCombo: 0,

		// RNG (random number generation) source
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),

		// Pellet count at the start
		numPellets: initPelletCount,
	}

	// Declare the initial locations of Pacman and the fruit
	gs.pacmanLoc = newLocationStateCopy(pacmanSpawnLoc)
	gs.fruitLoc = newLocationStateCopy(fruitSpawnLoc)

	// Initialize the ghosts
	for color := uint8(0); color < numColors; color++ {
		gs.ghosts[color] = newGhostState(&gs, color)
	}

	// Copy over maze bit arrays
	copy(gs.pellets[:], initPellets[:])
	copy(gs.walls[:], initWalls[:])

	// Return the new game state
	return &gs
}

/**************************** Curr Ticks Functions ****************************/

// Helper function to get the current ticks
func (gs *gameState) getCurrTicks() uint16 {

	// (Read) lock the current ticks
	gs.muTicks.RLock()
	defer gs.muTicks.RUnlock()

	// Return the current ticks
	return gs.currTicks
}

// Helper function to increment the current ticks
func (gs *gameState) nextTick() {

	// Get the current number of ticks
	currTicks := gs.getCurrTicks()

	// If the current ticks are at the maximum, return
	if currTicks == 0xffff {
		return
	} else if currTicks == 0xfffe {
		gs.pause()
		log.Println("\033[31mGAME: Max tick limit reached\033[0m")
	}

	// (Write) lock the current ticks
	gs.muTicks.Lock()
	{
		gs.currTicks++ // Update the current ticks
	}
	gs.muTicks.Unlock()
}

/**************************** Upd Period Functions ****************************/

// Helper function to get the update period
func (gs *gameState) getUpdatePeriod() uint8 {

	// (Read) lock the update period
	gs.muPeriod.RLock()
	defer gs.muPeriod.RUnlock()

	// Return the update period
	return gs.updatePeriod
}

// Helper function to set the update period
func (gs *gameState) setUpdatePeriod(period uint8) {

	// Send a message to the terminal
	log.Printf("\033[36mGAME: Update period changed (%d -> %d) (t = %d)\033[0m\n",
		gs.getUpdatePeriod(), period, gs.getCurrTicks())

	// (Write) lock the update period
	gs.muPeriod.Lock()
	{
		gs.updatePeriod = period // Update the update period
	}
	gs.muPeriod.Unlock()
}

/******************************* Mode Functions *******************************/

// See game_modes.go, there were a lot of mode functions so I moved them there

/**************************** Game Score Functions ****************************/

// Helper function to get the current score of the game
func (gs *gameState) getScore() uint16 {

	// (Read) lock the current score
	gs.muScore.RLock()
	defer gs.muScore.RUnlock()

	// Return the current score
	return gs.currScore
}

// (For performance) helper function to increment the current score of the game
func (gs *gameState) incrementScore(change uint16) {

	// Calculate the next score, capping at the maximum 16-bit unsigned int
	score := uint32(gs.currScore)
	score = min(score+uint32(change), 65535)

	// (Write) lock the current score
	gs.muScore.Lock()
	{
		gs.currScore = uint16(score) // Update the current score
	}
	gs.muScore.Unlock()
}

/**************************** Game Level Functions ****************************/

// Helper function to get the current level of the game
func (gs *gameState) getLevel() uint8 {

	// (Read) lock the current level
	gs.muLevel.RLock()
	defer gs.muLevel.RUnlock()

	// Return the current level
	return gs.currLevel
}

// Helper function to set the current level of the game
func (gs *gameState) setLevel(level uint8) {

	// (Write) lock the current level
	gs.muLevel.Lock()
	{
		gs.currLevel = level // Update the level

		// Adjust the initial update period accordingly
		suggestedPeriod := int(initUpdatePeriod) - 2*(int(level)-1)
		gs.setUpdatePeriod(uint8(max(1, suggestedPeriod)))
	}
	gs.muLevel.Unlock()
}

// Helper function to increment the game level
func (gs *gameState) incrementLevel() {

	// Keep track of the current level
	level := gs.getLevel()

	// If we are at the last level, don't increment it anymore
	if level == 255 {
		return
	}

	// Send a message to the terminal
	log.Printf("\033[32mGAME: Next level (%d -> %d) (t = %d)\033[0m\n",
		level, level+1, gs.getCurrTicks())

	// (Write) lock the current lives
	gs.muLevel.Lock()
	{
		gs.currLevel++ // Update the level

		// Adjust the initial update period accordingly
		suggestedPeriod := int(initUpdatePeriod) - 2*int(level)
		gs.setUpdatePeriod(uint8(max(1, suggestedPeriod)))
	}
	gs.muLevel.Unlock()
}

/**************************** Game Level Functions ****************************/

// Helper function to get the lives left
func (gs *gameState) getLives() uint8 {

	// (Read) lock the current lives
	gs.muLives.RLock()
	defer gs.muLives.RUnlock()

	// Return the current lives
	return gs.currLives
}

// Helper function to set the lives left
func (gs *gameState) setLives(lives uint8) {

	// Send a message to the terminal
	log.Printf("\033[36mGAME: Lives changed (%d -> %d)\033[0m\n",
		gs.getLives(), lives)

	// (Write) lock the current lives
	gs.muLives.Lock()
	{
		gs.currLives = lives // Update the lives
	}
	gs.muLives.Unlock()
}

// Helper function to decrement the lives left
func (gs *gameState) decrementLives() {

	// Keep track of how many lives Pacman has left
	lives := gs.getLives()

	// If there were no lives, no need to decrement any more
	if lives == 0 {
		return
	}

	// Send a message to the terminal
	log.Printf("\033[31mGAME: Pacman lost a life (%d -> %d) (t = %d)\033[0m\n",
		lives, lives-1, gs.getCurrTicks())

	// (Write) lock the current lives
	gs.muLives.Lock()
	{
		gs.currLives-- // Update the lives
	}
	gs.muLives.Unlock()
}

/****************************** Pellet Functions ******************************/

// Helper function to get the number of pellets
func (gs *gameState) getNumPellets() uint16 {

	// (Read) lock the number of pellets
	gs.muPellets.RLock()
	defer gs.muPellets.RUnlock()

	// Return the number of pellets
	return gs.numPellets
}

// Helper function to decrement the number of pellets
func (gs *gameState) decrementNumPellets() {

	// (Write) lock the number of pellets
	gs.muPellets.Lock()
	{
		if gs.numPellets != 0 {
			gs.numPellets--
		}
	}
	gs.muPellets.Unlock()
}

// Reset all the pellets on the board
func (gs *gameState) resetPellets() {

	// (Write) lock the pellets array and number of pellets
	gs.muPellets.Lock()
	{
		// Copy over pellet bit array
		copy(gs.pellets[:], initPellets[:])

		// Set the number of pellets to be the default
		gs.numPellets = initPelletCount
	}
	gs.muPellets.Unlock()
}

/************************** Fruit Spawning Functions **************************/

// Helper function to get the number of steps until the fruit disappears
func (gs *gameState) getFruitSteps() uint8 {

	// (Read) lock the fruit steps
	gs.muFruit.RLock()
	defer gs.muFruit.RUnlock()

	// Return the fruit steps
	return gs.fruitSteps
}

// Helper function to determine whether fruit exists
func (gs *gameState) fruitExists() bool {

	// Return whether fruit exists
	return gs.getFruitSteps() > 0
}

// Helper function to set the number of steps until the fruit disappears
func (gs *gameState) setFruitSteps(steps uint8) {

	// (Write) lock the fruit steps
	gs.muFruit.Lock()
	{
		gs.fruitSteps = steps // Set the fruit steps
	}
	gs.muFruit.Unlock()
}

// Helper function to decrement the number of fruit steps
func (gs *gameState) decrementFruitSteps() {

	// (Write) lock the fruit steps
	gs.muFruit.Lock()
	{
		if gs.fruitSteps != 0 {
			gs.fruitSteps-- // Decrease the fruit steps
		}
	}
	gs.muFruit.Unlock()
}

/***************************** Level Steps Passed *****************************/

// Helper function to get the number of steps until the level speeds up
func (gs *gameState) getLevelSteps() uint16 {

	// (Read) lock the level steps
	gs.muLevelSteps.RLock()
	defer gs.muLevelSteps.RUnlock()

	// Return the level steps
	return gs.levelSteps
}

// Helper function to set the number of steps until the level speeds up
func (gs *gameState) setLevelSteps(steps uint16) {

	// (Write) lock the level steps
	gs.muLevelSteps.Lock()
	{
		gs.levelSteps = steps // Set the level steps
	}
	gs.muLevelSteps.Unlock()
}

// Helper function to decrement the number of steps until the mode changes
func (gs *gameState) decrementLevelSteps() {

	// (Write) lock the level steps
	gs.muLevelSteps.Lock()
	{
		if gs.levelSteps != 0 {
			gs.levelSteps-- // Decrease the level steps
		}
	}
	gs.muLevelSteps.Unlock()
}

/***************************** Step-Related Events ****************************/

// Helper function to handle step-related events, if the mode steps hit 0
func (gs *gameState) handleStepEvents() {

	// Get the current mode steps
	modeSteps := gs.getModeSteps()

	// Get the current level steps
	levelSteps := gs.getLevelSteps()

	// If the mode steps are 0, change the mode
	if modeSteps == 0 {
		switch gs.getMode() {
		// chase -> scatter
		case chase:
			gs.setMode(scatter)
			gs.setModeSteps(modeDurations[scatter])
		// scatter -> chase
		case scatter:
			gs.setMode(chase)
			gs.setModeSteps(modeDurations[chase])
		case paused:
			switch gs.getLastUnpausedMode() {
			// chase -> scatter
			case chase:
				gs.setLastUnpausedMode(scatter)
				gs.setModeSteps(modeDurations[scatter])
			// scatter -> chase
			case scatter:
				gs.setLastUnpausedMode(chase)
				gs.setModeSteps(modeDurations[chase])
			}
		}

		// Reverse the directions of all ghosts to indicate a mode switch
		gs.reverseAllGhosts()
	}

	// If the level steps are 0, add a penalty by speeding up the game
	if levelSteps == 0 {

		// Log the change to the terminal
		log.Println("\033[31mGAME: Long-game penalty applied\033[0m")

		// Drop the update period by 2
		gs.setUpdatePeriod(uint8(max(1, int(gs.getUpdatePeriod())-2)))

		// Reset the level steps to the level penalty duration
		gs.setLevelSteps(levelPenaltyDuration)
	}

	// Decrement the mode steps
	gs.decrementModeSteps()

	// Decrement the level steps
	gs.decrementLevelSteps()

	// Decrement the fruit steps
	gs.decrementFruitSteps()
}
