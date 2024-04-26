# Library for UDP sockets
import socket

# Enums for command info
from enum import IntEnum

# Terminal colors
from terminalColors import *

class CommandType(IntEnum):
    STOP=0
    START=1
    FLUSH=2
    MOVE=3

class CommandDirection(IntEnum):
    NONE=-1
    NORTH=0
    EAST=1
    WEST=2
    SOUTH=3

dirMap = {
    b'w': CommandDirection.NORTH,
    b'a': CommandDirection.WEST,
    b's': CommandDirection.SOUTH,
    b'd': CommandDirection.EAST
}

class RobotSocket:

    def __init__(self, robotIP: str, robotPort: int) -> None:

        # Robot address
        self.robotIP = robotIP
        self.robotPort = robotPort

        # UDP Socket
        self.sock = socket.socket(socket.AF_INET, # Internet
                            socket.SOCK_DGRAM) # UDP
        self.sock.setblocking(False)

        # Received sequence number and data
        self.recvSeq: int
        self.recvData: bytes = bytes([0,0,0,0,0,0,0])

        # Data
        self.NULL: int = 0
        self.seq0: int = 1
        self.seq1: int = 0
        self.typ:  int = int(CommandType.FLUSH)
        self.val1: int = 0
        self.val2: int = 0
        self.done: bool = False

    #     self.doneEventSubscribers=[]

    # def notifyDoneEvent(self, done):
    #     for handler in self.doneEventSubscribers:


    # def registerDoneHandler(self, doneEventHandler):
    #     pass

    # def unRegisterDoneHandler(self, doneEventHandler):
    #     raise Exception("unimplemented")



    def moveNoCoal(self, command: bytes, row: int, col: int, dist: int) -> bool:


        if command != b'w' and command != b'a' and command != b's' and command != b'd':
            return False

        # Update the sequence number, if applicable
        self.updateSeq()

        # Overwrite the output for a move command
        self.typ  = int(CommandType.MOVE)
        self.val1 = dirMap[command]
        self.val2 = dist

        # Dispatch the message
        self.dispatch(row, col)

        print(f'{CYAN}sending command{NORMAL}', command, dist, '->', row, col, " seqno: ", int(self.seq1 << 8 | self.seq0))
        return True

    def flush(self, row: int, col: int) -> None:

        print('flush', row, col)

        # Update the sequence number, if applicable
        self.updateSeq()

        # Overwrite the output for a flush
        self.seq0 = self.recvData[2]
        self.seq1 = self.recvData[1]
        self.typ  = int(CommandType.FLUSH)
        self.val1 = 0
        self.val2 = 0

        # Dispatch the message
        self.dispatch(row, col)

    def start(self) -> None:

        print('start')

        # Update the sequence number, if applicable
        self.updateSeq()

        # Overwrite the output for a flush
        self.seq0 = self.recvData[2]
        self.seq1 = self.recvData[1]
        self.typ  = int(CommandType.START)
        self.val1 = 0
        self.val2 = 0

        # Dispatch the message
        self.dispatch(0, 0)

    def stop(self) -> None:

        print('stop')

        # Update the sequence number, if applicable
        self.updateSeq()

        # Overwrite the output for a flush
        self.seq0 = self.recvData[2]
        self.seq1 = self.recvData[1]
        self.typ  = int(CommandType.STOP)
        self.val1 = 0
        self.val2 = 0

        # Dispatch the message
        self.dispatch(0, 0)

    def wait(self) -> bool:
        try:
            while True:
                self.recvData, _ = self.sock.recvfrom(1024) # type: ignore
        except:
            pass

        # Received sequence number
        self.recvSeq = (self.recvData[1] << 8 | self.recvData[2]) # type: ignore

        # Is done
        self.done = not bool(self.recvData[5])

        return self.done

    def updateSeq(self) -> None:

        # Send the message only if up to date
        if self.recvSeq == (self.seq1 << 8 | self.seq0):

            print(f'{GREEN}ack #{self.recvSeq}{NORMAL}')

            # Increment the sequence number
            self.seq0 += 1

            # First overflow
            if self.seq0 > 127:
                self.seq0 = 0
                self.seq1 += 1

            # Second overflow
            if self.seq1 > 127:
                self.seq1 = 0

    def isPending(self) -> bool:
        return self.recvSeq < (self.seq1 << 8 | self.seq0)

    def dispatch(self, row: int, col: int) -> None:

        message = ""
        inputString = "{{[{:02x}][{:02x}][{:02x}][{:02x}][{:02x}][{:02x}][{:02x}][{:02x}]}}".format(
            self.NULL, self.seq1, self.seq0, self.typ, row, col, self.val1, self.val2
        )
        inputString = inputString + '\n'
        i = 0

        while i < len(inputString):
            currentChar = inputString[i]
            if (currentChar == "["):
                currentChar = chr(int("0x" + inputString[i+1:i+3], 16))
                i += 3
            message += currentChar
            i += 1

        message = bytes(message, "ascii")

        print(message)

        self.sock.sendto(message, (self.robotIP, self.robotPort))