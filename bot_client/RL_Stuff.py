import math
from gameState import GameState
from debugServer import DebugServer
import pathfinding
import numpy as np

class RLLearn_SARAS():
    def __init__(self, addr, port, training=False):
        #training     - whetyer or not its training
        self.training = training
        # alpha       - learning rate
        # epsilon     - exploration rate
        # gamma       - discount factor
        # numTraining - number of training episodes

        # let's think about what Q, the 2D array of learned state-action values, should look like
        # Pacman has 4 possible actions: up, down, left, right

        self.states = self.create_state_list()
        # The state includes various information,
        # such as the position of pacman, the position of the ghosts (or lack thereof, cuz they're dead),
        # the position of pellets, the frightened state of the ghosts

        # 28 * 36 = 1008 possible positions for pacman
        # 28 * 36 * 2 + 1 = 2017 possible states for the ghosts (position, frightened, or dead)
        # 28 * 36 = 1008 possible positions for each pellet
        # add em all up: 1008 + 4*2017 + 1008 = 10084
        # that's kind of a lot, maybe let's rethink the state space (after all, not all locations are reachable)

        # reference: https://github.com/wrhlearner/PacBot-2023/blob/master/src/Pi/botCode/HighLevelMarkov.py


            
        # a dictionary for storing Q(s,a)
        # a list records last state
        # a list records last action
        # a variable stores the score before last action

    def create_state_list():
        return []
        #this function should calculate all the legal states on the board that the pacman can be in, 
        #keeping in mind the 
        #returns a list of states
    
    def q_mapper(GameState):
        # figure out where a state is in the Q table
        #returns a dictionary mapping rewards to each state
        return

    
    def stepping(state, action):
        #Updating the new state, the reward for the step, whether pacman is done or not
        #should call 
        next_state = state
        reward = 0
        done = False
        return next_state, reward, done
    
    def get_possible_actions(state):
        #given a state, return a list of possible states
        return
    
    def get_action_random():
        #returns a random action from the state space
        return
        
    def get_action_greedy():
        return
        
    def get_action_epsilon():
        return
    
    def calculate_reward(state1, state2, action):
        #takes a board state, the next state and the action and calculates the ending rewards 
        return 0
    
    def action_to_command():
        return
    
    def train(train=True):
        
        # initialize Q(s,a)
        # take a random action
        # update Q(s,a)
        # choose the action maximises Q or a random action according to Ɛ-greedy function
        # repeat step 3 and 4 until the game ends
        # update Q(s,a) where s is the last state before the end, a is the last action taken
        #when you are at the last state, print out the statistics -> call evaluate()
        #if train is False, do not ever do e greedy, always choose the greedy actions

        return

    def evaluate(max_steps, episodes, train=True):
        #print out the average reward, what is the average reward for the episdoes, 
        #how many episdoes sucess, how many died
        #how many steps before death on average (maybe print the whole array)
        return  

    def update_values():
        #every time an episode ends, what is
        return

    