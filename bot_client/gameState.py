# Enum class (for game mode)
from enum import IntEnum

# Struct class (for processing)
from struct import unpack_from, pack

# Internal representation of walls
from walls import wallArr

# Buffer to collect messages to write to the server
from collections import deque

# Terminal colors for formatting output text
from terminalColors import *

class GameModes(IntEnum):
	'''
	Enum of possible game modes
	'''

	PAUSED  = 0
	SCATTER = 1
	CHASE   = 2

# Terminal colors, based on the game mode
GameModeColors = {
	GameModes.PAUSED:  DIM,
	GameModes.CHASE:   YELLOW,
	GameModes.SCATTER: GREEN
}

class GhostColors(IntEnum):
	'''
	Enum of possible ghost names
	'''

	RED    = 0
	PINK   = 1
	CYAN   = 2
	ORANGE = 3

# Scatter targets for each of the ghosts
#               R   P   C   O
SCATTER_ROW = [-3, -3, 31, 31]
SCATTER_COL = [25,  2, 27,  0]

class Directions(IntEnum):
	'''
	Enum of possible directions for the Pacman agent
	'''

	UP    = 0
	LEFT  = 1
	DOWN  = 2
	RIGHT = 3
	NONE  = 4

# Directions:        U   L   D   R  None
D_ROW: list[int] = [-1, -0, +1, +0, +0]
D_COL: list[int] = [-0, -1, +0, +1, +0]

reversedDirections: dict[Directions, Directions] = {
	Directions.UP:    Directions.DOWN,
	Directions.LEFT:  Directions.RIGHT,
	Directions.DOWN:  Directions.UP,
	Directions.RIGHT: Directions.LEFT,
	Directions.NONE:  Directions.NONE
}

class Location:
	'''
	Location of an entity in the game engine
	'''

	def __init__(self, state) -> None: # type: ignore
		'''
		Construct a new location state object
		'''

		# Relevant game state
		self.state: GameState = state

		# Row and column information
		self.rowDir: int  = 0
		self.row: int     = 32
		self.colDir: int  = 0
		self.col: int     = 32

	def update(self, loc_uint16: int) -> None:
		'''
		Update a location, based on a 2-byte serialization
		'''

		# Get the row and column bytes
		row_uint8: int = loc_uint16 >> 8
		col_uint8: int = loc_uint16 & 0xff

		# Get the row direction (2's complement of first 2 bits)
		self.rowDir = row_uint8 >> 6
		if self.rowDir >= 2:
			self.rowDir -= 4

		# Get the row value (last 6 bits)
		self.row = row_uint8 & 0x3f

		# Get the col direction (2's complement of first 2 bits)
		self.colDir = col_uint8 >> 6
		if self.colDir >= 2:
			self.colDir -= 4

		# Get the column value (last 6 bits)
		self.col = col_uint8 & 0x3f

	def at(self, row: int, col: int) -> bool:
		'''
		Determine whether a row and column intersect with this location
		'''

		# Check the compared position is not an empty location
		if (row >= 31) or (col >= 28):
			return False

		# Return whether the rows and columns both match
		return (self.row == row) and (self.col == col)

	def serialize(self) -> int:
		'''
		Serialize this location state into a 16-bit integer (two bytes)
		'''

		# Serialize the row byte
		row_uint8: int = (((self.rowDir & 0x03) << 6) | (self.row & 0x3f))

		# Serialize the column byte
		col_uint8: int = (((self.colDir & 0x03) << 6) | (self.col & 0x3f))

		# Return the full serialization
		return (row_uint8 << 8) | (col_uint8)

	def advance(self) -> None:
		'''
		Advance this location state for simulating another step transition
		'''

		# If the current position is out of bounds, ignore it
		if (self.row > 31) or (self.col > 28):
			return

		# Calculate the next row and column
		newRow = self.row + self.rowDir
		newCol = self.col + self.colDir

		# Move to the next row and column, if applicable
		if not self.state.wallAt(newRow, newCol):
			self.row = newRow
			self.col = newCol

	def setDirection(self, direction: Directions) -> None:
		'''
		Given a direction enum object, set the direction of this location
		'''

		# Set the direction of this location
		self.rowDir = D_ROW[direction]
		self.colDir = D_COL[direction]

	def getDirection(self) -> Directions:
		'''
		Return a direction enum object corresponding to this location
		'''

		# Return the matching direction, if applicable
		for direction in Directions:
			if self.row == D_ROW[direction] and self.col == D_COL[direction]:
				return direction

		# Return none if no direction matches
		return Directions.NONE

