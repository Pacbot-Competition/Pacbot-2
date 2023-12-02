# Heap Queues
from heapq import heappush, heappop

# Game state
from gameState import *

# Location mapping
import policies.astar.genPachattanDistDict as pacdist
import policies.astar.example as ex



'''
Cost Explanations:

Started at point	S
Targetting point 	T
Currently at point 	C

gcost = cost from S to C (past, known)
hcost = cost from C to T (future, predicted)

fcost = gcost + hcost

Start-------Current-------Target
S--------------C---------------T
|-----gcost----|-----hcost-----|
|------------fcost-------------|

Test D3:
4720
1330
3640

'''

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

# Pachattan distance
def distL3(loc1: Location, loc2: Location) -> int:
	key = pacdist.getKey(loc1, loc2)
	return ex.PACHATTAN[key]

# Squared Pachattan distance
def distSqL3(loc1: Location, loc2: Location) -> int:
	pacDist = distL3(loc1, loc2)
	return pacDist * pacDist


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

		if 0 > self.state.pacmanLoc.row or 32 <= self.state.pacmanLoc.row or 0 > self.state.pacmanLoc.col or 28 <= self.state.pacmanLoc.col:
			return 999999999

		# Heuristic cost for this location
		hCostTarget = distL3(self.state.pacmanLoc, self.target)

		# Heuristic cost to estimate ghost locations
		hCostGhost = 0

		# Catching frightened ghosts
		hCostScaredGhost = 999999999

		# Chasing fruit
		hCostFruit = 999999999

		# Add a penalty for being close to the ghosts
		for ghost in self.state.ghosts:
			if not ghost.spawning:
				if not ghost.isFrightened():
					hCostGhost += int(
						64 / max(distSqL3(
							self.state.pacmanLoc,
							ghost.location
						), 1)
					)
				else:
					hCostScaredGhost = min(
						distL3(self.state.pacmanLoc, ghost.location),
						hCostScaredGhost
					)

		# Check whether fruit exists, and add a target to it if so
		if self.state.fruitSteps > 0:
			hCostFruit = distL3(self.state.pacmanLoc, self.state.fruitLoc)

		# If there are frightened ghosts, chase them
		if hCostScaredGhost < 999999999:
			return min(hCostScaredGhost, hCostFruit) + hCostGhost

		# Otherwise, if there is a fruit on the board, target fruit
		if hCostFruit != 999999999:
			return hCostFruit + hCostGhost
		
		# Otherwise, chase the target
		return hCostTarget + hCostGhost

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

		if self.state.superPelletAt(3, 26):
			self.target = newLocation(5, 21)

        # check if top left pellet exists
		elif self.state.superPelletAt(3, 1):
			self.target = newLocation(5, 6)

        # check if bottom left pellet exists
		elif self.state.superPelletAt(23, 1):
			self.target = newLocation(20, 1)

        # check if bottom right pellet exists
		elif self.state.superPelletAt(23, 26):
			self.target = newLocation(20, 26)

		# Keep proceeding until a break point is hit
		while len(priorityQueue):

			# Pop the lowest f-cost node
			currNode = heappop(priorityQueue)

			# Reset to the current compressed state
			decompressGameState(self.state, currNode.compressedState)

			# If the g-cost of this node is high enough or we reached the target,
			# make the moves and return
			if currNode.bufLength >= 8 or self.hCost() <= 1:
				for index in range(min(2, currNode.bufLength)):
					# TODO: Avoid sending same command twice. Sometimes, there is a delay
					# before algo can recognize that an instruction has been executed,
					# so it might send the same action multiple times.
					self.state.queueAction(
						currNode.delayBuf[index],
						currNode.directionBuf[index]
					)

				print('decided')
				return

			# Loop over the directions
			for direction in Directions:

				# If the direction is none, skip it
				if direction == Directions.NONE:
					continue

				# Reset to the current compressed state
				decompressGameState(self.state, currNode.compressedState)

				# Check whether the direction is valid
				valid = self.state.simulateAction(6, direction)

				# If the state is valid, add it to the priority queue
				if valid:

					nextNode = AStarNode(
						compressGameState(self.state),
						fCost = currNode.gCost + 1 + self.hCost(),
						gCost = currNode.gCost + 1,
						directionBuf = currNode.directionBuf + [direction],
						delayBuf = currNode.delayBuf + [6],
						bufLength = currNode.bufLength + 1
					)

					# Add the next node to the priority queue
					heappush(priorityQueue, nextNode)