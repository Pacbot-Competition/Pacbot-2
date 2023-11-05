# Heap Queues
from heapq import heappush, heappop

# Game state
from gameState import *

'''
Distance metrics
'''

# Manhattan distance
def distL1(loc1: Location, loc2: Location) -> int:
	return abs(loc1.row - loc2.row) + abs(loc1.col - loc2.col)

# Squared Euclidean distance
def distSqL2(loc1: Location, loc2: Location) -> int:
	return (loc1.row - loc2.row) * (loc1.row - loc2.row) + \
		(loc1.col - loc2.col) * (loc1.col - loc2.col)

# Euclidean distance
def distL2(loc1: Location, loc2: Location) -> int:
	return ((loc1.row - loc2.row) * (loc1.row - loc2.row) + \
		(loc1.col - loc2.col) * (loc1.col - loc2.col)) ** 0.5

class AStarNode:
	'''
	Node class for running the A-Star Algorithm for Pacbot.
	'''

	def __init__(
		self,
		compressedState: GameStateCompressed,
		fCost: int,
		gCost: int,
		directionBuf: list[Directions],
		delayBuf: list[int],
		bufLength: int
	) -> None:

		# Compressed game state
		self.compressedState = compressedState

		# Costs
		self.fCost = fCost
		self.gCost = gCost

		# Message buffer
		self.directionBuf = directionBuf
		self.delayBuf = delayBuf
		self.bufLength = bufLength

	def __lt__(self, other) -> bool: # type: ignore
		return self.fCost < other.fCost # type: ignore

	def __repr__(self) -> str:
		return str(f'g = {self.gCost} ~ f = {self.fCost}')

def newLocation(row: int, col: int):
	'''
	Construct a new location state
	'''
	result = Location(0)
	result.row = row
	result.col = col
	return result

class AStarPolicy:
	'''
	Policy class for running the A-Star Algorithm for Pacbot.
	'''

	def __init__(
		self,
		state: GameState,
		target: Location
	) -> None:

		# Game state
		self.state: GameState = state

		# Target location
		self.target: Location = target

	def hCost(self) -> int:

		# Return the heuristic cost for this
		return distL2(self.state.pacmanLoc, self.target)

	async def act(self) -> None:

		# Make a priority queue of A-Star Nodes
		priorityQueue: list[AStarNode] = []

		# Construct an initial node
		initialNode = AStarNode(
			compressGameState(self.state),
			fCost = self.hCost(),
			gCost = 0,
			directionBuf = [],
			delayBuf = [],
			bufLength = 0
		)

		# Add the initial node to the priority queue
		heappush(priorityQueue, initialNode)

		# Keep proceeding until a break point is hit
		while len(priorityQueue):

			# Pop the lowest f-cost node
			currNode = heappop(priorityQueue)

			# Reset to the current compressed state
			decompressGameState(self.state, currNode.compressedState)

			# If the g-cost of this node is high enough or we reached the target,
			# make the moves and return
			if len(currNode.directionBuf) >= 10 or \
				distSqL2(self.state.pacmanLoc, self.target) == 0:

				for index in range(len(currNode.directionBuf)):
					self.state.queueAction(
						currNode.delayBuf[index],
						currNode.directionBuf[index]
					)

				return

			# Loop over the directions
			for direction in Directions:

				# If the direction is none, skip it
				if direction == Directions.NONE:
					continue

				# Reset to the current compressed state
				decompressGameState(self.state, currNode.compressedState)

				# Check whether the direction is valid
				valid = self.state.simulateAction(3, direction)

				# If the state is valid, add it to the priority queue
				if valid:

					nextNode = AStarNode(
						compressGameState(self.state),
						fCost = currNode.gCost + 1 + self.hCost(),
						gCost = currNode.gCost + 1,
						directionBuf = currNode.directionBuf + [direction],
						delayBuf = currNode.delayBuf + [3],
						bufLength = currNode.bufLength + 1
					)

					# Add the next node to the priority queue
					heappush(priorityQueue, nextNode)