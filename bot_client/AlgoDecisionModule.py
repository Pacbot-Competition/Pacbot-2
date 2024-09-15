# Asyncio (for concurrency)
import asyncio

import socket

# Game state
from gameState import *

from variables import *

from operator import itemgetter

from websockets.sync.client import ClientConnection

from serverMessage import ServerMessage

D_MESSAGES: list[bytes] = [b'w', b'a', b's', b'd', b'.']

PACBOT_IP = "192.168.0.100"
PACBOT_PORT = 1234

ALGO_IP = "127.0.0.1"
ALGO_PORT = 3000

class AlgoDecisionModule:
    def __init__(self, state: GameState) -> None:
        # Game state object to store the game information
        self.state = state
        self.previous_loc = None
        self.direction = Directions.RIGHT
        self.pacbot_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        #self.pacbot_sock.connect((PACBOT_IP,PACBOT_PORT))
        self.algo_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.algo_sock.connect((ALGO_IP,ALGO_PORT))

    def set_connection(self, connection: ClientConnection):
        self.connection = connection
        self._update_game_state()

    def _direction_to_str(self,dir:Directions):
        result=""
        match dir:
            case Directions.UP:
                result+="u"
            case Directions.RIGHT:
                result+="r"
            case Directions.LEFT:
                result+="l"
            case Directions.DOWN:
                result+="d"
            case Directions.NONE:
                result+="n"
        return result

    def _serialize_state(self):
        result = ""
        result+= self._direction_to_str(self.direction)
        match self.state.gameMode:
            case GameModes.PAUSED:
                result+="p"
            case GameModes.SCATTER:
                result+="s"
            case GameModes.CHASE:
                result+="c"
        result+=f"{self.state.pacmanLoc.col},{self.state.pacmanLoc.row},"
        result+=f"{self.state.fruitLoc.col},{self.state.fruitLoc.row}"
        for ghost in self.state.ghosts:
            match ghost.color:
                case GhostColors.RED:
                    result+="r"
                case GhostColors.PINK:
                    result+="p"
                case GhostColors.CYAN:
                    result+="c"
                case GhostColors.ORANGE:
                    result+="o"
            result+=f"{ghost.location.col},{ghost.location.row},"
            result+=f"{ghost.frightSteps}"
            result+=self._direction_to_str(ghost.plannedDirection)
        result+=":"
        pelletGrid = "".join(str(x) for x in self.state.pelletArr)
        result+=pelletGrid
            
        return result

    def _get_direction(self, p_loc, next_loc):
        if p_loc[0] == next_loc[0]:
            if p_loc[1] < next_loc[1]:
                return Directions.UP
            else:
                return Directions.DOWN
        else:
            if p_loc[0] < next_loc[0]:
                return Directions.RIGHT
            else:
                return Directions.LEFT

    def _update_game_state(self):
        p_loc = (self.state.pacmanLoc.col, 30-self.state.pacmanLoc.row)
        pass

    def _send_command_message_to_target(self, p_loc, target):
        direction = self._get_direction(p_loc, target)
        # self.state.queueAction(4,direction)
        self.connection.send(ServerMessage(D_MESSAGES[direction], 4).getBytes())

    def _send_stop_command(self):
        #self.state.queueAction(4,Directions.NONE)
        self.connection.send(ServerMessage(D_MESSAGES[4], 4).getBytes())

    def _send_socket_command_to_target(self, p_loc, target):
        direction = self._get_direction(p_loc, target)
        match direction:
            case Directions.UP:
                direction = b'n'
            case Directions.DOWN:
                direction = b's'
            case Directions.LEFT:
                direction = b'w'
            case Directions.RIGHT:
                direction = b'e'
        self.pacbot_sock.send(direction)

    def _send_socket_stop_command(self):
        self.pacbot_sock.send(b'x')

    def _listen_for_algo_command(self):
        self.algo_sock.send(self._serialize_state().encode())
        msg = self.algo_sock.recv(512)
        print(msg)

    def update_state(self):
        #TODO check if prev_loc has correct x and y order and whether y value need to be re calculated
        if not self.state.pacmanLoc.at(col=self.previous_loc.col,row=self.previous_loc.row):
            if self.previous_loc is not None:
                self.direction = self._get_direction((self.previous_loc.col, 30 - self.previous_loc.row), (self.state.pacmanLoc.col,30 - self.state.pacmanLoc.row))
            self.previous_loc = self.state.pacmanLoc if self.state else None

    def tick(self):
        if self.state.gameMode == GameModes.PAUSED:
            #self._send_socket_stop_command()
            self._send_stop_command()
            return
        if self.state:
            self._update_game_state()
            p_loc = (self.state.pacmanLoc.col, 30-self.state.pacmanLoc.row)
            #self._send_command_message_to_target(p_loc, next_loc)
            #self._send_socket_command_to_target(p_loc, next_loc)
            self._listen_for_algo_command()
            return
        self._send_socket_stop_command()
        #self._send_stop_command()

    async def decisionLoop(self) -> None:
		# Receive values as long as we have access
        while self.state.isConnected():
			# If the current messages haven't been sent out yet, skip this iteration
            # if len(self.state.writeServerBuf):
            #     await asyncio.sleep(0)
            #     continue

			# Lock the game state
            self.state.lock()

			# Write back to the server, as a test (move right)
            self.tick()

			# Unlock the game state
            self.state.unlock()

			# Print that a decision has been made
            print('decided')

			# Free up the event loop
            await asyncio.sleep(0.1)
