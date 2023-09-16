# Enum class (for game mode)
from enum import Enum

class GameMode(Enum):
	''' Enum of possible game modes '''
	PAUSED = 0
	SCATTER = 1
	CHASE = 2

# TODO
class Location:
	''' Location of an entity in the game engine '''

	# Constructor
	def __init__(self) -> None:
		pass

# TODO
class Ghost:
	''' Location and auxiliary info of a ghost in the game engine '''

	# Constructor
	def __init__(self) -> None:
		pass

class GameState:
	'''
	Game state object for the Pacbot client, decoding the serialization
	from the server to make querying the game state simple.
	'''

	# Constructor
	def __init__(self) -> None:

		# A list of important game state attributes (from game engine)
		self.currTicks: int
		self.updatePeriod: int
		self.gameMode: GameMode
		self.currScore: int
		self.currLevel: int
		self.currLives: int
		self.ghosts: list[Ghost]
		self.pacmanLoc: Location
		self.fruitLoc: Location
		self.pelletArr: list[int]

	# Update, given a bytes object from the client
	def update(self, state: bytes) -> None:
		print(state)

