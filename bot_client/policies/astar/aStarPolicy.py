# Heap Queues
from heapq import heappush, heappop

# Game state
from gameState import *

# Location mapping
import policies.astar.genPachattanDistDict as pacdist
import policies.astar.example as ex

# Big Distance
INF = 999999

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
'''

class DistTypes(IntEnum):
	'''
	Enum of distance types
	'''
	MANHATTAN_DISTANCE = 0
	EUCLIDEAN_DISTANCE = 1
	PACHATTAN_DISTANCE = 2

# Create new location with row, col
def newLocation(row: int, col: int):
	'''
	Construct a new location state
	'''
	result = Location(0)
	result.row = row
	result.col = col
	return result

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
		victimColor: GhostColors
	) -> None:

		# Compressed game state
		self.compressedState = compressedState

		# Costs
		self.fCost = fCost
		self.gCost = gCost

		# Estimated velocity
		self.estSpeed = 0
		self.direction = Directions.NONE

		# Message buffer
		self.directionBuf = directionBuf
		self.delayBuf = delayBuf
		self.bufLength = bufLength

		# Victim color (catching scared ghosts)
		self.victimColor: GhostColors = GhostColors.NONE

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

	def hCost(self) -> int:
		# make sure pacman in bounds (TODO: Why do we have to do this?)
		if 0 > self.state.pacmanLoc.row or 32 <= self.state.pacmanLoc.row or 0 > self.state.pacmanLoc.col or 28 <= self.state.pacmanLoc.col:
			return 999999999

		# Heuristic cost for this location
		hCostTarget = 0

		# Heuristic cost to estimate ghost locations
		hCostGhost = 0

		# Catching frightened ghosts
		hCostScaredGhost = 0

		# Chasing fruit
		hCostFruit = 0

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
					hCostScaredGhost += self.dist(self.state.pacmanLoc, ghost.location)

		# Check whether fruit exists, and then add it to target
		if self.state.fruitSteps > 0:
			hCostFruit = self.dist(self.state.pacmanLoc, self.state.fruitLoc)

		# If there are frightened ghosts, chase them
		if hCostScaredGhost > 0:
			return int(hCostTarget + hCostGhost + hCostScaredGhost + hCostFruit)

		# Otherwise, chase the target
		hCostTarget = self.dist(self.state.pacmanLoc, self.target)
		return int(hCostTarget + hCostGhost + hCostFruit)

	def hCostExtend(self, gCost: int, bufLen: int, victimColor: GhostColors) -> int:
		'''
		Extends the existing g_cost delta to estimate a new h-cost due to
		distance Pachattan distance and estimated speed
		'''

		# make sure pacman in bounds
		if 0 > self.state.pacmanLoc.row or 32 <= self.state.pacmanLoc.row or 0 > self.state.pacmanLoc.col or 28 <= self.state.pacmanLoc.col:
			return 999999999

		# Dist to target
		distTarget: int = self.dist(self.state.pacmanLoc, self.target)

		# Dist to nearest scared ghost
		distScared: int = INF
		if victimColor != GhostColors.NONE:
			distScared = self.dist(self.state.pacmanLoc, self.state.ghosts[victimColor].location)

		# Dist to fruit
		distFruit: int = 999999
		if self.state.fruitSteps > 0:
			distFruit = self.dist(self.state.pacmanLoc, self.state.fruitLoc)

		# Distance to our chosen target: the minimum
		dist: int = distScared if (distScared < INF) else min(distTarget, distFruit)
		gCostPerStep: float = 3

		# If the buffer is too small, then the gCostPerStep should be 2 on average
		if bufLen >= 3:
			gCostPerStep = gCost / bufLen

		# Return the result: (g-cost) / (buffer length) * (dist to target)
		return int(gCostPerStep * dist)

	def fCostMultiplier(self) -> float:

		# Constant for the multiplier
		K: int = 8

		# Multiplier addition term
		multTerm: int = 0

		# Calculate closest non-frightened ghost
		for ghost in self.state.ghosts:
			if not ghost.spawning:
				if not ghost.isFrightened():
					multTerm += int(
						K >> self.dist(self.state.pacmanLoc, ghost.location)
					)

		# Return the multiplier (1 + constant / distance squared)
		return 1 + multTerm

	async def act(self, predicted_delay: int = 6, victimColor: GhostColors = GhostColors.NONE) -> GhostColors:

		# Make a priority queue of A-Star Nodes
		priorityQueue: list[AStarNode] = []

		# Choose a scared ghost to attack
		if victimColor == GhostColors.NONE:
			closest, closestDist = GhostColors.NONE, INF
			for color in GhostColors:
				if self.state.ghosts[color].isFrightened() and not self.state.ghosts[color].spawning:
					dist = self.dist(self.state.pacmanLoc, self.state.ghosts[color].location)
					if dist < closestDist:
						closest = color
						closestDist = dist
			victimColor = closest

		# Reset the victim color
		elif not self.state.ghosts[victimColor].isFrightened() or self.state.ghosts[victimColor].spawning:
			victimColor = GhostColors.NONE

		# Construct an initial node
		initialNode = AStarNode(
			compressGameState(self.state),
			fCost = self.hCostExtend(0, 0, victimColor),
			gCost = 0,
			directionBuf = [],
			delayBuf = [],
			bufLength = 0,
			victimColor = victimColor
		)

		# Add the initial node to the priority queue
		heappush(priorityQueue, initialNode)

		# check if top right pellet exists
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
			# target the nearest pellet
			self.target = self.getNearestPellet()

		# Keep proceeding until a break point is hit
		while len(priorityQueue):

			# Pop the lowest f-cost node
			currNode = heappop(priorityQueue)

			# Reset to the current compressed state
			decompressGameState(self.state, currNode.compressedState)

			# If the g-cost of this node is high enough or we reached the target,
			# make the moves and return
			if currNode.bufLength >= 8 or (currNode.victimColor != GhostColors.NONE and currNode.bufLength >= 4):
				for index in range(min(currNode.bufLength, 4)):
					self.state.queueAction(
						currNode.delayBuf[index],
						currNode.directionBuf[index]
					)
				return currNode.victimColor

			# Loop over the directions
			for direction in Directions:

				# Reset to the current compressed state
				decompressGameState(self.state, currNode.compressedState)

				valid = self.state.simulateAction(predicted_delay, direction)

				# Make the scared ghost 'victim' color the same as the current node
				victimColor = currNode.victimColor

				# Choose a scared ghost to attack
				if victimColor == GhostColors.NONE:
					closest, closestDist = GhostColors.NONE, INF
					for color in GhostColors:
						if self.state.ghosts[color].isFrightened() and not self.state.ghosts[color].spawning:
							dist = self.dist(self.state.pacmanLoc, self.state.ghosts[color].location)
							if dist < closestDist:
								closest = color
								closestDist = dist
					victimColor = closest

				# If the state is valid, add it to the priority queue
				if valid:
					nextNode = AStarNode(
						compressGameState(self.state),
						fCost = int((self.hCostExtend(currNode.gCost, currNode.bufLength, victimColor) + currNode.gCost + 1) * self.fCostMultiplier()),
						gCost = currNode.gCost + 1,
						directionBuf = currNode.directionBuf + [direction],
						delayBuf = currNode.delayBuf + [predicted_delay],
						bufLength = currNode.bufLength + 1,
						victimColor = victimColor
					)

					# Add the next node to the priority queue
					heappush(priorityQueue, nextNode)

		return victimColor