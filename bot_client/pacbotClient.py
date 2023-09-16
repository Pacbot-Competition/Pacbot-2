# JSON (for reading config.json)
import json

# Asyncio (for asynchronous communication)
import asyncio

# Websockets imports
from websockets.sync.client import connect, ClientConnection
from websockets.exceptions import ConnectionClosedError

# Restore the ability to use Ctrl+C within asyncio
import signal
signal.signal(signal.SIGINT, signal.SIG_DFL)

# Font color modifiers
GREEN = '\033[31m'
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
	Asyncio implementation of a websocket client to communicate with the
	Pacbot game server.
	'''

	# Constructor
	def __init__(self, connect_url: str) -> None:
		self.connect_url: str = connect_url
		self._socket_open: bool = False
		self.connection: ClientConnection

	# Connect and run
	async def run(self) -> None:

		# Connect to the websocket server
		await self.connect()

		try: # Try receiving messages indefinitely
			await self.recv_loop()
		finally: # Disconnect once the connection is over
			await self.disconnect()

	# Connect to the websocket server
	async def connect(self) -> None:

		try: # Connect to the specified URL
			self.connection = connect(self.connect_url)
			self._socket_open = True
		except ConnectionRefusedError:
			print(
				f"{GREEN}Websocket connection refused [{self.connect_url}]\n"
				f"Are the address and port correct, and is the "
				f"server running?{NORMAL}"
			)
			return

	# Disconnect from the websocket server
	async def disconnect(self) -> None:

		# Close the connection
		self.connection.close()
		self._socket_open = False

	# Receive loop for capturing messages
	async def recv_loop(self) -> None:

		# Loop indefinitely
		while self._socket_open:
			try: # Receive values as long as the connection is open
				message = self.connection.recv()
				print(message)
			except ConnectionClosedError: # Break once the conneciton is closed
				break

# Main loop
async def main():

	# Get the URL to connect to
	connect_url = get_connect_url()
	pbc = PacbotClient(connect_url)
	await pbc.run()

	# Once the connection is closed, end the event loop
	loop = asyncio.get_event_loop()
	loop.stop()

if __name__ == '__main__':

	# Run the event loop forever
	loop = asyncio.get_event_loop()
	loop.create_task(main())
	loop.run_forever()