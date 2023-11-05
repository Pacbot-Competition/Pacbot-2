# Asyncio (for concurrency)
import asyncio

# Game state
from gameState import *

# A-Star Policy
from aStarPolicy import *

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
		self.policy = AStarPolicy(state, newLocation(5, 21))

	async def decisionLoop(self) -> None:
		'''
		Decision loop for Pacbot
		'''

		# Receive values as long as we have access
		while self.state.isConnected():

			'''
			WARNING: 'await' statements should be routinely placed
			to free the event loop to receive messages, or the
			client may fall behind on updating the game state!
			'''

			# If the current messages haven't been sent out yet, skip this iteration
			if len(self.state.writeServerBuf):
				await asyncio.sleep(0)
				continue

			# Lock the game state
			self.state.lock()

			# Figure out which actions to take, according to the policy
			await self.policy.act()

			# Unlock the game state
			self.state.unlock()

			# Free up the event loop (a good chance to talk to the bot!)
			await asyncio.sleep(0)
