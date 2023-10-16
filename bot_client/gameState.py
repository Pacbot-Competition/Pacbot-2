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

class GameMode(IntEnum):
	'''
	Enum of possible game modes
	'''

	PAUSED = 0
	SCATTER = 1
	CHASE = 2

# Terminal colors, based on the game mode
GameModeColors = {
	GameMode.PAUSED:  DIM,
	GameMode.CHASE:   YELLOW,
	GameMode.SCATTER: GREEN
}

class GhostColors(IntEnum):
	'''
	Enum of possible ghost names
	'''

	RED = 0
	PINK = 1
	CYAN = 2
	ORANGE = 3

class Direction(IntEnum):
	'''
	Enum of possible directions
	'''

	UP = 0
	LEFT = 1
	DOWN = 2
	RIGHT = 3

class Location:
	'''
	Location of an entity in the game engine
	'''

	def __init__(self) -> None:
		'''
		Construct a new location state object
		'''

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

	def at(self, row: int, col: int):
		'''
		Determine whether a row and column intersect with this location
		'''

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

class Ghost:
	'''
	Location and auxiliary info of a ghost in the game engine
	'''

	def __init__(self, color: GhostColors) -> None:
		'''
		Construct a new ghost state object
		'''

		# Set the color for this ghost
		self.color: GhostColors = color
		self.location: Location = Location()
		self.frightSteps: int = 0
		self.spawning: bool = bool(True)

	def updateAux(self, auxInfo: int) -> None:
		'''
		Update auxiliary info (fright steps and spawning flag, 1 byte)
		'''

		self.frightSteps = auxInfo & 0xff
		self.spawning = bool(auxInfo >> 7)

	def serializeAux(self) -> int:
		'''
		Serialize auxiliary info (fright steps and spawning flag, 1 byte)
		'''

		return (self.spawning << 7) | (self.frightSteps)

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
		self.gameMode: GameMode = GameMode.PAUSED
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
		self.ghosts: list[Ghost] = [Ghost(color) for color in GhostColors]
		self.format += 'HBHBHBHB'

		# 2 byte location
		self.pacmanLoc: Location = Location()
		self.format += 'H'

		# 2 byte location
		self.fruitLoc: Location = Location()
		self.format += 'H'

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

			# Pellet info
			*self.pelletArr
		)

	def update(self, serializedState: bytes, lockOverride: bool = False) -> None:
		'''
		Update this game state, given a bytes object from the client
		'''

		# If the state is locked, don't update it
		if self._locked and not lockOverride:
			return

		print(lockOverride)
		print(serializedState)

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

		# Pellet info
		self.pelletArr = list[int](unpacked)[18:]

		# Display the game state (i.e., terminal printer)
		# self.display()

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

	def wallAt(self, row: int, col: int) -> bool:
		'''
		Helper function to check if a wall is at a given location
		'''

		return bool((self.wallArr[row] >> col) & 1)

	def display(self):
		'''
		Helper function to display the game state in the terminal
		'''

		# Print the tick number, colored based on the mode
		print(f'{GameModeColors[self.gameMode]}-------'\
				f' time = {self.currTicks:5d} -------\033[0m')

		# Loop over all 31 rows
		for row in range(31):

			# For each cell, choose a character based on the entities in it
			for col in range(28):

				# Red ghost
				if self.ghosts[GhostColors.RED].location.at(row, col):
					scared = self.ghosts[GhostColors.RED].frightSteps > 0
					print(f'{RED if not scared else BLUE}@{NORMAL}', end='')

				# Pink ghost
				elif self.ghosts[GhostColors.PINK].location.at(row, col):
					scared = self.ghosts[GhostColors.PINK].frightSteps > 0
					print(f'{PINK if not scared else BLUE}@{NORMAL}', end='')

				# Cyan ghost
				elif self.ghosts[GhostColors.CYAN].location.at(row, col):
					scared = self.ghosts[GhostColors.CYAN].frightSteps > 0
					print(f'{CYAN if not scared else BLUE}@{NORMAL}', end='')

				# Orange ghost
				elif self.ghosts[GhostColors.ORANGE].location.at(row, col):
					scared = self.ghosts[GhostColors.ORANGE].frightSteps > 0
					print(f'{ORANGE if not scared else BLUE}@{NORMAL}', end='')

				# Pacman
				elif self.pacmanLoc.at(row, col):
					print(f'{YELLOW}P{NORMAL}', end='')

			  # Fruit
				elif self.fruitLoc.at(row, col):
					print(f'{GREEN}f{NORMAL}', end='')

				# Wall
				elif self.wallAt(row, col):
					print(f'{DIM}#{NORMAL}', end='')

				# Super pellet
				elif self.superPelletAt(row, col):
					print('●', end='')

				# Pellet
				elif self.pelletAt(row, col):
					print('·', end='')

				# Empty space
				else:
					print(' ', end='')

			# New line at end of row
			print()

		# New line at end of display
		print()