# Asyncio (for concurrency)
import asyncio

import socket

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

class OldDecisionModule:
    def __init__(self, state: GameState) -> None:
        # Game state object to store the game information
        self.state = state
        self.previous_loc = None
        self.direction = Directions.RIGHT
        self.grid = copy.deepcopy(grid)
        #self.sock = socket.socket()
        #self.sock.connect(("192.168.0.100",1337))
        print("connected")
        self.num_powerup = 4
        self.last_life = 3
        self.depth = 6

    def set_connection(self, connection: ClientConnection):
        self.connection = connection

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

    def _target_is_invalid(self, target_loc):
        return self.grid[target_loc[0]][target_loc[1]] in [I, n]

    def _simulate_surrounded_action(self,state:GameState,pacmanDir:Directions) -> tuple[GameState,bool,bool]:
        '''takes in the planned direction of the pacbot and the current state. Simulates the action and returns a tuple containing three results:The updated simulated state, first boolean is whether the action is valid (is the pacbot alive). Second boolean is whether the action brings pacbot out of surrounded state'''
        simulated_state = copy.deepcopy(state)
        alive = simulated_state.simulateAction(TICK_ESTIMATE_BY_LEVEL[self.state.currLevel - 1],pacmanDir) #simulated_state.updatePeriod - simulated_state.currTicks%simulated_state.updatePeriod

        return (simulated_state,alive,surrounded)

    def _state_to_loc(self,state:GameState) -> tuple[int,int]:
        return (state.pacmanLoc.col,30-state.pacmanLoc.row)
    

    def _update_game_state(self):
        p_loc = (self.state.pacmanLoc.col, 30-self.state.pacmanLoc.row)
        if p_loc[0] < 0 or p_loc[1] < 0 or p_loc[0] >= len(self.grid) or p_loc[1] >= len(self.grid[0]):
            return
        if self.grid[p_loc[0]][p_loc[1]] in [o, O]:
            self.grid[p_loc[0]][p_loc[1]] = e

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
        self.sock.send(direction)
        print(direction)

    def _send_socket_stop_command(self):
        self.sock.send(b'x')
        print("stay in place")

    def update_state(self):
        #TODO check if prev_loc has correct x and y order and whether y value need to be re calculated
        if not self.state.pacmanLoc.at(col=self.previous_loc.col,row=self.previous_loc.row):
            if self.previous_loc is not None:
                self.direction = self._get_direction((self.previous_loc.col, 30 - self.previous_loc.row), (self.state.pacmanLoc.col,30 - self.state.pacmanLoc.row))
            self.previous_loc = self.state.pacmanLoc if self.state else None

    # New stuff here
    def evaluationFunction(self,state:GameState):
        """Calculate distance to the nearest food"""
        min_food_distance = 10000 # larger than 31^2 + 28^2, largest squared distance
        for x in range(31):
            for y in range(28):
                if state.fruitAt(y,x):
                    dist = (y - state.pacmanLoc.row)**2 + (x - state.pacmanLoc.col)**2
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

        return state.currScore() * 5 - min_food_distance

    def deepSearch(self, depth, state: GameState):
        p_loc = self._state_to_loc(state)
        targets = [p_loc, (p_loc[0] - 1, p_loc[1]), (p_loc[0] + 1, p_loc[1]), (p_loc[0], p_loc[1] - 1), (p_loc[0], p_loc[1] + 1)]
        directions =  [Directions.NONE, Directions.LEFT, Directions.RIGHT, Directions.DOWN, Directions.UP]
        heuristics = []
        for i in range(len(targets)):
            target_loc = targets[i]
            if self._target_is_invalid(target_loc):
                continue
            simulated_state = copy.deepcopy(state)
            alive = simulated_state.simulateAction(TICK_ESTIMATE_BY_LEVEL[self.state.currLevel - 1],directions[i])
            heuristics.append(self.deepSearch(depth+1,simulated_state))

        if len(heuristics) == 0 or state.currLives == 0 or depth == self.depth or state.numPellets() == 0:
            return self.evaluationFunction(state) - depth * 100

        max_val = max(heuristics)
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
            self._update_game_state()
            p_loc = (self.state.pacmanLoc.col, 30-self.state.pacmanLoc.row)
            next_loc = self._find_best_target(p_loc)
            if next_loc != p_loc:
                self._send_command_message_to_target(p_loc, next_loc)
                #self._send_socket_command_to_target(p_loc, next_loc)
                if self.grid[next_loc[0]][next_loc[1]] == O:
                    self.num_powerup -= 1
                print(self._get_direction(p_loc,next_loc))
                return
        #self._send_socket_stop_command()
        self._send_stop_command()

    async def decisionLoop(self) -> None:
		# Receive values as long as we have access
        resetted = False
        last_time = time.time()
        while self.state.isConnected():
			# If the current messages haven't been sent out yet, skip this iteration
            # if len(self.state.writeServerBuf):
            #     await asyncio.sleep(0)
            #     continue
            if not resetted:
                print("reset")
                #self.sock.send(b'r')
                resetted = True

			# Lock the game state
            self.state.lock()

			# Write back to the server, as a test (move right)
            self.tick()

			# Unlock the game state
            self.state.unlock()

			# Print that a decision has been made

			# Free up the event loop
            # curr_time = time.time()
            # print(curr_time - last_time)
            # last_time = curr_time
            await asyncio.sleep(0.001)
