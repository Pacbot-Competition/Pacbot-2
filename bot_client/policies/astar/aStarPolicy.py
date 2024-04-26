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
def newLocation(row: int, col: int, state: GameState):
	'''
	Construct a new location state
	'''
	result = Location(state)
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
		victimCaught: bool = False,
		targetCaught: bool = False
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
		self.victimCaught: bool = victimCaught

		# Determines whether the target was caught
		self.targetCaught: bool = targetCaught

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
		self.expectedLoc: Location = newLocation(23, 13, self.state)
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

		self.initialSeqComplete = False


	def getNearestPellet(self) -> Location:

		# Check bounds
		first = self.state.pacmanLoc
		if self.state.wallAt(first.row, first.col):
			return self.state.pacmanLoc

		#  BFS traverse
		queue = [first]
		visited = {first.hash()}
		while queue:

			# pop from queue
			currLoc = queue.pop(0)

			# Base Case: Found a pellet
			if self.state.pelletAt(currLoc.row, currLoc.col) and \
				not self.state.superPelletAt(currLoc.row, currLoc.col):
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
				valid = nextLoc.advance() and not self.state.superPelletAt(nextLoc.row, nextLoc.col)

				# avoid same node twice and check this is a valid move
				if nextLoc.hash() not in visited and valid:
					queue.append(nextLoc)
					visited.add(nextLoc.hash())

		#print('No nearest...')
		return first

	def scaryVictim(self, victimColor: GhostColors) -> bool:

		V = self.state.ghosts[victimColor]

		if (victimColor == GhostColors.NONE) or self.state.wallAt(V.location.row, V.location.col):
			return False

		if (V.spawning):
			return True

		for color in GhostColors:
			G = self.state.ghosts[color]
			if (color != victimColor) and (not G.spawning) and (not G.isFrightened()):
				if not self.state.wallAt(G.location.row, G.location.col):
					if self.dist(V.location, G.location) <= 2:
						#print('re-assigning victim')
						return True

		return False

	def getNearestVictim(self) -> GhostColors:

		closest, closestDist = GhostColors.NONE, INF
		for color in GhostColors:
			if self.state.ghosts[color].isFrightened() and not self.state.ghosts[color].spawning:
				dist = self.dist(self.state.pacmanLoc, self.state.ghosts[color].location)
				if dist < closestDist and not self.scaryVictim(color):
					closest = color
					closestDist = dist

		# Return the closest scared ghost
		return closest

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
		if (distTarget == 0):
			return -10000000

		# Dist to nearest scared ghost
		distScared: int = INF
		if victimColor != GhostColors.NONE and not self.state.ghosts[victimColor].spawning:
			distScared = self.dist(self.state.pacmanLoc, self.state.ghosts[victimColor].location)

		# Dist to fruit
		distFruit: int = 999999
		if self.state.fruitSteps > 0:
			distFruit = self.dist(self.state.pacmanLoc, self.state.fruitLoc)

		# Distance to our chosen target: the minimum
		dist: int = distScared if (distScared < INF) else (distTarget if (distTarget < distFruit // 20) else distFruit)
		gCostPerStep: float = 2

		# If the buffer is too small, then the gCostPerStep should be 2 on average
		if bufLen >= 4:
			gCostPerStep = gCost / bufLen

		# Return the result: (g-cost) / (buffer length) * (dist to target)
		return int(gCostPerStep * dist)

	def fCostMultiplier(self) -> float:

		# Constant for the multiplier
		K: int = 64 # TODO: considering changing to 128

		# Multiplier addition term
		multTerm: int = 0

		# Location of the ghost lair
		lairLoc: Location = Location(self.state)
		lairLoc.row = 11
		lairLoc.col = 13

		# Check if any ghosts are frightened
		fright = False
		for ghost in self.state.ghosts:
			if ghost.isFrightened():
				fright = True

		# Calculate closest non-frightened ghost
		for ghost in self.state.ghosts:
			if not ghost.spawning:
				if not ghost.isFrightened():
					multTerm += int(
						K >> self.dist(self.state.pacmanLoc, ghost.location)
					)
			elif not fright and not self.state.wallAt(self.state.pacmanLoc.row, self.state.pacmanLoc.col):
				multTerm += int(
					K >> self.dist(self.state.pacmanLoc, lairLoc)
				)

		# Return the multiplier (1 + constant / distance squared)
		return 1 + multTerm

	def selectTarget(self, pelletTarget: Location) -> None:

		chase = self.state.gameMode == GameModes.CHASE

		if self.initialSeqComplete:
			pass

		# check if top left pellet exists
		if self.state.superPelletAt(3, 1) and chase:
			self.target = newLocation(5, 1, self.state)

		# check if top right pellet exists
		elif self.state.superPelletAt(3, 26) and chase:
			self.target = newLocation(5, 26, self.state)

		# check if bottom left pellet exists
		elif self.state.superPelletAt(23, 1) and chase:
			self.target = newLocation(20, 3, self.state)

		# check if bottom right pellet exists
		elif self.state.superPelletAt(23, 26) and chase:
			self.target = newLocation(20, 24, self.state)

		# no super pellets
		else:
			# target the nearest pellet
			self.target = pelletTarget

	async def act(self, predicted_delay: int, victimColor: GhostColors, pelletTarget: Location) -> tuple[GhostColors, Location]:

		# Make a priority queue of A-Star Nodes
		priorityQueue: list[AStarNode] = []

		# Starting row and col
		startRow = self.state.pacmanLoc.row
		startCol = self.state.pacmanLoc.col

		# Construct an initial node
		initialNode = AStarNode(
			compressGameState(self.state),
			fCost = self.hCostExtend(0, 0, victimColor),
			gCost = 0,
			directionBuf = [],
			delayBuf = [],
			bufLength = 0
		)

		# Add the initial node to the priority queue
		heappush(priorityQueue, initialNode)

		# Select a target
		self.selectTarget(pelletTarget)

		# Select a victim, if applicable
		victimCaught = (victimColor != GhostColors.NONE) and ((not self.state.ghosts[victimColor].isFrightened()) or self.state.ghosts[victimColor].spawning)
		if victimColor == GhostColors.NONE or self.scaryVictim(victimColor) or victimCaught:
			victimColor = self.getNearestVictim()

		# Select a new target, if applicable
		targetCaught = self.state.wallAt(pelletTarget.row, pelletTarget.col) or not self.state.pelletAt(pelletTarget.row, pelletTarget.col) or self.state.fruitLoc.at(self.state.pacmanLoc.row, self.state.pacmanLoc.col)
		if targetCaught:
			pelletTarget = self.getNearestPellet()

		# Flag for first iteration
		firstIt = True

		# Lag for first iteration
		firstItLag = 0

		# Lag for turns
		turnLag = 10

		# Keep proceeding until a break point is hit
		while len(priorityQueue):

			# Pop the lowest f-cost node
			currNode = heappop(priorityQueue)

			# Reset to the current compressed state
			decompressGameState(self.state, currNode.compressedState)

			# If the g-cost of this node is high enough or we reached the target,
			# make the moves and return
			if currNode.bufLength >= 14:

				# testLoc = newLocation(startRow, startCol, self.state)
				# for index in range(1):
				# 	testLoc.setDirection(currNode.directionBuf[index])
				# 	testLoc.advance()

				# 	self.state.queueAction(
				# 		currNode.delayBuf[index] - (index == 0),
				# 		currNode.directionBuf[index],
				# 		testLoc.row,
				# 		testLoc.col
				# 	)

				testLoc = newLocation(startRow, startCol, self.state)
				lastDir = currNode.directionBuf[0]
				dist = 0

				for index in range(len(currNode.directionBuf)):
					# coalesce
					if lastDir != currNode.directionBuf[index]:
						break
					dist += 1

					# get target location
					testLoc.setDirection(currNode.directionBuf[index])
					testLoc.advance()


				self.state.queueAction(
					currNode.delayBuf[0] - (0 == 0),
					lastDir,
					dist,
					testLoc.row,
					testLoc.col,
				)

				#print(['RED', 'PINK', 'CYAN', 'ORANGE', 'NONE'][victimColor], pelletTarget)
				return victimColor, pelletTarget

			if currNode.victimCaught:
				# testLoc = newLocation(startRow, startCol, self.state)
				# lastDir = currNode.directionBuf[0]
				# dist = 0

				# for index in range(len(currNode.directionBuf)):
				# 	# coalesce
				# 	if lastDir != currNode.directionBuf[index]:
				# 		break
				# 	dist += 1

				# 	# get target location
				# 	testLoc.setDirection(currNode.directionBuf[index])
				# 	testLoc.advance()


				# self.state.queueAction(
				# 	currNode.delayBuf[0] - (0 == 0),
				# 	lastDir,
				# 	dist,
				# 	testLoc.row,
				# 	testLoc.col,
				# )

				#print('victim caught?')

				# force pb to take this path
				priorityQueue.clear()

				if currNode.targetCaught:
					#print('target caught')
					pelletTarget = self.getNearestPellet()

				#print(['RED', 'PINK', 'CYAN', 'ORANGE', 'NONE'][victimColor], pelletTarget)
				# return victimColor, pelletTarget

			elif currNode.targetCaught and (victimColor == GhostColors.NONE):

				# force pb to take this path
				priorityQueue.clear()

				# choose new target
				pelletTarget = self.getNearestPellet()


				# testLoc = newLocation(startRow, startCol, self.state)
				# lastDir = currNode.directionBuf[0]
				# dist = 0

				# for index in range(len(currNode.directionBuf)):
				# 	# coalesce
				# 	if lastDir != currNode.directionBuf[index]:
				# 		break
				# 	dist += 1

				# 	# get target location
				# 	testLoc.setDirection(currNode.directionBuf[index])
				# 	testLoc.advance()


				# self.state.queueAction(
				# 	currNode.delayBuf[0] - (0 == 0),
				# 	lastDir,
				# 	dist,
				# 	testLoc.row,
				# 	testLoc.col,
				# )

				# #print('target caught')
				# pelletTarget = self.getNearestPellet()

				# #print(['RED', 'PINK', 'CYAN', 'ORANGE', 'NONE'][victimColor], pelletTarget)
				# return GhostColors.NONE, pelletTarget


			# Get Pacman's current direction
			prevDir = self.state.pacmanLoc.getDirection()

			# Determines if waiting (none) is allowed as a move
			waitAllowed = (victimColor == GhostColors.NONE)

			# TODO: EVALUATE THIS

			# Loop over the directions
			for direction in Directions:

				# If direction is none, continue
				if (direction == Directions.NONE) and (not waitAllowed):
					continue

				# Reset to the current compressed state
				decompressGameState(self.state, currNode.compressedState)

				turnPenalty = 0
				evadePenalty = 0
				if (prevDir != direction):
					turnPenalty = 2

					if (victimColor != GhostColors.NONE) and not self.state.ghosts[victimColor].spawning:

						loc: Location = Location(self.state)
						loc.update(self.state.pacmanLoc.serialize())
						loc.setDirection(direction)
						dist1 = self.dist(loc, self.state.ghosts[victimColor].location)
						loc.advance()
						dist2 = self.dist(loc, self.state.ghosts[victimColor].location)

						if (dist1 < dist2):
							evadePenalty = 10

				npBefore = self.state.numPellets()
				nspBefore = self.state.numSuperPellets()
				valid = self.state.simulateAction(predicted_delay + firstItLag * firstIt + turnPenalty * turnLag, direction)
				npAfter = self.state.numPellets()
				nspAfter = self.state.numSuperPellets()
				ateNormalPellet = (npBefore > npAfter) and (nspBefore == nspAfter)

				# Determines if the target was caught
				targetCaught = self.state.wallAt(pelletTarget.row, pelletTarget.col) or (not self.state.pelletAt(pelletTarget.row, pelletTarget.col)) or pelletTarget.at(self.state.pacmanLoc.row, self.state.pacmanLoc.col) or (self.state.fruitLoc.at(self.state.pacmanLoc.row, self.state.pacmanLoc.col))

				# Determines if the scared ghost 'victim' was caught
				victimCaught = (victimColor != GhostColors.NONE) and ((not self.state.ghosts[victimColor].isFrightened()) or self.state.ghosts[victimColor].spawning)

				# Select a new target
				self.selectTarget(pelletTarget)

				# Determine if there is a frightened ghost to chase
				victimExists = (victimColor == GhostColors.NONE)

				# If the state is valid, add it to the priority queue
				if valid:
					nextNode = AStarNode(
						compressGameState(self.state),
						fCost = int((self.hCostExtend(currNode.gCost, currNode.bufLength, victimColor) + currNode.gCost + 1) * self.fCostMultiplier()),
						gCost = currNode.gCost + 2 + 4 * ((not ateNormalPellet) and (not victimExists)) + 2 * (turnPenalty and victimExists) + 5 * evadePenalty,
						directionBuf = currNode.directionBuf + [direction],
						delayBuf = currNode.delayBuf + [predicted_delay + firstItLag * firstIt + turnPenalty * turnLag],
						bufLength = currNode.bufLength + 1,
						victimCaught = victimCaught,
						targetCaught = targetCaught
					)

					# Add the next node to the priority queue
					heappush(priorityQueue, nextNode)

			firstIt = False

		#print("Trapped...")

		return victimColor, pelletTarget