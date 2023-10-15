import math
from gameState import GameState

def get_distance(posA, posB):
    dx = posB[0] - posA[0]
    dy = posB[1] - posA[1]
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
