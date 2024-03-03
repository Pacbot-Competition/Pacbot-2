import math
from gameState import GameState
from debugServer import DebugServer
import heapq

class PriorityQueue:
    def __init__(self):
        self._queue = []
        self._index = 0

    def push(self, item, priority):
        heapq.heappush(self._queue, (priority, self._index, item))
        self._index += 1

    def pop(self):
        return heapq.heappop(self._queue)[-1]
    
    def empty(self):
        return len(self._queue) == 0

def get_distance(posA, posB):
    rowA,colA = posA if type(posA) == tuple else (posA.row, posA.col)
    rowB,colB = posB if type(posB) == tuple else (posB.row, posB.col)

    drow = rowA - rowB
    dcol = colA - colB
    dist = math.sqrt(dcol * dcol + drow * drow)
    return dist

def estimate_heuristic(node_pos, target_pos, cell_avoidance_map):
    return get_distance(node_pos, target_pos) + (cell_avoidance_map[node_pos] if cell_avoidance_map is not None else 0)

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

def get_walkable_tiles(g: GameState):
	walkable_cells = set()
	for row in range(31):
		for col in range(28):
			if not g.wallAt(row, col):
				walkable_cells.add((row, col))
	return walkable_cells

def build_cell_avoidance_map(g: GameState):
    cell_avoidance_map = {}

    ghost_positions = list(map(lambda ghost: (ghost.location.row, ghost.location.col), g.ghosts))

    for tile in get_walkable_tiles(g):
        ghost_proximity = 0
        for ghost_pos in ghost_positions:
            dist = get_distance(tile, ghost_pos)
            if dist == 0:
                ghost_proximity += 1000
            else:
                ghost_proximity += 1 / dist * 500

        pellet_boost = 0
        if g.pelletAt(tile[0], tile[1]):
            pellet_boost = 50
        if g.superPelletAt(tile[0], tile[1]):
            pellet_boost = 200

        cell_avoidance_map[tile] = ghost_proximity - pellet_boost

    return cell_avoidance_map

def show_cell_avoidance_map(cell_avoidance_map):
    new_cell_colors = []
    for cell, score in cell_avoidance_map.items():
        score = min(max(-255, score), 255)
        color = (score, 0, 0) if score > 0 else (0, -score, 0)
        new_cell_colors.append((cell, color))

    DebugServer.instance.set_cell_colors(new_cell_colors)

def find_path(start, target, g: GameState):
    cell_avoidance_map = build_cell_avoidance_map(g)
    show_cell_avoidance_map(cell_avoidance_map)

    frontier = PriorityQueue()
    frontier.push(start, 0)
    expanded = []
    reached = {start: {"cost": estimate_heuristic(start, target, cell_avoidance_map), "parent": None}}
    path = []

    while not frontier.empty():
        # Pop highest priority from the frontier
        currentNode = frontier.pop()

        # If current node is target, retrace path
        if currentNode == target:
            retrace = currentNode
            path = []
            while retrace is not start:
                path.append(retrace)
                retrace = reached[retrace]["parent"]
            path.reverse()
            print(path)
            return path

        # Add current, non-goal node to the expanded list
        expanded.append(currentNode)
        
        # Add neighboring nodes to the frontier
        neighbors = get_neighbors(g, currentNode)
        for neighbor in neighbors:
            # Get cumulative cost, g
            tentative_gScore = reached[currentNode]['cost'] + get_distance(currentNode, neighbor)
            
            if neighbor not in reached:
                reached[neighbor] = {"cost": tentative_gScore, "parent": currentNode}
                frontier.push(neighbor, (tentative_gScore + estimate_heuristic(neighbor, target, cell_avoidance_map)))
    
    return None

if __name__ == '__main__':
    g = GameState()

    start = (1,1)
    target = (6,6)

    path = find_path(start, target, g)
    #print(path)