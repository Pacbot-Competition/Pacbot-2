import math
from gameState import GameState

def get_distance(posA, posB):
    dy = posB[0] - posA[0]
    dx = posB[1] - posA[1]
    dist = math.sqrt(dx * dx + dy * dy)
    return dist

def get_neighbors(g: GameState):
    pacman_loc = g.pacmanLoc
    y,x = (pacman_loc.row,pacman_loc.col)
    dirs = []
    if not g.wallAt(y+1,x):
        dirs.append(((y+1,x),'s'))
    if not g.wallAt(y-1,x):
        dirs.append(((y-1,x),'w'))
    if not g.wallAt(y,x+1):
        dirs.append(((y,x+1),'d'))
    if not g.wallAt(y,x-1):
        dirs.append(((y,x-1),'a'))
    print(dirs)

class node:
    def __init__ (self, f, g, parent_node=None):
        self.f = f
        self.g = g
        self.h = f + g
        self.parent_node = parent_node

    def get_f():
        return f

def pathfind(start, end):
    open_cells = []
    closed_cells = []
    path = []

    open_cells.append(start)
    while len(open_cells) > 0:
        current_node = open_cells[0]
        min_f = current_node.get_f()

        for i in range(0, len(open_cells)):
            if open_cells[i].get_f < min_f:
                current_node = open_cells[i]
                min_f = open_cells[i].get_f
        
        open_cells.remove(current_node)
        closed_cells.append(current_node)

        if current_node == end:


            
