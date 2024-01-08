import math
from gameState import GameState
from debugServer import DebugServer

def get_distance(posA, posB):
    rowA,colA = posA if type(posA) == tuple else (posA.row, posA.col)
    rowB,colB = posB if type(posB) == tuple else (posB.row, posB.col)

    drow = rowA - rowB
    dcol = colA - colB
    dist = math.sqrt(dcol * dcol + drow * drow)
    return dist

def estimate_heuristic(node_pos, target_pos):
    return get_distance(node_pos, target_pos) # For now, just use the euclidean distance

def get_neighbors(g: GameState, location=None):
    if location is None:
        location = g.pacmanLoc
    
    row,col = location if type(location) == tuple else (location.row, location.col)
    neighbors = []
    if not g.wallAt(row + 1,col):
        neighbors.append((row + 1,col))
    if not g.wallAt(row - 1,col):
        neighbors.append((row - 1,col))
    if not g.wallAt(row, col + 1):
        neighbors.append((row, col + 1))
    if not g.wallAt(row, col - 1):
        neighbors.append((row, col - 1))
    return neighbors

def find_path(start, target, g: GameState, debug_server: DebugServer = None):
    if debug_server is not None:
        debug_server.reset_cell_colors()

    open_nodes = set()
    open_nodes.add(start)

    parents = {}

    g_map = {}
    g_map[start] = 0

    f_map = {}
    f_map[start] = estimate_heuristic(start, target)

    while len(open_nodes) > 0:
        # Find the node with the lowest f score
        current = None
        current_f = None
        for node, score in f_map.items():
            if node in open_nodes and (current is None or score < current_f):
                current = node
                current_f = score

        if current == target:
            path = []
            while current in parents:
                path.append(current)
                current = parents[current]
            path.reverse()
            return path
        
        open_nodes.remove(current)
        if debug_server is not None:
            debug_server.set_cell_color(current[0], current[1], 'orange')

        for neighbor in get_neighbors(g, current):
            tentative_gScore = g_map[current] + get_distance(current, neighbor)
            if neighbor not in g_map or tentative_gScore < g_map[neighbor]:
                parents[neighbor] = current
                g_map[neighbor] = tentative_gScore
                f_map[neighbor] = tentative_gScore + estimate_heuristic(neighbor, target)
                open_nodes.add(neighbor)

                if debug_server is not None:
                    debug_server.set_cell_color(neighbor[0], neighbor[1], 'yellow')

    return None

if __name__ == '__main__':
    g = GameState()

    start = (1,1)
    target = (6,6)

    path = find_path(start, target, g)
    print(path)