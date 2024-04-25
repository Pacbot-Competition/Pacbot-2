class ServerMessage:
  '''
	Sample implementation of a message object for communication with
  the Pacbot server.
	'''

  def __init__(self, messageBytes: bytes, numTicks: int, dist: int, row: int, col: int):
    '''
		Construct a new server object
		'''
    self.messageBytes = messageBytes
    self.waitTicks = numTicks
    self.dist = dist
    self.row = row
    self.col = col

  def tick(self):
    '''
		Each tick, decrement the waiting ticks, and return whether the message
    is ready to send to the server.
		'''
    self.waitTicks -= 1
    return (self.waitTicks <= 0)
  
  def skipDelay(self):
    '''
    Skip the delay for a message sent to the robot
    '''
    self.waitTicks = 0

  def getBytes(self):
    '''
		Return the bytes of this server message.
		'''
    return self.messageBytes
