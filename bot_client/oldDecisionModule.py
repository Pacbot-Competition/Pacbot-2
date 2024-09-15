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

PELLET_WEIGHT = 0.65
GHOST_WEIGHT = 0.35
FRIGHTENED_GHOST_WEIGHT = 2 * GHOST_WEIGHT
GHOST_CUTOFF = 10
GHOST_FRIGHTENED_CUTOFF = 13
CHERRY_WEIGHT = 2 * PELLET_WEIGHT
REVERSE_WEIGHT = 15
SURROUNDED_WEIGHT = 1.5

SURROUND_RADIUS = 15
SURROUND_DIFF = 5

TICK_ESTIMATE_BY_LEVEL = [12,int(12*1.5),12*2,int(12*2.5),12*3]

D_MESSAGES: list[bytes] = [b'w', b'a', b's', b'd', b'.']

class OldDecisionModule:
    def __init__(self, state: GameState) -> None:
        # Game state object to store the game information
        self.state = state
        self.previous_loc = None
        self.direction = Directions.RIGHT
        self.grid = copy.deepcopy(grid)
        self.escape_queue = []
        #self.sock = socket.socket()
        #self.sock.connect(("192.168.0.100",1337))
        print("connected")
        self.num_powerup = 4
        self.last_life = 3

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

    def  _find_paths_to_closest_ghosts(self, pac_loc:tuple[int,int]):
        ghosts = self.state.ghosts
        state_paths = [(ghost.frightSteps > 2, bfs(self.grid, pac_loc, (ghost.location.col, 30 - ghost.location.row), GHOST_CUTOFF)) for ghost in ghosts]
        return [sp for sp in state_paths if sp[1] is not None]

    def  _find_paths_to_closest_ghosts_state(self, state:GameState):
        ghosts = state.ghosts
        state_paths = [(ghost.frightSteps > 2, bfs(self.grid, (state.pacmanLoc.col,30-state.pacmanLoc.row), (ghost.location.col, 30 - ghost.location.row), GHOST_CUTOFF)) for ghost in ghosts]
        return [sp for sp in state_paths if sp[1] is not None]

    def _find_distance_of_closest_pellet(self, target_loc):
        return len(bfs(self.grid, target_loc, [o])) - 1

    def _find_distance_of_closest_powerup(self, target_loc):
        return len(bfs(self.grid, target_loc, [O])) - 1

    def _find_distance_to_cherry(self, target_loc):
        return len(bfs(self.grid, target_loc, (self.state.fruitLoc.col, 30- self.state.fruitLoc.row))) - 1

    def _target_is_invalid(self, target_loc):
        return self.grid[target_loc[0]][target_loc[1]] in [I, n]

    def _is_power_pellet_closer(self, path):
        for coord in path:
            if self.grid[coord[0]][coord[1]] == O:
                return True
        return False

    def _get_num_turns(self, p_dir, n_dir):
        lat = [Directions.LEFT, Directions.RIGHT]
        lng = [Directions.DOWN, Directions.UP]

        if p_dir == n_dir:
            return 0
        elif (p_dir in lat and n_dir in lat) or (p_dir in lng and n_dir in lng):
            return 2
        else:
            return 1

    def _get_target_with_min_turning_direction(self, mins):
        turns = [(self._get_num_turns(self.direction, direct), targ) for direct, targ in mins]
        return min(turns, key=itemgetter(0))[1]

    def _is_surrounded(self,p_loc:tuple[int,int]):
        paths_to_ghosts = self._find_paths_to_closest_ghosts(p_loc)
        dists_to_ghosts = [len(x[1]) for x in paths_to_ghosts if not x[0]]
        for i in range(len(dists_to_ghosts)):
            for j in range(i + 1,len(dists_to_ghosts)):
                if dists_to_ghosts[i] < SURROUND_RADIUS and dists_to_ghosts[j] < SURROUND_RADIUS and abs(dists_to_ghosts[i] - dists_to_ghosts[j]) < SURROUND_DIFF:
                    return True
        return False
    
    def _is_surrounded_state(self,state:GameState, surround_radius=SURROUND_RADIUS):
        paths_to_ghosts = self._find_paths_to_closest_ghosts_state(state)
        dists_to_ghosts = [len(x[1]) for x in paths_to_ghosts if not x[0]]
        count = 0
        for i in range(len(dists_to_ghosts)):
            if dists_to_ghosts[i] < surround_radius:
                count += 1
            # for j in range(i + 1,len(dists_to_ghosts)):
            #     if dists_to_ghosts[i] < SURROUND_RADIUS and dists_to_ghosts[j] < SURROUND_RADIUS and abs(dists_to_ghosts[i] - dists_to_ghosts[j]) < SURROUND_DIFF:
            #         return True
        return count >= 2
        return False

    def _simulate_surrounded_action(self,state:GameState,pacmanDir:Directions) -> tuple[GameState,bool,bool]:
        '''takes in the planned direction of the pacbot and the current state. Simulates the action and returns a tuple containing three results:The updated simulated state, first boolean is whether the action is valid (is the pacbot alive). Second boolean is whether the action brings pacbot out of surrounded state'''
        simulated_state = copy.deepcopy(state)
        alive = simulated_state.simulateAction(TICK_ESTIMATE_BY_LEVEL[self.state.currLevel - 1],pacmanDir) #simulated_state.updatePeriod - simulated_state.currTicks%simulated_state.updatePeriod
        surrounded = self._is_surrounded_state(simulated_state, surround_radius=int(1.25*SURROUND_RADIUS))

        return (simulated_state,alive,surrounded)

    def _state_to_loc(self,state:GameState) -> tuple[int,int]:
        return (state.pacmanLoc.col,30-state.pacmanLoc.row)
    
    def _find_escape_route(self,max_dist=int(SURROUND_RADIUS*1.5)) -> list[tuple[int,int]]:
        visited = []
        queue = [(self.state, [])]

        while len(queue) > 0:
            nxt = queue.pop(0)
            curr_state = nxt[0] # self.state - all info
            visited.append(self._state_to_loc(curr_state))
            new_path = copy.deepcopy(nxt[1])
            new_path.append(self._state_to_loc(curr_state))

            long_path = ()
            max_ghost_dist = 0

            p_loc = self._state_to_loc(curr_state)

            targets = [p_loc, (p_loc[0] - 1, p_loc[1]), (p_loc[0] + 1, p_loc[1]), (p_loc[0], p_loc[1] - 1), (p_loc[0], p_loc[1] + 1)]
            directions =  [Directions.NONE, Directions.LEFT, Directions.RIGHT, Directions.DOWN, Directions.UP]

            valid_directions = []
            for i in range(5):
                if not self._target_is_invalid(targets[i]):
                    valid_directions.append(directions[i])
            
            for dir in valid_directions:
                new_state_result = self._simulate_surrounded_action(curr_state,dir) # simulated state, alive, surrounded
                new_pos = self._state_to_loc(new_state_result[0])
                if (self.grid[new_pos[0]][new_pos[1]] == 'O') and new_state_result[1]:
                    new_path.append(self._state_to_loc(new_state_result[0]))
                    print("Condition 1")
                    return new_path
                if (new_state_result[1] and not new_state_result[2]):
                    new_path.append(self._state_to_loc(new_state_result[0]))
                    print("Condition 2")
                    return new_path
                elif new_state_result[1] and self._state_to_loc(new_state_result[0]) not in visited and len(new_path) <= max_dist:
                    queue.append((new_state_result[0],new_path))

                if(len(new_path) >= max_dist):

                    ghost_paths = self._find_paths_to_closest_ghosts_state(curr_state)

                    ghost_dists = [len(x[1]) for x in ghost_paths if not x[0]]

                    median = np.median(ghost_dists)

                    if median > max_ghost_dist:

                        max_ghost_dist = median

                        long_path = new_path

        print("Condition 3")
        return long_path


    def _find_best_target(self, p_loc):
        targets = [p_loc, (p_loc[0] - 1, p_loc[1]), (p_loc[0] + 1, p_loc[1]), (p_loc[0], p_loc[1] - 1), (p_loc[0], p_loc[1] + 1)]
        directions =  [Directions.NONE, Directions.LEFT, Directions.RIGHT, Directions.DOWN, Directions.UP]
        heuristics = []
        for i in range(len(targets)):
            target_loc = targets[i]
            if self._target_is_invalid(target_loc):
                heuristics.append(float('inf'))
                continue
            if self.state.numPellets() - self.num_powerup == 0:
                dist_to_pellet = self._find_distance_of_closest_powerup(target_loc)
            else:
                dist_to_pellet = self._find_distance_of_closest_pellet(target_loc)
            
            paths_to_ghosts = self._find_paths_to_closest_ghosts(target_loc)

            closest_ghost = (None, float('inf'))
            ghosts = []
            for state, path in paths_to_ghosts:
                dist = len(path) - 1
                closest_ghost = (state, dist) if dist < closest_ghost[1] else closest_ghost
                ghosts.append((state, dist))
                if self._is_power_pellet_closer(path):
                    if target_loc == p_loc:
                        return path[1]
                    else:
                        return path[0]

            ghost_heuristic = 0
            for state, dist in ghosts:
                #'''
                if dist < GHOST_CUTOFF:
                    if state == False:
                        ghost_heuristic += pow((GHOST_CUTOFF - dist), 2) * GHOST_WEIGHT
                    else:
                        ghost_heuristic += pow((GHOST_CUTOFF - dist), 2) * -1 * FRIGHTENED_GHOST_WEIGHT
                #'''
                '''
                if state and dist < GHOST_FRIGHTENED_CUTOFF: 
                    ghost_heuristic += pow((GHOST_CUTOFF - closest_ghost[1]), 2) * -1 * FRIGHTENED_GHOST_WEIGHT
                elif not state and dist < GHOST_CUTOFF:
                    ghost_heuristic += pow((GHOST_CUTOFF - closest_ghost[1]), 2) * GHOST_WEIGHT
                '''


            pellet_heuristic = dist_to_pellet * PELLET_WEIGHT
            total_heuristic = pellet_heuristic + ghost_heuristic
 
            if self.state.fruitLoc.col != 32:
                total_heuristic += self._find_distance_to_cherry(target_loc) * CHERRY_WEIGHT
            heuristics.append(total_heuristic)
        # print(heuristics)
        mins = []
        min_heur = float('inf')
        for i, heur in enumerate(heuristics):

            if directions[i] == Directions.UP and self.state.pacmanLoc.getDirection() == Directions.DOWN:
                heur += REVERSE_WEIGHT
            elif directions[i] == Directions.DOWN and self.state.pacmanLoc.getDirection() == Directions.UP:
                heur += REVERSE_WEIGHT
            elif directions[i] == Directions.LEFT and self.state.pacmanLoc.getDirection() == Directions.RIGHT:
                heur += REVERSE_WEIGHT
            elif directions[i] == Directions.RIGHT and self.state.pacmanLoc.getDirection() == Directions.LEFT:
                heur += REVERSE_WEIGHT
            if heur < min_heur:
                min_heur = heur
                mins = [(directions[i], targets[i])]
            elif heur == min_heur:
                mins.append((directions[i], targets[i]))
        return self._get_target_with_min_turning_direction(mins)

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

    def tick(self):
        if self.state.gameMode == GameModes.PAUSED:
                
            #self.sock.send(b'p')
            print("game paused")
            self.escape_queue = []
            if self.state.numPellets() == 244:
                self.grid = copy.deepcopy(grid)
                self.num_powerup = 4
            return
        
        if self.state:
            if self.state.pacmanLoc.row == 32 :
                return
            self._update_game_state()
            p_loc = (self.state.pacmanLoc.col, 30-self.state.pacmanLoc.row)
            # surrounded = self._is_surrounded(p_loc)
            # if len(self.escape_queue) == 0 and surrounded and self.state.currTicks > 6:
            #     print("generating route")
            #     route = self._find_escape_route()
            #     if route != None and len(route) > 0:
            #         route.pop(0)
            #     self.escape_queue = route if route is not None else []
            # if len(self.escape_queue) != 0:
            #     print("escaping",p_loc,self.escape_queue)
            #     next_loc = self.escape_queue.pop(0)
            #     # if (p_loc == next_loc) :
            #     #     self.escape_queue.pop(0)
            #     # else:
            #     self._send_socket_command_to_target(p_loc,next_loc)
            #     #self._send_command_message_to_target(p_loc,self.escape_queue.pop(0))
            #     return
            # global FRIGHTENED_GHOST_WEIGHT 
            # if self.num_powerup <= 2:
            #     FRIGHTENED_GHOST_WEIGHT = 2 * GHOST_WEIGHT
            # else:
            #     FRIGHTENED_GHOST_WEIGHT = .3 * GHOST_WEIGHT
            
            # print(self.num_powerup,FRIGHTENED_GHOST_WEIGHT)
            
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
            curr_time = time.time()
            print(curr_time - last_time)
            last_time = curr_time
            await asyncio.sleep(0.001)