class Ghost:
	'''
	Location and auxiliary info of a ghost in the game engine
	'''

	def __init__(self, color: GhostColors, state) -> None: # type: ignore
		'''
		Construct a new ghost state object
		'''

		# Relevant game state
		self.state: GameState = state

		# Ghost information
		self.color: GhostColors = color
		self.location: Location = Location(state) # type: ignore
		self.frightSteps: int = 0
		self.spawning: bool = bool(True)

		# (For simulation) Planned next direction the ghost will take
		self.plannedDirection: Directions = Directions.NONE

	def updateAux(self, auxInfo: int) -> None:
		'''
		Update auxiliary info (fright steps and spawning flag, 1 byte)
		'''

		self.frightSteps = auxInfo & 0x3f
		self.spawning = bool(auxInfo >> 7)

	def serializeAux(self) -> int:
		'''
		Serialize auxiliary info (fright steps and spawning flag, 1 byte)
		'''

		return (self.spawning << 7) | (self.frightSteps)

	def isFrightened(self) -> bool:
		'''
		Return whether this ghost is frightened
		'''

		return (self.frightSteps > 0)

	def move(self) -> None:
		'''
		Update the ghost's position for simulation purposes
		'''

		# As an approximation since we don't have enough info, assume the location
		# of the ghost will not change much if it is spawning (as it might be
		# trapped in the ghost house) - this holds for short-term simulations into
		# the future, but feel free to adjust it if you have a better way to
		# predict how spawning ghosts will behave
		if self.spawning:
			return

		# Advance the ghost's location
		self.location.advance()

		# Set the current direction to the guess of the planned direction
		self.location.setDirection(self.plannedDirection)

		# If the ghost is frightened, drop its steps by 1
		if self.isFrightened():
			self.frightSteps -= 1

	def guessPlan(self) -> None:
		'''
		Use incomplete knowledge of the current game state to predict where the
		ghosts might aim at the next step
		'''

		# For the same reason as in move(), ignore spawning ghosts during short-
		# term projections into the future
		if self.spawning:
			return

		# If the ghost is at an empty location, ignore it
		if self.location.row >= 32 or self.location.col >= 32:
			return

		# Row and column at the next step
		nextRow: int = self.location.row + self.location.rowDir
		nextCol: int = self.location.col + self.location.colDir

		# Pacman row and column
		pacmanRow: int = self.state.pacmanLoc.row
		pacmanCol: int = self.state.pacmanLoc.col
		pacmanRowDir: int = self.state.pacmanLoc.rowDir
		pacmanColDir: int = self.state.pacmanLoc.colDir

		# Red ghost's location
		redRow: int = self.state.ghosts[GhostColors.RED].location.row
		redCol: int = self.state.ghosts[GhostColors.RED].location.col

		# Target row and column
		targetRow: int = 0
		targetCol: int = 0

		# Choose a target for the ghost based on its color
		if self.state.gameMode == GameModes.CHASE:

			# Red targets Pacman
			if self.color == GhostColors.RED:
				targetRow = pacmanRow
				targetCol = pacmanCol

			# Pink targets the space 4 ahead of Pacman
			elif self.color == GhostColors.PINK:
				targetRow = pacmanRow + 4 * pacmanRowDir
				targetRow = pacmanCol + 4 * pacmanColDir

			# Cyan targets the position of red, reflected about the position 2 spaces
			# ahead of Pacman
			elif self.color == GhostColors.CYAN:
				targetRow = 2 * pacmanRow + 4 * pacmanRowDir - redRow
				targetCol = 2 * pacmanCol + 4 * pacmanColDir - redCol

			# Orange targets Pacman, but only if Pacman is farther than 8 spaces away
			elif self.color == GhostColors.ORANGE:
				distSqToPacman = (nextRow - pacmanRow) * (nextRow - pacmanRow) + \
													(nextCol - pacmanCol) * (nextCol - pacmanCol)
				targetRow = pacmanRow if (distSqToPacman < 64) else \
											SCATTER_ROW[GhostColors.ORANGE]
				targetCol = pacmanCol if (distSqToPacman < 64) else \
											SCATTER_COL[GhostColors.ORANGE]

		# In scatter mode, each ghost tracks a fixed target at a corner of the maze
		if self.state.gameMode == GameModes.SCATTER:
			targetRow = SCATTER_ROW[self.color]
			targetCol = SCATTER_COL[self.color]

		# Calculate the distance squared to the target, for all 4 moves
		minDist = 0xfffffff
		maxDist = -1
		minDir  = Directions.UP
		maxDir  = Directions.UP
		for direction in Directions:
			if direction != Directions.NONE:

				# Avoid reversals, as ghosts are not typically allowed to reverse
				if D_ROW[direction] + self.location.rowDir != 0 or \
					D_COL[direction] + self.location.colDir != 0:

					# Check whether this new location would be valid (not in a wall)
					newRow = nextRow + D_ROW[direction]
					newCol = nextCol + D_COL[direction]
					if not self.state.wallAt(newRow, newCol):

						# Compare the distance squared to the target to the current best;
						# if it is better, choose it to be the new ghost plan
						distSqToTarget = (newRow - targetRow) * (newRow - targetRow) + \
															(newCol - targetCol) * (newCol - targetCol)
						if distSqToTarget < minDist:
							minDir  = direction
							minDist = distSqToTarget
						elif distSqToTarget >= maxDist:
							maxDir  = direction
							maxDist = distSqToTarget

		# Update the best direction to be the plan
		self.plannedDirection = minDir if (not self.isFrightened()) else maxDir

