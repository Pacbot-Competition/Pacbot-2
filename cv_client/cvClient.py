# JSON (for reading config.json)
import json

# Asyncio (for concurrency)
import asyncio

# Websockets (for communication with the server)
from websockets.sync.client import connect, ClientConnection # type: ignore
from websockets.exceptions import ConnectionClosed # type: ignore
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

	def __init__(self, connectURL: str) -> None:
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
		self.cameraModule: CameraModule = CameraModule(self.state)

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
					self.sendLoop(),
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

		# Release the camera and its reader thread
		self.cameraModule.release()

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

				# Receive a message from the connection (unused). The blocking
				# recv() runs in a worker thread so it does not stall the event
				# loop (and starve the send and decision loops)
				_: Data = await asyncio.to_thread(self.connection.recv)

				# Free the event loop to allow another decision
				await asyncio.sleep(0)

			# Break once the connection is closed (clean or abnormal closures
			# both subclass ConnectionClosed)
			except ConnectionClosed:
				print('Connection lost...')
				self.state.setConnectionStatus(False)
				self._socketOpen = False

				# Wake the send loop so it can observe the closed connection
				# and exit, instead of blocking on the queue forever
				self.state.signalClose()
				break

	async def sendLoop(self) -> None:
		'''
		Send loop for flushing queued messages to the server as they are
		produced, independent of inbound traffic
		'''

		# Send values as long as the connection is open
		while self.isOpen():

			# Try to send queued messages (and skip to except on error)
			try:

				# Block until the next message is queued. A None sentinel (from
				# signalClose) or a closed socket means it is time to exit.
				response = await self.state.writeServerBuf.get()
				if response is None or not self.isOpen():
					break

				# The blocking send() runs in a worker thread so it does not
				# stall the event loop
				await asyncio.to_thread(self.connection.send, response)

			# Break once the connection is closed (clean or abnormal closures
			# both subclass ConnectionClosed)
			except ConnectionClosed:
				print('Connection lost...')
				self.state.setConnectionStatus(False)
				self._socketOpen = False
				break

# Main function
async def main():

	# Get the URL to connect to
	connectURL = getConnectURL()
	client = CvClient(connectURL)
	await client.run()

	# Once the connection is closed, end the event loop
	loop = asyncio.get_event_loop()
	loop.stop()

if __name__ == '__main__':

	# Run the event loop forever
	loop = asyncio.new_event_loop()
	loop.create_task(main())
	loop.run_forever()