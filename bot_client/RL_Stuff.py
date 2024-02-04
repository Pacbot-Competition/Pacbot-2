import math
from gameState import GameState
from debugServer import DebugServer
import pathfinding
import numpy as np

class RLLearn_SARAS():
    def __init__(self, addr, port, training=False):
        self.training = training

        # let's think about what Q, the 2D array of learned state-action values, should look like
        # Pacman has 4 possible actions: up, down, left, right

        # The state includes various information,
        # such as the position of pacman, the position of the ghosts (or lack thereof, cuz they're dead),
        # the position of pellets, the frightened state of the ghosts

        # 28 * 36 = 1008 possible positions for pacman
        # 28 * 36 * 2 + 1 = 2017 possible states for the ghosts (position, frightened, or dead)
        # 28 * 36 = 1008 possible positions for each pellet
        # add em all up: 1008 + 4*2017 + 1008 = 10084
        # that's kind of a lot, maybe let's rethink the state space (after all, not all locations are reachable)

        # reference: https://github.com/wrhlearner/PacBot-2023/blob/master/src/Pi/botCode/HighLevelMarkov.py
        return

    """
    def q_mapper(GameState g):
        # figure out where a state is in the Q table
        return
    """
    
    def stepping():
        #Updating the new state, the reward for the step, whether pacman is done or not
        return
    
    def get_action_random():
        return
        
    def get_action_greedy():
        return
        
    def get_action_epsilon():
        return
    
    def calculate_reward(state1, state2, action):
        return 0
    
    def action_to_command():
        return
    
    def train():
        return


    def evaluate():
        return

    