class GameStateCompressed:
	'''
	Compressed copy of the game state, for easier storage for path planning.
	'''

	def __init__(
		self,
		serialized: bytes,
		ghostPlans: dict[GhostColors, Directions]
	) -> None:
		'''
		Construct a new compressed game state object
		'''

		# Serialization of the game state, in bytes
		self.serialized: bytes = serialized

		# Store tentative ghost plans
		self.ghostPlans: dict[GhostColors, Directions] = ghostPlans

class GameState:
	'''
	Game state object for the Pacbot client, decoding the serialization
	from the server to make querying the game state simple.
	'''

	def __init__(self) -> None:
		'''
		Construct a new game state object
		'''

		# Big endian format specifier
		self.format: str = '>'

		# Internal variable to lock the state
		self._locked: bool = False

		# Keep track of whether the client is connected
		self._connected: bool = False

		# Buffer of messages to write back to the server
		self.writeServerBuf: deque[bytes] = deque[bytes](maxlen=6)

		# Internal representation of walls:
		# 31 * 4 bytes = 31 * (32-bit integer bitset)
		self.wallArr: list[int] = wallArr

		#--- Important game state attributes (from game engine) ---#

		# 2 bytes
		self.currTicks: int = 0
		self.format += 'H'

		# 1 byte
		self.updatePeriod: int = 12
		self.format += 'B'

		# 1 byte
		self.gameMode: GameModes = GameModes.PAUSED
		self.format += 'B'

		# 2 bytes
		self.modeSteps: int = 0
		self.modeDuration: int = 255
		self.format += 'BB'

		# 2 bytes
		self.currScore: int = 0
		self.format += 'H'

		# 1 byte
		self.currLevel: int = 0
		self.format += 'B'

		# 1 byte
		self.currLives: int = 3
		self.format += 'B'

		# 4 * 3 bytes = 4 * (2 bytes location + 1 byte aux info)
		self.ghosts: list[Ghost] = [Ghost(color, self) for color in GhostColors]
		self.format += 'HBHBHBHB'

		# 2 byte location
		self.pacmanLoc: Location = Location(self)
		self.format += 'H'

		# 2 byte location
		self.fruitLoc: Location = Location(self)
		self.format += 'H'

		# 2 bytes
		self.fruitSteps: int = 0
		self.fruitDuration: int = 30
		self.format += 'BB'

		# 31 * 4 bytes = 31 * (32-bit integer bitset)
		self.pelletArr: list[int] = [0 for _ in range(31)]
		self.format += (31 * 'I')

	def lock(self) -> None:
		'''
		Lock the game state, to prevent updates
		'''

		# Lock the state by updating the internal state variable
		self._locked = True

	def unlock(self) -> None:
		'''
		Unlock the game state, to allow updates
		'''

		# Unlock the state by updating the internal state variable
		self._locked = False

	def isLocked(self) -> bool:
		'''
		Check if the game state is locked
		'''

		# Return the internal 'locked' state variable
		return self._locked

	def setConnectionStatus(self, connected: bool) -> None:
		'''
		Set the connection status of this game state's client
		'''

		# Update the internal 'connected' state variable
		self._connected = connected

	def isConnected(self) -> bool:
		'''
		Check if the client attached to the game state is connected
		'''

		# Return the internal 'connected' state variable
		return self._connected

	def serialize(self) -> bytes:
		'''
		Serialize this game state into a bytes object (for policy state storage)
		'''

		# Return a serialization with the same format as server updates
		return pack(

			# Format string
			self.format,

			# General game info
			self.currTicks,
			self.updatePeriod,
			self.gameMode,
			self.modeSteps,
			self.modeDuration,
			self.currScore,
			self.currLevel,
			self.currLives,

			# Red ghost info
			self.ghosts[GhostColors.RED].location.serialize(),
			self.ghosts[GhostColors.RED].serializeAux(),

			# Pink ghost info
			self.ghosts[GhostColors.PINK].location.serialize(),
			self.ghosts[GhostColors.PINK].serializeAux(),

			# Cyan ghost info
			self.ghosts[GhostColors.CYAN].location.serialize(),
			self.ghosts[GhostColors.CYAN].serializeAux(),

			# Orange ghost info
			self.ghosts[GhostColors.ORANGE].location.serialize(),
			self.ghosts[GhostColors.ORANGE].serializeAux(),

			# Pacman location info
			self.pacmanLoc.serialize(),

			# Fruit location info
			self.fruitLoc.serialize(),
			self.fruitSteps,
			self.fruitDuration,

			# Pellet info
			*self.pelletArr
		)

	def getGhostPlans(self) -> dict[GhostColors, Directions]:
		'''
		Return the ghosts' planned directions to compress the game state
		'''

		return {ghost.color: ghost.plannedDirection for ghost in self.ghosts}

	def update(self, serializedState: bytes, lockOverride: bool = False) -> None:
		'''
		Update this game state, given a bytes object from the client
		'''

		# If the state is locked, don't update it
		if self._locked and not lockOverride:
			return

		# Unpack the values based on the format string
		unpacked: tuple[int, ...] = unpack_from(self.format, serializedState, 0)

		# General game info
		self.currTicks    = unpacked[0]
		self.updatePeriod = unpacked[1]
		self.gameMode     = unpacked[2]
		self.modeSteps    = unpacked[3]
		self.modeDuration = unpacked[4]
		self.currScore    = unpacked[5]
		self.currLevel    = unpacked[6]
		self.currLives    = unpacked[7]

		# Red ghost info
		self.ghosts[GhostColors.RED].location.update(unpacked[8])
		self.ghosts[GhostColors.RED].updateAux(unpacked[9])

		# Pink ghost info
		self.ghosts[GhostColors.PINK].location.update(unpacked[10])
		self.ghosts[GhostColors.PINK].updateAux(unpacked[11])

		# Cyan ghost info
		self.ghosts[GhostColors.CYAN].location.update(unpacked[12])
		self.ghosts[GhostColors.CYAN].updateAux(unpacked[13])

		# Orange ghost info
		self.ghosts[GhostColors.ORANGE].location.update(unpacked[14])
		self.ghosts[GhostColors.ORANGE].updateAux(unpacked[15])

		# Pacman location info
		self.pacmanLoc.update(unpacked[16])

		# Fruit location info
		self.fruitLoc.update(unpacked[17])
		self.fruitSteps = unpacked[18]
		self.fruitDuration = unpacked[19]

		# Pellet info
		self.pelletArr = list[int](unpacked)[20:]

		# Reset our guesses of the planned ghost directions
		for ghost in self.ghosts:
			ghost.plannedDirection = Directions.NONE

	def updateGhostPlans(self, ghostPlans: dict[GhostColors, Directions]):
		'''
		Update this game state, given a list of ghost planned directions
		'''

		for ghost in self.ghosts:
			ghost.plannedDirection = ghostPlans[ghost.color]

	def pelletAt(self, row: int, col: int) -> bool:
		'''
		Helper function to check if a pellet is at a given location
		'''

		return bool((self.pelletArr[row] >> col) & 1)

	def superPelletAt(self, row: int, col: int) -> bool:
		'''
		Helper function to check if a super pellet is at a given location
		'''

		return self.pelletAt(row, col) and \
			((row == 3) or (row == 23)) and ((col == 1) or (col == 26))

	def fruitAt(self, row: int, col: int) -> bool:
		'''
		Helper function to check if a fruit is at a given location
		'''

		return (self.fruitSteps > 0) and (row == self.fruitLoc.row) and \
			(col == self.fruitLoc.col)

	def numPellets(self) -> int:
		'''
		Helper function to compute how many pellets are left in the maze
		'''

		return sum(row_arr.bit_count() for row_arr in self.pelletArr)

	def collectFruit(self, row: int, col: int) -> None:
		'''
		Helper function to collect a fruit for simulation purposes
		'''

		# Remove the fruit if we have collected it
		if self.fruitAt(row, col):
			self.currScore += 100
			self.fruitSteps = 0
			self.fruitLoc.row = 32
			self.fruitLoc.col = 32

		# Decrease the fruit steps to bring it closer to despawning
		if self.fruitSteps > 0:
			self.fruitSteps -= 1

		# If the fruit steps counter has expired, despawn it
		if self.fruitSteps == 0:
			self.fruitLoc.row = 32
			self.fruitLoc.col = 32

	def collectPellet(self, row: int, col: int) -> None:
		'''
		Helper function to collect a pellet for simulation purposes
		'''

		# Return if there are no pellets to collect
		if not self.pelletAt(row, col):
			return

		# Determine the type of pellet (super / normal)
		superPellet: bool = self.superPelletAt(row, col)

		# Remove the pellet at this location
		self.pelletArr[row] &= (~(1 << col))

		# Increase the score by this amount
		self.currScore += (50 if superPellet else 10)

		# Spawn the fruit based on the number of pellets, if applicable
		numPellets = self.numPellets()
		if numPellets == 174 or numPellets == 74:
			self.fruitSteps = 30
			self.fruitLoc.row = 17
			self.fruitLoc.col = 13

		# When <= 20 pellets are left, keep the game in chase mode
		if numPellets <= 20:
			if self.gameMode == GameModes.SCATTER:
				self.gameMode = GameModes.CHASE

		# Scare the ghosts, if applicable
		if superPellet:
			for ghost in self.ghosts:
				ghost.frightSteps = 40
				ghost.plannedDirection = reversedDirections[ghost.plannedDirection]

	def wallAt(self, row: int, col: int) -> bool:
		'''
		Helper function to check if a wall is at a given location
		'''

		return bool((self.wallArr[row] >> col) & 1)

	def display(self):
		'''
		Helper function to display the game state in the terminal
		'''

		# Begin by outputting the tick number, colored based on the mode
		out: str = f'{GameModeColors[self.gameMode]}-------'\
				f' time = {self.currTicks:5d} -------\033[0m\n'

		# Loop over all 31 rows
		for row in range(31):

			# For each cell, choose a character based on the entities in it
			for col in range(28):

				# Red ghost
				if self.ghosts[GhostColors.RED].location.at(row, col):
					scared = self.ghosts[GhostColors.RED].isFrightened()
					out += f'{RED if not scared else BLUE}@{NORMAL}'

				# Pink ghost
				elif self.ghosts[GhostColors.PINK].location.at(row, col):
					scared = self.ghosts[GhostColors.PINK].isFrightened()
					out += f'{PINK if not scared else BLUE}@{NORMAL}'

				# Cyan ghost
				elif self.ghosts[GhostColors.CYAN].location.at(row, col):
					scared = self.ghosts[GhostColors.CYAN].isFrightened()
					out += f'{CYAN if not scared else BLUE}@{NORMAL}'

				# Orange ghost
				elif self.ghosts[GhostColors.ORANGE].location.at(row, col):
					scared = self.ghosts[GhostColors.ORANGE].isFrightened()
					out += f'{ORANGE if not scared else BLUE}@{NORMAL}'

				# Pacman
				elif self.pacmanLoc.at(row, col):
					out += f'{YELLOW}P{NORMAL}'

				# Fruit
				elif self.fruitLoc.at(row, col):
					out += f'{GREEN}f{NORMAL}'

				# Wall
				elif self.wallAt(row, col):
					out += f'{DIM}#{NORMAL}'

				# Super pellet
				elif self.superPelletAt(row, col):
					out += '●'

				# Pellet
				elif self.pelletAt(row, col):
					out += '·'

				# Empty space
				else:
					out += ' '

			# New line at end of row
			out += '\n'

		# Print the output, with a new line at end of display
		print(out)

	def safetyCheck(self) -> bool:
		'''
		Helper function to check whether Pacman is safe in the current game state
		(i.e., Pacman is not directly colliding with a non-frightened ghost)
		'''

		# Retrieve Pacman's coordinates
		pacmanRow = self.pacmanLoc.row
		pacmanCol = self.pacmanLoc.col

		# Check for collisions
		for ghost in self.ghosts:
			if ghost.location.at(pacmanRow, pacmanCol):
				if not ghost.isFrightened(): # Collision; Pacman loses
					return False
				else: # 'Respawn' the ghost
					ghost.location.row = 32
					ghost.location.col = 32
					ghost.spawning = True

		# Otherwise, Pacman is safe
		return True

	def simulateAction(self, numTicks: int, pacmanDir: Directions) -> bool:
		'''
		Helper function to advance the game state (predicting the new ghost
		positions, modes, and other information) and move Pacman one space in a
		chosen direction, as a high-level path planning step

		Returns: whether this action is safe (True) or unsafe (False), in terms
		of colliding with non-frightened ghosts.
		'''

		# Try to plan the ghost directions if we expect them to be none
		for ghost in self.ghosts:
			if ghost.plannedDirection == Directions.NONE:
				ghost.guessPlan()

		# Loop over every tick
		for tick in range(1, numTicks+1):

			# Keep ticking until an update
			if (self.currTicks + tick) % self.updatePeriod != 0:
				continue

			# Update the ghost positions (and reduce frightened steps if applicable)
			for ghost in self.ghosts:
				ghost.move()

			# Return if Pacman collides with a non-frightened ghost
			if not self.safetyCheck():
				return False

			# Update the mode steps counter, and change the mode if necessary
			if self.modeSteps > 0:
				self.modeSteps -= 1

			if self.modeSteps == 0:

				# Scatter -> Chase
				if self.gameMode == GameModes.SCATTER:
					self.gameMode = GameModes.CHASE
					self.modeSteps = 180
					self.modeDuration = 180

				# Chase -> Scatter
				elif self.gameMode == GameModes.CHASE and self.numPellets() > 20:
					self.gameMode = GameModes.SCATTER
					self.modeSteps = 60
					self.modeDuration = 60

				# Reverse the planned directions of all ghosts
				for ghost in self.ghosts:
						ghost.plannedDirection = reversedDirections[ghost.plannedDirection]

			# Guess the next ghost moves (will likely be inaccurate, due to inferring
			# unknown information from other features of the game state)
			for ghost in self.ghosts:
				ghost.guessPlan()

		# If Pacman is not given a direction to move towards, return
		if pacmanDir == Directions.NONE:
			return True

		# Set the direction of Pacman, as chosen, and try to move one step
		self.pacmanLoc.setDirection(pacmanDir)
		self.pacmanLoc.advance()
		self.collectFruit(self.pacmanLoc.row, self.pacmanLoc.col)
		self.collectPellet(self.pacmanLoc.row, self.pacmanLoc.col)

		# If there are no pellets left, return
		if self.numPellets() == 0:
			return True

		# Return if Pacman collides with a non-frightened ghost
		if not self.safetyCheck():
			return False

		# Increment the number of ticks by the chosen amount
		self.currTicks += numTicks

		# Return that Pacman was safe during this transition
		return True

def compressGameState(state: GameState) -> GameStateCompressed:
	'''
	Function to compress the game state into a smaller object, for easier storage
	'''

	return GameStateCompressed(state.serialize(), state.getGhostPlans())

def decompressGameState(state: GameState, compressed: GameStateCompressed):
	'''
	Function to de-compress game state information for path planning
	'''

	# Serialization (bytes) to state
	state.update(compressed.serialized, lockOverride=True)

	# Unpack the ghost plans
	state.updateGhostPlans(compressed.ghostPlans)