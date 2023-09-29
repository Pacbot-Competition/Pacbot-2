# JSON (for reading config.json)
import json

# Asyncio (for concurrency)
import asyncio

# Websockets (for communication with the server)
from websockets.sync.client import connect, ClientConnection # type: ignore
from websockets.exceptions import ConnectionClosedError # type: ignore
from websockets.typing import Data # type: ignore

# Game state
from gameState import GameState

# Decision module
from decisionModule import DecisionModule

# Restore the ability to use Ctrl + C within asyncio
import signal
signal.signal(signal.SIGINT, signal.SIG_DFL)

# Font color modifiers
RED = '\033[31m'
NORMAL = '\033[0m'

# Get the connect URL from the config.json file
def get_connect_url() -> str:

	# Read the configuration file
	with open('../config.json', 'r', encoding='UTF-8') as config_json:
		config_dict = json.load(config_json)

	# Return the websocket connect address
	return f'ws://{config_dict["ServerIP"]}:{config_dict["WebSocketPort"]}'

class PacbotClient:
	'''
	Sample implementation of a websocket client to communicate with the
	Pacbot game server, using asyncio.
	'''

	def __init__(self, connect_url: str) -> None:
		'''
		Construct a new Pacbot client object
		'''

		# Connection URL (starts with ws://)
		self.connect_url: str = connect_url

		# Private variable to store whether the socket is open
		self._socket_open: bool = False

		# Connection object to communicate with the server
		self.connection: ClientConnection

		# Game state object to store the game information
		self.state: GameState = GameState()

		# Decision module (policy) to make high-level decisions
		self.policy: DecisionModule = DecisionModule(self.state)

	async def run(self) -> None:
		'''
		Connect to the server, then run
		'''

		# Connect to the websocket server
		await self.connect()

		try: # Try receiving messages indefinitely
			await asyncio.gather(
				self.recv_loop(),
				self.policy.decision_loop()
			)
		finally: # Disconnect once the connection is over
			await self.disconnect()

	async def connect(self) -> None:
		'''
		Connect to the websocket server
		'''

		# Connect to the specified URL
		try:
			self.connection = connect(self.connect_url)
			self._socket_open = True

		# If the connection is refused, log and return
		except ConnectionRefusedError:
			print(
				f'{RED}Websocket connection refused [{self.connect_url}]\n'
				f'Are the address and port correct, and is the '
				f'server running?{NORMAL}'
			)
			return

	async def disconnect(self) -> None:
		'''
		Disconnect from the websocket server
		'''

		# Close the connection
		if self._socket_open:
			self.connection.close()
		self._socket_open = False

	# Return whether the connection is open
	def is_open(self) -> bool:
		'''
		Check whether the connection is open (unused)
		'''
		return self._socket_open

	async def recv_loop(self) -> None:
		'''
		Receive loop for capturing messages from the server
		'''

		# Receive values as long as the connection is open
		while self._socket_open:

			# Try to receive messages (and skip to except in case of an error)
			try:

				# Receive a message from the connection
				message: Data = self.connection.recv()

				# Convert the message to bytes, if necessary
				message_bytes: bytes
				if isinstance(message, bytes):
					message_bytes = message # type: ignore
				else:
					message_bytes = message.encode('ascii') # type: ignore

				# Update the state, given this message from the server
				self.state.update(message_bytes)

				# Free the event loop to allow another decision
				await asyncio.sleep(0)

			# Break once the connection is closed
			except ConnectionClosedError:
				break

# Main function
async def main():

	# Get the URL to connect to
	connect_url = get_connect_url()
	client = PacbotClient(connect_url)
	await client.run()

	# Once the connection is closed, end the event loop
	loop = asyncio.get_event_loop()
	loop.stop()

if __name__ == '__main__':

	# Run the event loop forever
	loop = asyncio.get_event_loop()
	loop.create_task(main())
	loop.run_forever()