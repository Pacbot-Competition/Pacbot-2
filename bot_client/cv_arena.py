import numpy as np
import cv2 as cv
from typing import Tuple

# Variables used in the arena grid
I = 1  # wall
o = 2  # pellet
e = 3  # empty space
O = 4  # powerpoint
n = 5  # untouchable
P = 6  # pacbot
v = 7  # visited
arena = np.rot90([
    [I, I, I, I, I, I, I, I, I, I, I, I, e, e, e, I,
     I, I, e, e, e, I, I, I, I, I, I, I, I, I, I],
    [I, o, o, o, o, I, I, O, o, o, o, I, e, e, e, I,
     v, I, e, e, e, I, o, o, o, o, o, O, o, o, I],
    [I, o, I, I, o, I, I, o, I, I, o, I, e, e, e, I,
     v, I, e, e, e, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, o, o, o, I, I, o, I, e, e, e, I,
     v, I, e, e, e, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, I, I, I, o, I, e, e, e, I,
     v, I, e, e, e, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, I, I, I, o, I, I, I, I, I,
     I, I, I, I, I, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, o, o, o, o, o, o, o, o, o, o, o,
     o, o, o, o, o, o, o, o, o, o, o, o, o, o, I],
    [I, o, I, I, I, I, I, o, I, I, o, I, I, I, I, I,
     v, I, I, I, I, I, I, I, I, o, I, I, I, o, I],
    [I, o, I, I, I, I, I, o, I, I, o, I, I, I, I, I,
     v, I, I, I, I, I, I, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, o, o, o, I, I, o, v, v, v, v, v,
     v, v, v, v, I, I, o, o, o, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, o, I, I, o, I, I, v, I, I,
     I, I, I, v, I, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, o, I, I, o, I, I, v, I, n,
     n, n, I, v, I, I, o, I, I, o, I, I, I, o, I],
    [I, o, o, o, o, I, I, o, o, o, o, I, I, v, I, n,
     n, n, I, v, v, v, o, I, I, o, o, o, o, o, I],
    [I, o, I, I, I, I, I, v, I, I, I, I, I, v, I, n,
     n, n, n, v, I, I, I, I, I, o, I, I, I, I, I],
    [I, o, I, I, I, I, I, v, I, I, I, I, I, v, I, n,
     n, n, n, v, I, I, I, I, I, o, I, I, I, I, I],
    [I, o, o, o, o, I, I, o, o, o, o, I, I, v, I, n,
     n, n, I, v, v, v, o, I, I, o, o, o, o, o, I],
    [I, o, I, I, o, I, I, o, I, I, o, I, I, v, I, n,
     n, n, I, v, I, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, o, I, I, o, I, I, v, I, I,
     I, I, I, v, I, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, o, o, o, I, I, o, v, v, v, v, v,
     v, v, v, v, I, I, o, o, o, o, I, I, I, o, I],
    [I, o, I, I, I, I, I, o, I, I, o, I, I, I, I, I,
     v, I, I, I, I, I, I, I, I, o, I, I, I, o, I],
    [I, o, I, I, I, I, I, o, I, I, o, I, I, I, I, I,
     v, I, I, I, I, I, I, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, o, o, o, o, o, o, o, o, o, o, o,
     o, o, o, o, o, o, o, o, o, o, o, o, o, o, I],
    [I, o, I, I, o, I, I, I, I, I, o, I, I, I, I, I,
     I, I, I, I, I, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, I, I, I, o, I, e, e, e, I,
     v, I, e, e, e, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, o, o, o, I, I, o, I, e, e, e, I,
     v, I, e, e, e, I, o, I, I, o, I, I, I, o, I],
    [I, o, I, I, o, I, I, o, I, I, o, I, e, e, e, I,
     v, I, e, e, e, I, o, I, I, o, I, I, I, o, I],
    [I, o, o, o, o, I, I, O, o, o, o, I, e, e, e, I,
     v, I, e, e, e, I, o, o, o, o, o, O, o, o, I],
    [I, I, I, I, I, I, I, I, I, I, I, I, e, e, e, I, I, I, e, e, e, I, I, I, I, I, I, I, I, I, I]], k=2)


direction = {
    'right': (1, 0),
    'down': (0, 1),
    'left': (-1, 0),
    'up': (0, -1),
}


