# Heap Queues
from heapq import heappush, heappop
import math

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

class DistTypes(IntEnum):
	'''
	Enum of distance types
	'''
	MANHATTAN_DISTANCE = 0
	EUCLIDEAN_DISTANCE = 1
	PACHATTAN_DISTANCE = 2

# INNER_RING_LOCATIONS = {
# 	"9,12","9,15","10,12","10,15",
# 	"14,7","14,8","14,19","14,20",
# 	"18,9","19,9","18,18","19,18",
# 	"11,9","11,10","11,11","11,12","11,13","11,14","11,15","11,16","11,17","11,18",
# 	"12,9","13,9","14,9","15,9","16,9",
# 	"12,18","13,18","14,18","15,18","16,18",
# 	"17,9","17,10","17,11","17,12","17,13","17,14","17,15","17,16","17,17","17,18",
# }




# Create new location with row, col
def newLocation(row: int, col: int):
	'''
	Construct a new location state
	'''
	result = Location(0)
	result.row = row
	result.col = col
	return result

GHOST_SPAWN_LOCATION = newLocation(11, 13) # in between two squares so really 13.5


# Manhattan distance
def distL1(loc1: Location, loc2: Location) -> int:
	return abs(loc1.row - loc2.row) + abs(loc1.col - loc2.col)

# Manhattan distance
def distSqL1(loc1: Location, loc2: Location) -> int:
	dr = abs(loc1.row - loc2.row)
	dc = abs(loc1.col - loc2.col)
	return dr*dr + dc*dc

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
		bufLength: int,
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



