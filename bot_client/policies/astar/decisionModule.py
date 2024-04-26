# JSON (for reading config.json)
import json

# Asyncio (for concurrency)
import asyncio

# Game state
from gameState import *

# A-Star Policy
from policies.astar.aStarPolicy import *

# Get the FPS of the server from the config.json file
def getGameFPS() -> int:

	# Read the configuration file
	with open('../config.json', 'r', encoding='UTF-8') as configFile:
		config = json.load(configFile)

	# Return the FPS
	return config["GameFPS"]

class DecisionModule:
	'''
	Sample implementation of a decision module for high-level
	programming for Pacbot, using asyncio.
	'''

	def __init__(self, state: GameState) -> None:
		'''
		Construct a new decision module object
		'''

		# Game state object to store the game information
		self.state = state

		# Policy object, with the game state
		self.policy = AStarPolicy(state, newLocation(5, 21, self.state))

	async def decisionLoop(self) -> None:
		'''
		Decision loop for Pacbot
		'''

		# wait = True
		# gameFPS = getGameFPS()
		victimColor = GhostColors.NONE
		pelletTarget = Location(self.state)
		pelletTarget.row = 23
		pelletTarget.col = 6 # start by moving to the left??

		# Receive values as long as we have access
		while self.state.isConnected():

			'''
			WARNING: 'await' statements should be routinely placed
			to free the event loop to receive messages, or the
			client may fall behind on updating the game state!
			'''

			# If the current messages haven't been sent out yet, skip this iteration
			if not self.state.isFound() or self.state.isLocked() or self.state.isPaused():
				await asyncio.sleep(0)
				continue

			# Lock the game state
			self.state.lock()

			# Figure out which actions to take, according to the policy

			print("astar calculating...")
			victimColor, pelletTarget = await self.policy.act(4, victimColor, pelletTarget)
			print("astar done")

			if not len(self.state.writeServerBuf):
				print(f"{YELLOW}astar failed, trying again{NORMAL}")
				await asyncio.sleep(0)
				self.state.unlock()
				continue

			# Unlock the game state
			self.state.unlock()

			self.state.setClientMode(ClientMode.PLANNED)

			# Free up the event loop
			await asyncio.sleep(0.005)

			# wait = True
