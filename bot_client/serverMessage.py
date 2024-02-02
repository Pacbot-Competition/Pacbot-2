class ServerMessage:
  '''
	Sample implementation of a message object for communication with
  the Pacbot server.
	'''

  def __init__(self, messageBytes: bytes, numTicks: int):
    '''
		Construct a new server object
		'''
    self.messageBytes = messageBytes
    self.waitTicks = numTicks

  def tick(self, staleTick=False):
    '''
		Each tick, decrement the waiting ticks, and return whether the message
    is ready to send to the server.
		'''
    if not staleTick:
      self.waitTicks -= 1
    return (self.waitTicks <= 0)

  def getBytes(self):
    '''
		Return the bytes of this server message.
		'''
    return self.messageBytes
