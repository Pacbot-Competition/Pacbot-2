# Asyncio (for concurrency)
import asyncio

# Game state
from gameState import *

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

			# Write back to the server, as a test (move right)
			self.state.queueAction(4, Directions.RIGHT)

			# Unlock the game state
			self.state.unlock()

			# Print that a decision has been made
			print('decided')

			# Free up the event loop
			await asyncio.sleep(0)