def wall_correction(pac_pos: Tuple[int, int]) -> Tuple[int, int]:
    '''
    Takes in an incorrect node coordinate in a tuple (x,y) and returns the closest valid node coordinate instead.
    Valid nodes include: pellet, powerpoint, and visited (see constants above)
    The incorrect node is most likely caused by the pacbot being located inside a wall node.
    '''
    # TODO: implement this function
    # TODO: implement this function
    # Check if the current position is already valid

    offsets = [(x, y) for x in range(-1, 2) for y in range(-1, 2)]
    best_dist = 100
    best_pos = (-1, -1)
    for offset in offsets:
      new_pos = (int(pac_pos[0] + offset[0]), int(pac_pos[1] + offset[1]))
      dist = abs(pac_pos[0] - (new_pos[0] + 0.5)) + abs(pac_pos[1] - (new_pos[1] + 0.5))
      if dist < best_dist and 0 <= new_pos[0] < len(arena) and 0 <= new_pos[1] < len(arena[0]) and arena[new_pos[0], new_pos[1]] in [o, O, v]:
        best_dist = dist
        best_pos = new_pos

    return best_pos

    if arena[pac_pos[0], pac_pos[1]] in [o, O, v]:
        return pac_pos
    
    #check 8 surrounding positions
    if 0 <= pac_pos[0] < arena.shape[0] and 0 <= pac_pos[1] + 1< arena.shape[1] and arena[pac_pos[0], pac_pos[1] + 1] in [o, O, v]:
      search_pos = (pac_pos[0], pac_pos[1] + 1)
      return search_pos

    elif 0 <= pac_pos[0] - 1 < arena.shape[0] and 0 <= pac_pos[1]< arena.shape[1] and arena[pac_pos[0] - 1, pac_pos[1]] in [o, O, v]:
      search_pos = (pac_pos[0] - 1 , pac_pos[1])
      return search_pos
    
    elif 0 <= pac_pos[0] + 1 < arena.shape[0] and 0 <= pac_pos[1]< arena.shape[1] and arena[pac_pos[0] + 1, pac_pos[1]] in [o, O, v]:
      search_pos = (pac_pos[0] +1, pac_pos[1])
      return search_pos

    elif 0 <= pac_pos[0] < arena.shape[0] and 0 <= pac_pos[1] - 1< arena.shape[1] and arena[pac_pos[0], pac_pos[1] - 1] in [o, O, v]:
      search_pos = (pac_pos[0], pac_pos[1] - 1)
      return search_pos

    elif 0 <= pac_pos[0] -1 < arena.shape[0] and 0 <= pac_pos[1] + 1< arena.shape[1] and arena[pac_pos[0] - 1, pac_pos[1] + 1] in [o, O, v]:
      search_pos = (pac_pos[0] -1, pac_pos[1] + 1)
      return search_pos
    
    elif 0 <= pac_pos[0] + 1 < arena.shape[0] and 0 <= pac_pos[1] + 1 < arena.shape[1] and arena[pac_pos[0] + 1, pac_pos[1] + 1] in [o, O, v] :
      search_pos = (pac_pos[0] + 1, pac_pos[1] + 1)
      return search_pos
    
    elif 0 <= pac_pos[0] - 1 < arena.shape[0] and 0 <= pac_pos[1] - 1< arena.shape[1] and arena[pac_pos[0] -1 , pac_pos[1] - 1] in [o, O, v]:
      search_pos = (pac_pos[0] -1, pac_pos[1] - 1)
      return search_pos
    
    
    elif 0 <= pac_pos[0] + 1 < arena.shape[0] and 0 <= pac_pos[1] - 1< arena.shape[1] and arena[pac_pos[0] + 1, pac_pos[1] - 1] in [o, O, v]:
      search_pos = (pac_pos[0]  + 1, pac_pos[1] - 1)
      return search_pos
    

    return (-1,-1)


def trace_missing_path(start_pos: Tuple[int, int], end_pos: Tuple[int, int]) -> list[Tuple[int, int]]:
    '''
    Takes in two not continuous node coordinates (it takes more than 1 action to get from node 1 to node 2) and returns the most likely path the robot has taken to reach the end node from the start node. 
    This is required in cases where pacbot moves almost diagonally at corners and confuses the system, making it look like pacbot suddenly moved two nodes in one move. 
    Returned path is a list of node traversed during the movement.
    '''
    path_list = [start_pos]
    current_pos = start_pos
    while current_pos != end_pos:
        possible_locs = get_possible_next_loc(current_pos)
        best_loc = (-1,-1)
        best_dist = 1000000
        for loc in possible_locs:
            dist = euclidian_dist_squared(end_pos, loc)
            if dist < best_dist:
                best_loc = loc
                best_dist = dist
        path_list.append(best_loc)
        current_pos = best_loc
    return path_list

def get_possible_next_loc(pac_pos:Tuple[int,int]) -> list[Tuple[int,int]]:
    possible_locs = []
    position = pac_pos
    for key in direction:
        new_pos = add_tuples(position, direction[key])
        if arena[new_pos[0]][new_pos[1]] == I or arena[new_pos[0]][new_pos[1]] == n or arena[new_pos[0]][new_pos[1]] == e:
            continue
        possible_locs.append(new_pos)
    return possible_locs

def add_tuples(a:Tuple[int,int],b:Tuple[int,int])->Tuple[int,int]:
    return (a[0]+b[0],a[1]+b[1])

def euclidian_dist_squared(a:Tuple[int,int], b:Tuple[int,int]):
    return (a[0]-b[0])**2 + (a[1] - b[1])**2

if __name__ == "__main__":
    print(arena)
