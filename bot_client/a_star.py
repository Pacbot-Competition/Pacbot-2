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
    if not g.wallAt(y + 1,x):
        dirs.append((y + 1,x))
    if not g.wallAt(y - 1,x):
        dirs.append((y - 1,x))
    if not g.wallAt(y, x + 1):
        dirs.append((y, x + 1))
    if not g.wallAt(y, x - 1):
        dirs.append((y, x - 1))
    return dirs

class node:
    def __init__ (self, y, x, f, g, parent_node=None):
        self.y = y
        self.x = x
        self.f = f
        self.g = g
        self.h = f + g
        self.parent_node = parent_node

    def __init__ (self, y, x, parent_node=None):
        self.y = y
        self.x = x
        self.f = 0
        self.g = 0
        self.h = 0
        self.parent_node = parent_node

    def get_f(self):
        return self.f
    
    def get_y(self):
        return self.y
    
    def get_x(self):
        return self.x

def pathfind(start, end):
    open_cells = []
    closed_cells = []
    path = []

    open_cells.append(start)
    while len(open_cells) > 0:
        print("test")
        current_node = open_cells[0]
        min_f = current_node.get_f()

        for i in range(0, len(open_cells)):
            if open_cells[i].get_f() < min_f:
                current_node = open_cells[i]
                min_f = open_cells[i].get_f()
        
        open_cells.remove(current_node)
        closed_cells.append(current_node)

        if current_node == end:
            return
        
        children = [node(current_node.y + 1, current_node.x),
                    node(current_node.y - 1, current_node.x),
                    node(current_node.y, current_node.x + 1),
                    node(current_node.y, current_node.x - 1)]

        for child in children:
            if child in closed_cells:
                continue

            child.g = current_node.g + get_distance((child.x, child.y), (current_node.x, current_node.y))
            child.f = get_distance((child.x, child.y), (end.x, end.y))
            child.h = child.g + child.f

            if child in open_cells:
                if child.g > current_node.g:
                    continue
            
            open_cells.append(child)

    return 1

print(pathfind(node(1,0, None), node(5,0, None)))