from gameState import GameState

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