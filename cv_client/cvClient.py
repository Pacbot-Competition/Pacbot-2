# sys (for argv)
import sys

# JSON (for reading config.json)
import json

# Asyncio (for concurrency)
import asyncio

# Typing
from typing import Any

# Websockets (for communication with the server)
from websockets.sync.client import connect, ClientConnection # type: ignore
from websockets.exceptions import ConnectionClosedError # type: ignore
from websockets.typing import Data # type: ignore

# Decision module
from cameraModule import CameraModule

# Import connection state object
from connectionState import ConnectionState

# Restore the ability to use Ctrl + C within asyncio
import signal
signal.signal(signal.SIGINT, signal.SIG_DFL)

# Terminal colors for formatting output text
from terminalColors import *

# Get the connect URL from the config.json file
def getConnectURL() -> str:

	# Read the configuration file
	with open('../config.json', 'r', encoding='UTF-8') as configFile:
		config = json.load(configFile)

	# Return the websocket connect address
	return f'ws://{config["ServerIP"]}:{config["WebSocketPort"]}'

class CvClient:
	'''
	Implementation of a websocket client to communicate with the
	Pacbot game server, using asyncio.
	'''

	def __init__(self, connectURL: str, name: Any) -> None:
		'''
		Construct a new Pacbot client object
		'''

		# Connection URL (starts with ws://)
		self.connectURL: str = connectURL

		# Private variable to store whether the socket is open
		self._socketOpen: bool = False

		# Connection object to communicate with the server
		self.connection: ClientConnection

		# Game state object to store the game information
		self.state: ConnectionState = ConnectionState()

		# Decision module (policy) to make high-level decisions
		self.cameraModule: CameraModule = CameraModule(self.state, name)

	async def run(self) -> None:
		'''
		Connect to the server, then run
		'''

		# Connect to the websocket server
		await self.connect()

		try: # Try receiving messages indefinitely
			if self._socketOpen:
				await asyncio.gather(
					self.receiveLoop(),
					self.cameraModule.decisionLoop()
				)
		finally: # Disconnect once the connection is over
			await self.disconnect()

	async def connect(self) -> None:
		'''
		Connect to the websocket server
		'''

		# Connect to the specified URL
		try:
			self.connection = connect(self.connectURL)
			self._socketOpen = True
			self.state.setConnectionStatus(True)

		# If the connection is refused, log and return
		except ConnectionRefusedError:
			print(
				f'{RED}Websocket connection refused [{self.connectURL}]\n'
				f'Are the address and port correct, and is the '
				f'server running?{NORMAL}'
			)
			return

	async def disconnect(self) -> None:
		'''
		Disconnect from the websocket server
		'''

		# Close the connection
		if self._socketOpen:
			self.connection.close()
		self._socketOpen = False
		self.state.setConnectionStatus(False)

	# Return whether the connection is open
	def isOpen(self) -> bool:
		'''
		Check whether the connection is open (unused)
		'''
		return self._socketOpen

	async def receiveLoop(self) -> None:
		'''
		Receive loop for capturing messages from the server
		'''

		# Receive values as long as the connection is open
		while self.isOpen():

			# Try to receive messages (and skip to except in case of an error)
			try:

				# Receive a message from the connection (unused)
				_: Data = self.connection.recv()

				# Write a response back to the server if necessary
				if self.state.writeServerBuf:
					response: bytes = self.state.writeServerBuf.popleft()
					self.connection.send(response)

				# Free the event loop to allow another decision
				await asyncio.sleep(0)

			# Break once the connection is closed
			except ConnectionClosedError:
				print('Connection lost...')
				self.state.setConnectionStatus(False)
				break

# Main function
async def main():

	# Get the URL to connect to
	connectURL = getConnectURL()
	cameraID = int(sys.argv[1])
	client = CvClient(connectURL, cameraID)
	await client.run()

	# Once the connection is closed, end the event loop
	loop = asyncio.get_event_loop()
	loop.stop()

if __name__ == '__main__':

	# Run the event loop forever
	loop = asyncio.new_event_loop()
	loop.create_task(main())
	loop.run_forever()