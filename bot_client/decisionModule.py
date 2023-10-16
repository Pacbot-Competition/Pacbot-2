# Asyncio (for concurrency)
import asyncio

# Game state
from gameState import GameState

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

			# WARNING: 'await' statements should be routinely placed
			# to free the event loop to receive messages, or the
			# client may fall behind on updating the game state!

			# Lock the game state
			self.state.lock()

			# Replace this with the actual decisions for Pacbot
			await asyncio.sleep(0.1)

			self.state.update(self.state.serialize(), lockOverride=True)
			self.state.display()

			# Unlock the game state
			# self.state.unlock()

			# Writing back to the server, as a test (move right)
			# self.state.writeServerBuf.append(b'd')

			# Free up the event loop (a good chance to talk to the bot!)
			await asyncio.sleep(1000)

			# (REMOVE THIS) Unlock the game state
			self.state.unlock()
