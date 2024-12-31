# Asyncio (for concurrency)
import asyncio

import socket

import random

# Game state
from gameState import *

from variables import *

from operator import itemgetter

from websockets.sync.client import ClientConnection

from serverMessage import ServerMessage

from search import bfs
import copy
from grid import grid

import numpy as np

import time

D_MESSAGES: list[bytes] = [b'w', b'a', b's', b'd', b'.']
TICK_ESTIMATE_BY_LEVEL = [12,int(12*1.5),12*2,int(12*2.5),12*3]

FOOD_POSITIONS = []
gs = GameState()
for col in range(28):
    for row in range(31):
        if gs.pelletAt(row,col):
            FOOD_POSITIONS.append((col,row))

class DeepDecisionModule:
    def __init__(self, state: GameState) -> None:
        # Game state object to store the game information
        self.state = state
        #self.sock = socket.socket()
        #self.sock.connect(("192.168.0.100",1337))
        self.depth = 6

    def set_connection(self, connection: ClientConnection):
        self.connection = connection
        print("connected")

    def _get_direction(self, p_loc:Location, next_loc:tuple):
        if p_loc.col == next_loc[0]:
            if p_loc.row < next_loc[1]:
                return Directions.DOWN
            else:
                return Directions.UP
        else:
            if p_loc.col < next_loc[0]:
                return Directions.RIGHT
            else:
                return Directions.LEFT

    def _update_game_state(self):
        message = self.connection.recv()

        # Convert the message to bytes, if necessary
        messageBytes: bytes
        if isinstance(message, bytes):
            messageBytes = message # type: ignore
        else:
            messageBytes = message.encode('ascii') # type: ignore

        # Update the state, given this message from the server
        self.state.update(messageBytes,True)

    def _send_command_message_to_target(self, direction):
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
        self.sock.send(direction)
        print(direction)

    def _send_socket_stop_command(self):
        self.sock.send(b'x')
        print("stay in place")

    # New stuff here
    def evaluationFunction(self,state:GameState):
        #curr_time = time.time_ns()
        """Calculate distance to the nearest food"""
        min_food_distance = 10000 # larger than 31^2 + 28^2, largest squared distance
        for col,row in FOOD_POSITIONS:
            if state.pelletAt(row,col):
                dist = (row - state.pacmanLoc.row)**2 + (col - state.pacmanLoc.col)**2
                if dist < min_food_distance:
                    min_food_distance = dist

        """Calculate the distance to nearest ghost"""
        ghostPositions = [ghost.location for ghost in state.ghosts]
        scaredTimes = [ghost.frightSteps for ghost in state.ghosts]
        if len(ghostPositions) > 0:
            distanceToGhost = [abs(state.pacmanLoc.col - loc.col) + abs(state.pacmanLoc.row - loc.row) for loc in ghostPositions]
            min_ghost_distance = distanceToGhost[np.argmin(distanceToGhost)]
            nearestGhostScaredTime = scaredTimes[np.argmin(distanceToGhost)]
            # avoid certain death
            if min_ghost_distance <= 1 and nearestGhostScaredTime == 0:
                return -999999
            # eat a scared ghost
            if min_ghost_distance <= 1 and nearestGhostScaredTime > 0:
                return 999999

        #print(f"eval:{time.time_ns()-curr_time}")

        return state.currScore * 5 - min_food_distance

    def deepSearch(self, depth, state: GameState):

        #curr_time = time.time_ns()
        
        if state.currLives == 0 or depth == self.depth or state.numPellets() == 0:
            return self.evaluationFunction(state) - depth * 100
        p_loc = state.pacmanLoc
        targets = [(p_loc.col,p_loc.row), (p_loc.col-1, p_loc.row), (p_loc.col+1, p_loc.row), (p_loc.col, p_loc.row+1), (p_loc.col, p_loc.row - 1)]
        directions =  [Directions.NONE, Directions.LEFT, Directions.RIGHT, Directions.DOWN, Directions.UP]
        heuristics = []
        for i in range(len(targets)):
            target_loc = targets[i]
            if state.wallAt(target_loc[1],target_loc[0]):
                continue
            sim_state = GameState()
            sim_state.update(state.serialize(),True)
            #simulated_state = copy.deepcopy(state)
            alive = sim_state.simulateAction(TICK_ESTIMATE_BY_LEVEL[state.currLevel - 1],directions[i])
            heuristics.append(self.deepSearch(depth+1,sim_state))

        if len(heuristics) == 0:
            return self.evaluationFunction(state) - depth * 100

        max_val = max(heuristics)

        #print(f"deep:{time.time_ns()-curr_time}")
        
        return max_val + self.evaluationFunction(state) - depth * 100

    def tick(self):
        if self.state.gameMode == GameModes.PAUSED:
                
            #self.sock.send(b'p')
            print("game paused")
            if self.state.numPellets() == 244:
                self.grid = copy.deepcopy(grid)
                self.num_powerup = 4
            return
        
        if self.state:
            if self.state.pacmanLoc.row == 32 :
                return
            p_loc = self.state.pacmanLoc

            if self.state.numPellets() <= self.depth:
                self.depth = self.state.numPellets() - 1
            targets = [(p_loc.col,p_loc.row), (p_loc.col-1, p_loc.row), (p_loc.col+1, p_loc.row), (p_loc.col, p_loc.row+1), (p_loc.col, p_loc.row - 1)]
            directions =  [Directions.NONE, Directions.LEFT, Directions.RIGHT, Directions.DOWN, Directions.UP]
            action_scores = []

            curr_time = time.time()

            for i in range(len(targets)):
                target_loc = targets[i]
                if self.state.wallAt(target_loc[1],target_loc[0]):
                    continue
                sim_state = GameState()
                sim_state.update(self.state.serialize(),True)
                #simulated_state = copy.deepcopy(self.state)
                alive = sim_state.simulateAction(TICK_ESTIMATE_BY_LEVEL[self.state.currLevel - 1],directions[i])
                action_scores.append(self.deepSearch(0, sim_state))
                

            max_action = max(action_scores)
            max_indices = [index for index in range(len(action_scores)) if action_scores[index] == max_action]
            chosenIndex = max_indices[-1]

            print(time.time() - curr_time)

            next_loc = targets[chosenIndex]
            print(directions[chosenIndex])
            if next_loc != p_loc:
                self._send_command_message_to_target(directions[chosenIndex])
                #self._send_socket_command_to_target(p_loc, next_loc)
                return
        #self._send_socket_stop_command()
        self._send_stop_command()

    async def decisionLoop(self) -> None:
		# Receive values as long as we have access
        resetted = False
        while self.state.isConnected():
            self._update_game_state()
			# If the current messages haven't been sent out yet, skip this iteration
            # if len(self.state.writeServerBuf):
            #     await asyncio.sleep(0)
            #     continue
            if not resetted:
                print("reset")
                #self.sock.send(b'r')
                resetted = True

			# Lock the game state
            #self.state.lock()

			# Write back to the server, as a test (move right)
            self.tick()

			# Unlock the game state
            # self.state.unlock()

			# Print that a decision has been made

			#Free up the event loop
            
            await asyncio.sleep(0.01)
