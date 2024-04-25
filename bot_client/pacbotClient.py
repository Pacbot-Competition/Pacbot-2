# JSON (for reading config.json)
import json

# Asyncio (for concurrency)
import asyncio

# Websockets (for communication with the server)
from websockets.sync.client import connect, ClientConnection # type: ignore
from websockets.exceptions import ConnectionClosedError # type: ignore
from websockets.typing import Data # type: ignore

# Game state
from gameState import GameState, ClientMode

# Decision module
from policies.astar.decisionModule import DecisionModule

# Robot socket
from robotSocket import RobotSocket

# Server messages
from serverMessage import *

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

# Get the simulation flag from the config.json file
def getSimulationFlag() -> bool:

	# Read the configuration file
	with open('../config.json', 'r', encoding='UTF-8') as configFile:
		config = json.load(configFile)

	# Return the websocket connect address
	return config["PythonSimulation"]

# Get the robot address from the config.json file
def getRobotAddress() -> tuple[str, int]:

	# Read the configuration file
	with open('../config.json', 'r', encoding='UTF-8') as configFile:
		config = json.load(configFile)

	# Return the websocket connect address
	return config["RobotIP"], config['RobotPort']

class PacbotClient:
	'''
	Sample implementation of a websocket client to communicate with the
	Pacbot game server, using asyncio.
	'''

	def __init__(self, connectURL: str, simulationFlag: bool, robotAddress: tuple[str, int]) -> None:
		'''
		Construct a new Pacbot client object
		'''

		# Connection URL (starts with ws://)
		self.connectURL: str = connectURL

		# Simulation flag (bool)
		self.simulationFlag: bool = simulationFlag

		# Robot IP and port
		self.robotIP: str = robotAddress[0]
		self.robotPort: int = robotAddress[1]

		# Private variable to store whether the socket is open
		self._socketOpen: bool = False

		# Connection object to communicate with the server
		self.connection: ClientConnection

		# Game state object to store the game information
		self.state: GameState = GameState()

		# Decision module (policy) to make high-level decisions
		self.decisionModule: DecisionModule = DecisionModule(self.state)

		# Robot socket (comms) to dispatch low-level commands
		self.robotSocket: RobotSocket = RobotSocket(self.robotIP, self.robotPort)

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
					self.commsLoop(),
					self.decisionModule.decisionLoop()
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

				# Receive a message from the connection
				message: Data = self.connection.recv()

				# Convert the message to bytes, if necessary
				messageBytes: bytes
				if isinstance(message, bytes):
					messageBytes = message # type: ignore
				else:
					messageBytes = message.encode('ascii') # type: ignore

				# Update the state, given this message from the server
				self.state.update(messageBytes)

				# Write a response back to the server if necessary
				if (self.simulationFlag):
					if self.state.writeServerBuf and self.state.writeServerBuf[0].tick():
						response: bytes = self.state.writeServerBuf.popleft().getBytes()
						self.connection.send(response)

				# Free the event loop to allow another decision
				await asyncio.sleep(0)

			# Break once the connection is closed
			except ConnectionClosedError:
				print('Connection lost...')
				self.state.setConnectionStatus(False)
				break

	async def commsLoop(self) -> None:
		'''
		Communication loop for sending messages to the robot
		'''

		# Quit if in simulation
		if (self.simulationFlag):
			print("Simulation Mode: No Robot")
			return

		# Keep track if the first iteration has taken place
		firstIt = True

		# Keep track of whether we've sent a message before the robot told us it's done
		# sent = False

		# Keep sending messages as long as the server connection is open
		while self.isOpen():

			# Try to receive messages (and skip to except in case of an error)
			try:

				# Wait until the bot stops sending messages, and check if it's done
				robotIsDone = self.robotSocket.wait()

				if (self.state.isSent() and robotIsDone):
					print('robot just told us it\'s done')
					self.state.setClientMode(ClientMode.DONE)

				# Not ready to send a new message yet
				if (self.state.clientMode != ClientMode.PLANNED):
					await asyncio.sleep(0)
					continue

				# Handle first iteration (flush)
				if firstIt:
					self.robotSocket.start()
					while (self.state.isLocked() or not self.state.isDone()):
						await asyncio.sleep(0)
					self.state.lock()
					self.robotSocket.flush(self.state.pacmanLoc.row, self.state.pacmanLoc.col)
					self.state.unlock()
					self.state.setClientMode(ClientMode.SENT)
					firstIt = False

				# Otherwise, send out relevant messages
				else:
					if self.state.writeServerBuf and self.state.writeServerBuf[0].tick():
						srvmsg: ServerMessage = self.state.writeServerBuf.popleft()
						msg = srvmsg.getBytes()
						dist, row, col = srvmsg.dist, srvmsg.row, srvmsg.col
						self.robotSocket.moveNoCoal(msg, row, col, dist)
						if self.state.writeServerBuf:
							self.state.writeServerBuf[0].skipDelay()
						self.state.setClientMode(ClientMode.SENT)

				# Free the event loop to allow another decision
				await asyncio.sleep(0.025)

			# Break once the connection is closed
			except ConnectionClosedError:
				print('Comms lost...')
				self.state.setConnectionStatus(False)
				break

# Main function
async def main():

	# Get the URL to connect to
	connectURL = getConnectURL()
	simulationFlag = getSimulationFlag()
	robotAddress = getRobotAddress()
	client = PacbotClient(connectURL, simulationFlag, robotAddress)
	await client.run()

	# Once the connection is closed, end the event loop
	loop = asyncio.get_event_loop()
	loop.stop()

if __name__ == '__main__':

	# Run the event loop forever
	loop = asyncio.new_event_loop()
	loop.create_task(main())
	loop.run_forever()