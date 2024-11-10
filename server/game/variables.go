package game

// The number of rows in the pellets and walls states
const mazeRows int8 = 31

// The number of columns in the pellets and walls states
const mazeCols int8 = 28

// The update period that the game starts with by default
const initUpdatePeriod uint8 = 12

// The number of steps (update periods) that pass before the level speeds up
const levelDuration uint16 = 960 // 8 minutes at 24 fps, update period = 12

// The number of steps (update periods) before a level speeds up further
const levelPenaltyDuration uint16 = 240 // 2 min (24fps, update period = 12)

// The mode that the game starts on by default
const initMode uint8 = scatter

// The lengths of the game modes, in units of steps (update periods)
var modeDurations [numModes]uint8 = [...]uint8{
	255, // paused
	60,  // scatter - 30 seconds at 24 fps
	180, // chase   - 90 seconds at 24 fps
}

// The level that Pacman starts on by default
const initLevel uint8 = 1

// The number of lives that Pacman starts with
const initLives uint8 = 3

// The coordinates where the ghost house exit is located
const ghostHouseExitRow int8 = 12
const ghostHouseExitCol int8 = 13

// Spawn position for Pacman
var pacmanSpawnLoc = newLocationState(23, 13, right)

// Spawn position for the fruit
var fruitSpawnLoc = newLocationState(17, 13, none)

// The number of steps that the fruit stays on the maze for
const fruitDuration uint8 = 30

// The points earned upon collecting a fruit
const fruitPoints uint16 = 100

// "Invalid" location - serializes to 0x00100000 0x00100000
var emptyLoc = newLocationState(32, 32, none)

// Spawn positions for the ghosts
var ghostSpawnLocs [numColors]*locationState = [...]*locationState{
	newLocationState(11, 13, left), // red
	newLocationState(13, 13, down), // pink
	newLocationState(14, 11, up),   // cyan
	newLocationState(14, 15, up),   // orange
}

// Scatter targets for the ghosts - should remain constant
var ghostScatterTargets [numColors]*locationState = [...]*locationState{
	newLocationState(-3, 25, none), // red
	newLocationState(-3, 2, none),  // pink
	newLocationState(31, 27, none), // cyan
	newLocationState(31, 0, none),  // orange
}

// The number of steps that the ghosts stay in the trapped state for
var ghostTrappedSteps [numColors]uint8 = [...]uint8{
	0,  // red
	5,  // pink
	16, // cyan
	32, // orange
}

// The number of steps that the ghosts stay in the frightened state for
const ghostFrightSteps uint8 = 40

// The number of pellets in a typical game of Pacman
const initPelletCount uint16 = 244

// The number of pellets at which to spawn the first fruit
const fruitThreshold1 uint16 = 174

// The number of pellets at which to spawn the second fruit
const fruitThreshold2 uint16 = 74

// The number of pellets at which to make the ghosts angry
const angerThreshold1 uint16 = 20

// The number of pellets at which to make the ghosts angrier
const angerThreshold2 uint16 = 10

// The points earned when collecting a pellet
const pelletPoints uint16 = 10

// The points earned when collecting a pellet
const superPelletPoints uint16 = 50

// The multiplier for the combo from catching successive frightened ghosts
const comboMultiplier uint16 = 200

// Column-wise, this may look backwards; column 0 is at bit 0 on the right
// (Tip: Ctrl+F '1' to see the initial pellet locations)
var initPellets [mazeRows]uint32 = [...]uint32{
	//                middle
	// col:             vv    8 6 4 2 0
	0b0000_0000000000000000000000000000, // row 0
	0b0000_0111111111111001111111111110, // row 1
	0b0000_0100001000001001000001000010, // row 2
	0b0000_0100001000001001000001000010, // row 3
	0b0000_0100001000001001000001000010, // row 4
	0b0000_0111111111111111111111111110, // row 5
	0b0000_0100001001000000001001000010, // row 6
	0b0000_0100001001000000001001000010, // row 7
	0b0000_0111111001111001111001111110, // row 8
	0b0000_0000001000000000000001000000, // row 9
	0b0000_0000001000000000000001000000, // row 10
	0b0000_0000001000000000000001000000, // row 11
	0b0000_0000001000000000000001000000, // row 12
	0b0000_0000001000000000000001000000, // row 13
	0b0000_0000001000000000000001000000, // row 14
	0b0000_0000001000000000000001000000, // row 15
	0b0000_0000001000000000000001000000, // row 16
	0b0000_0000001000000000000001000000, // row 17
	0b0000_0000001000000000000001000000, // row 18
	0b0000_0000001000000000000001000000, // row 19
	0b0000_0111111111111001111111111110, // row 20
	0b0000_0100001000001001000001000010, // row 21
	0b0000_0100001000001001000001000010, // row 22
	0b0000_0111001111111001111111001110, // row 23
	0b0000_0001001001000000001001001000, // row 24
	0b0000_0001001001000000001001001000, // row 25
	0b0000_0111111001111001111001111110, // row 26
	0b0000_0100000000001001000000000010, // row 27
	0b0000_0100000000001001000000000010, // row 28
	0b0000_0111111111111111111111111110, // row 29
	0b0000_0000000000000000000000000000, // row 30
}

// Column-wise, this may look backwards; column 0 is at bit 0 on the right
// (Tip: Ctrl+F '0' to see the valid Pacman locations)
var initWalls [mazeRows]uint32 = [...]uint32{
	//                middle
	// col:             vv    8 6 4 2 0
	0b0000_1111111111111111111111111111, // row 0
	0b0000_1000000000000110000000000001, // row 1
	0b0000_1011110111110110111110111101, // row 2
	0b0000_1011110111110110111110111101, // row 3
	0b0000_1011110111110110111110111101, // row 4
	0b0000_1000000000000000000000000001, // row 5
	0b0000_1011110110111111110110111101, // row 6
	0b0000_1011110110111111110110111101, // row 7
	0b0000_1000000110000110000110000001, // row 8
	0b0000_1111110111110110111110111111, // row 9
	0b0000_1111110111110110111110111111, // row 10
	0b0000_1111110110000000000110111111, // row 11
	0b0000_1111110110111111110110111111, // row 12
	0b0000_1111110110111111110110111111, // row 13
	0b0000_1111110000111111110000111111, // row 14
	0b0000_1111110110111111110110111111, // row 15
	0b0000_1111110110111111110110111111, // row 16
	0b0000_1111110110000000000110111111, // row 17
	0b0000_1111110110111111110110111111, // row 18
	0b0000_1111110110111111110110111111, // row 19
	0b0000_1000000000000110000000000001, // row 20
	0b0000_1011110111110110111110111101, // row 21
	0b0000_1011110111110110111110111101, // row 22
	0b0000_1000110000000000000000110001, // row 23
	0b0000_1110110110111111110110110111, // row 24
	0b0000_1110110110111111110110110111, // row 25
	0b0000_1000000110000110000110000001, // row 26
	0b0000_1011111111110110111111111101, // row 27
	0b0000_1011111111110110111111111101, // row 28
	0b0000_1000000000000000000000000001, // row 29
	0b0000_1111111111111111111111111111, // row 30
}
