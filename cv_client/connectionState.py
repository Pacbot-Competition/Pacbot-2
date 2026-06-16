# Asyncio (for the outbound message queue)
import asyncio

class ConnectionState:
	'''
	Connection state object for the CV client, keeping track of when
	it is connected to the server
	'''

	def __init__(self) -> None:
		'''
		Construct a new game state object
		'''

		# Bounded queue of messages to write back to the server. The send loop
		# awaits this queue, so there is no polling; when the queue is full the
		# oldest position is dropped in favor of the newest.
		self.writeServerBuf: asyncio.Queue = asyncio.Queue(maxsize=64)

		# Whether the client is currently connected to the server
		self._connected: bool = False

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

		self._enqueue(bytes([ord('x'), row, col]))

	def signalClose(self) -> None:
		'''
		Queue a None sentinel to wake a blocked send loop so it can observe the
		closed connection and exit.
		'''

		self._enqueue(None)

	def _enqueue(self, item: 'bytes | None') -> None:
		'''
		Put an item on the outbound queue without blocking. If the queue is
		full (the send loop has fallen behind), drop the oldest message so the
		newest position always gets through.
		'''

		if self.writeServerBuf.full():
			try:
				self.writeServerBuf.get_nowait()
			except asyncio.QueueEmpty:
				pass
		try:
			self.writeServerBuf.put_nowait(item)
		except asyncio.QueueFull:
			pass
