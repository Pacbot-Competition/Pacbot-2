# Buffer to collect messages to write to the server
from collections import deque

class ConnectionState:
	'''
	Connection state object for the CV client, keeping track of when
	it is connected to the server
	'''

	def __init__(self) -> None:
		'''
		Construct a new game state object
		'''

		# Buffer of messages to write back to the server
		self.writeServerBuf: deque[bytes] = deque[bytes](maxlen=64)

	def setConnectionStatus(self, connected: bool) -> None:
		'''
		Set the connection status of this game state's client
		'''

		# Update the internal 'connected' state variable
		self._connected = connected

	def isConnected(self) -> bool:
		'''
		Check if the client attached to the game state is connected
		'''

		# Return the internal 'connected' state variable
		return self._connected
	
	def send(self, row: int, col: int) -> None:
		'''
		Helper function to queue a message to be sent to the server, with a
		given Pacbot location, represented as a row and column.
		'''

		self.writeServerBuf.append(
			bytes([ord('x'), row, col])
		)