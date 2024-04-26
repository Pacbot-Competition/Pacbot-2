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
			# Your client will not need to do this (in fact it won't be able to
   			# since it shouldn't connect as a trusted client).	
			# self.state.queueAction(4, Directions.RIGHT)

			if self.state.gameMode == GameModes.PAUSED:
				print("paused")
			else:
				# Print where the red ghost is
				print(f'{self.state.ghosts[0].location.row} {self.state.ghosts[0].location.col}')

			# Unlock the game state
			self.state.unlock()

			# Free up the event loop
			await asyncio.sleep(0)
