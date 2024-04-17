# Asyncio (for concurrency)
import asyncio

# Import connection state object
from connectionState import ConnectionState

# Import the wall array
from walls import wallArr

# Random number generator
from random import randint

class CameraModule:
	'''
	Sample implementation of a decision module for computer vision
	for Pacbot, using asyncio.
	'''

	def __init__(self, state: ConnectionState) -> None:
		'''
		Construct a new decision module object
		'''

		# Game state object to store the game information
		self.state = state

	async def decisionLoop(self) -> None:
		'''
		Decision loop for CV
		'''

		# Receive values as long as we have access
		while self.state.isConnected():

			# Write back to the server, as a test (move right)
			self.state.send(randint(0, 30), randint(0, 27))

			# Free up the event loop
			await asyncio.sleep(0)

	@classmethod
	def getClosestCoords(cls, floatRow: float, floatCol: float) -> tuple[int, int]:

		# Translation; account for the fact that cell coordinates are centered
		floatRow -= 0.5
		floatCol -= 0.5

		# Round to the nearest coordinates
		intRow, intCol = round(floatRow), round(floatCol)

		# If the coordinates are not in a wall, return them
		if not CameraModule.wallAt(intRow, intCol):
			return (intRow, intCol)

		# Otherwise, check the neighbors in a 3x3 square
		neighbors: list[tuple[float, tuple[int, int]]] = []
		for row in [intRow - 1, intRow, intRow + 1]:
			for col in [intCol - 1, intCol, intCol + 1]:
				if not CameraModule.wallAt(row, col):
					distSq = (row - floatRow) * (row - floatRow) + (col - floatCol) + (col - floatCol)
					neighbors.append((distSq, (row, col)))

		# Sort the neighbors by distance
		neighbors = sorted(neighbors)

		# Return the closest neighbor not in a wall
		if len(neighbors):
			return neighbors[0][1]

		# If there's no option other than walls, return a bogus row-col pair
		return (32, 32)

	@classmethod
	def wallAt(cls, row: int, col: int) -> bool:
		'''
		Helper function to check if a wall is at a given location
		'''

		# Check if the position is off the grid, and return true if so
		if (row < 0 or row >= 31) or (col < 0 or col >= 28):
			return True

		# Return whether there is a wall at the location
		return bool((wallArr[row] >> col) & 1)