class AStarPolicy:
	'''
	Policy class for running the A-Star Algorithm for Pacbot.
	'''

	def __init__(
		self,
		state: GameState,
		target: Location,
		distType: DistTypes = DistTypes.PACHATTAN_DISTANCE
	) -> None:

		# Game state
		self.state: GameState = state
		self.stateCopy: GameState = state

		# Target location
		self.target: Location = target

		# Expected location
		self.expectedLoc: Location = newLocation(23, 13)
		self.error_sum = 0
		self.error_count = 0
		self.dropped_command_count = 0

		# Distance metrics
		self.distType = distType
		match self.distType:
			case DistTypes.MANHATTAN_DISTANCE:
				self.dist = distL1
				self.distSq = distSqL1
			case DistTypes.EUCLIDEAN_DISTANCE:
				self.dist = distL2
				self.distSq = distSqL2
			case DistTypes.PACHATTAN_DISTANCE:
				self.dist = distL3
				self.distSq = distSqL3
			case _: # pachattan
				self.distType = DistTypes.PACHATTAN_DISTANCE
				self.dist = distL3
				self.distSq = distSqL3

	def getNearestPellet(self) -> Location:

		# Check bounds
		first = self.state.pacmanLoc
		if self.state.wallAt(first.row, first.col):
			return self.state.pacmanLoc

		#  BFS traverse
		queue = [first]
		visited = set(str(first))
		while queue:

			# pop from queue
			currLoc = queue.pop(0)

			# Base Case: Found a pellet
			if self.state.pelletAt(currLoc.row, currLoc.col):
				return currLoc

			# Loop over the directions
			for direction in Directions:

				# If the direction is none, skip it
				if direction == Directions.NONE:
					continue

				# Increment direction
				nextLoc = Location(self.state)
				nextLoc.col = currLoc.col
				nextLoc.row = currLoc.row
				nextLoc.setDirection(direction)
				valid = nextLoc.advance()

				# avoid same node twice and check this is a valid move
				if str(nextLoc) not in visited and valid:
					queue.append(nextLoc)
					visited.add(str(nextLoc))
		return first
		
	def pelletAtSafe(self, row, col):
		if self.state.wallAt(row, col):
			return False
		return self.state.pelletAt(row, col)

	def pelletTarget(self) -> Location:
		# calc vectors from every pellet in quad to pacman
		# target median pellet
		return Location(self.state)


	def hCost(self, realPacLoc, pelletExists=False) -> int:

		if 0 > self.state.pacmanLoc.row or 32 <= self.state.pacmanLoc.row or 0 > self.state.pacmanLoc.col or 28 <= self.state.pacmanLoc.col:
			return 999999999

		# Heuristic cost for this location
		hCostTarget = 0

		# Heuristic cost to estimate ghost locations
		hCostGhost = 0

		# Catching frightened ghosts
		# hCostScaredGhost = 999999999
		hCostScaredGhost = 0

		# Chasing fruit
		hCostFruit = 0
		
		# Pellet heuristic
		hCostPellet = 1

		# Ghost Spawn heuristic
		hCostGhostSpawn = 0

		# Add a penalty for being close to the ghosts
		for ghost in self.state.ghosts:
			if not ghost.spawning:
				if not ghost.isFrightened():
					hCostGhost += int(
						64 / max(self.distSq(
							self.state.pacmanLoc,
							ghost.location
						), 1)
					)
				else:
					# hCostScaredGhost = min(
					# 	# 0.25*math.pow(self.dist(self.state.pacmanLoc, ghost.location), 1.1),
					# 	# 5*self.dist(self.state.pacmanLoc, ghost.location),
					# 	# 10*math.log(self.dist(self.state.pacmanLoc, ghost.location)),
					# 	self.dist(self.state.pacmanLoc, ghost.location),
					# 	hCostScaredGhost
					# )
					hCostScaredGhost += self.dist(self.state.pacmanLoc, ghost.location)
					# row = realPacLoc.row - ghost.location.row
					# if row != 0: row /= abs(row)
					# col = realPacLoc.col - ghost.location.col
					# if col != 0: col /= abs(col)
					# offsetLoc = newLocation(int(ghost.location.row + row), int(ghost.location.col + col))
					# if not self.state.wallAt(offsetLoc.row, offsetLoc.col):
					# 	hCostScaredGhost += self.dist(self.state.pacmanLoc, offsetLoc)
					# else:
					# 	hCostScaredGhost += self.dist(self.state.pacmanLoc, ghost.location)
			# else:
			# 	# Check if ghost spawning
			# 	hCostGhostSpawn = int(2 / max(self.distSq(self.state.pacmanLoc, GHOST_SPAWN_LOCATION), 1))

		# Check whether fruit exists, and add a target to it if so
		if self.state.fruitSteps > 0:
			# hCostFruit = math.pow(self.dist(self.state.pacmanLoc, self.state.fruitLoc), 1.1)
			# hCostFruit = 5*self.dist(self.state.pacmanLoc, self.state.fruitLoc)
			# hCostFruit = 10 * math.log(self.dist(self.state.pacmanLoc, self.state.fruitLoc))
			hCostFruit = self.dist(self.state.pacmanLoc, self.state.fruitLoc)

		# Check if pellet at current node
		self.state.pelletAt(row=self.state.pacmanLoc.row, col=self.state.pacmanLoc.col)

		# Compute hCostPellet
		if pelletExists:
			hCostPellet = 0 # TODO: do something more sophisticated than adding a constant

		# If there are frightened ghosts, chase them
		# if hCostScaredGhost < 999999999:
		if hCostScaredGhost > 0:
			# return int(hCostTarget + min(hCostScaredGhost, hCostFruit) + hCostGhost + hCostGhostSpawn)
			return int(hCostTarget + hCostGhost + hCostScaredGhost + hCostFruit + hCostPellet + hCostGhostSpawn)

		# Otherwise, if there is a fruit on the board, target fruit
		# if hCostFruit != 0:
		# 	return int(hCostTarget + hCostFruit + hCostGhost + hCostGhostSpawn)
		
		# Otherwise, chase the target
		hCostTarget = self.dist(self.state.pacmanLoc, self.target)
		return int(hCostTarget + hCostGhost + hCostFruit + hCostPellet + hCostGhostSpawn)

	async def act(self, predicted_delay=6) -> None:

		# Make a priority queue of A-Star Nodes
		priorityQueue: list[AStarNode] = []

		# Construct an initial node
		initialNode = AStarNode(
			compressGameState(self.state),
			fCost = self.hCost(self.state.pacmanLoc),
			gCost = 0,
			directionBuf = [],
			delayBuf = [],
			bufLength = 0
		)

		# counter to avoid calc nearest pellet a bunch of times
		counter = 1

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
		# no super pellets
		else:
			# avoid calc every time (wait 20 decisions)
			if counter == 1 or \
					self.target.row == self.state.pacmanLoc.row and \
					self.target.col == self.state.pacmanLoc.col:
				self.target = self.getNearestPellet()
				counter = 0
			else:
				counter += 1
		
		
		print("-"*15)
		print("expected: " + str(self.expectedLoc))
		if str(self.expectedLoc) != str(self.state.pacmanLoc):
			print("actual: " + str(self.state.pacmanLoc) + " - non match! (Expected " + str(self.expectedLoc) + ")")
		else:
			print("actual: " + str(self.state.pacmanLoc))
		origLoc = newLocation(self.state.pacmanLoc.row, self.state.pacmanLoc.col)
		origLoc.state = self.state


		self.error_sum += distL1(origLoc, self.expectedLoc)
		self.error_count += 1
		# self.dropped_command_count += distL3(origLoc, self.expectedLoc) # not a perfect measure
		print("average error: " + str(self.error_sum/self.error_count))

		# print("dropped command count: " + str(self.dropped_command_count))
				
		# self.state.pacmanLoc.row = self.expectedLoc.row
		# self.state.pacmanLoc.col = self.expectedLoc.col

		realPacLoc = self.state.pacmanLoc
		

		# Keep proceeding until a break point is hit
		while len(priorityQueue):

			# Pop the lowest f-cost node
			currNode = heappop(priorityQueue)

			# Reset to the current compressed state
			decompressGameState(self.state, currNode.compressedState)

			# If the g-cost of this node is high enough or we reached the target,
			# make the moves and return
			if currNode.bufLength >= 10 or self.hCost(self.state.pacmanLoc) <= 1:
				for index in range(min(2, currNode.bufLength)):
					self.state.queueAction(
						currNode.delayBuf[index],
						currNode.directionBuf[index]
					)
					origLoc.setDirection(currNode.directionBuf[index])
					origLoc.advance()
					print(currNode.directionBuf[index])


				self.expectedLoc = origLoc
				return
			
			# get current direction (we will use this to negatively weight changing directions)
			currDir = self.state.pacmanLoc.getDirection()

			# Loop over the directions
			for direction in Directions:

				# Reset to the current compressed state
				decompressGameState(self.state, currNode.compressedState)

				# TODO: Fix failing when pacbot dies
				# Check if there's a pellet at curr location + direction
				if direction == Directions.UP:
					pelletExists = self.pelletAtSafe(row=self.state.pacmanLoc.row - 1, col=self.state.pacmanLoc.col)
				elif direction == Directions.LEFT:
					pelletExists = self.pelletAtSafe(row=self.state.pacmanLoc.row, col=self.state.pacmanLoc.col - 1)
				elif direction == Directions.DOWN:
					pelletExists = self.pelletAtSafe(row=self.state.pacmanLoc.row + 1, col=self.state.pacmanLoc.col)
				elif direction == Directions.RIGHT:
					pelletExists = self.pelletAtSafe(row=self.state.pacmanLoc.row, col=self.state.pacmanLoc.col + 1)
				else:
					pelletExists = self.pelletAtSafe(row=self.state.pacmanLoc.row, col=self.state.pacmanLoc.col)

				# Check whether the direction is valid
				valid = self.state.simulateAction(predicted_delay, direction)
				
				# If the state is valid, add it to the priority queue
				if valid:
					# calculate the cost of changing direction
					changeDirCost = 0 if (currDir == self.state.pacmanLoc.getDirection()) else 4

					nextNode = AStarNode(
						compressGameState(self.state),
						fCost = currNode.gCost + 1 + changeDirCost + self.hCost(realPacLoc, pelletExists),
						gCost = currNode.gCost + 1 + changeDirCost,
						directionBuf = currNode.directionBuf + [direction],
						delayBuf = currNode.delayBuf + [predicted_delay],
						bufLength = currNode.bufLength + 1
					)

					# Add the next node to the priority queue
					heappush(priorityQueue, nextNode)