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

    open_nodes = set()
    open_nodes.add(start)

    parents = {}

    g_map = {}
    g_map[start] = 0

    f_map = {}
    f_map[start] = estimate_heuristic(start, target, cell_avoidance_map)

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
            return tuple(path)
        
        open_nodes.remove(current)

        for neighbor in get_neighbors(g, current):
            tentative_gScore = g_map[current] + get_distance(current, neighbor)
            if neighbor not in g_map or tentative_gScore < g_map[neighbor]:
                parents[neighbor] = current
                g_map[neighbor] = tentative_gScore
                f_map[neighbor] = tentative_gScore + estimate_heuristic(neighbor, target, cell_avoidance_map)
                open_nodes.add(neighbor)

    return None

if __name__ == '__main__':
    g = GameState()

    start = (1,1)
    target = (6,6)

    path = find_path(start, target, g)
    print(path)