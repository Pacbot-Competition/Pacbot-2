import heapq
import math

from gameState import GameState

class PriorityQueue:
    """
    Priority queue implementation with heapq
    (For use in A* search)
    """
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
    

def get_walkable_tiles(g: GameState):
    """
    Get all walkable (non-wall) cells in the game grid
    @param:
        - g: GameState object
    @return:
        - set, a set of tuples representing the walkable cells
    """
    walkable_cells = set()
    for row in range(31):
        for col in range(28):
            if not g.wallAt(row, col):
                walkable_cells.add((row, col))
    return walkable_cells


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


def get_distance(posA, posB):
    """
    Get Euclidean distance between two cells
    @param:
        - posA: tuple, (row, col) of the first cell
        - posB: tuple, (row, col) of the second cell
    @return:
        - float, the Euclidean distance between the two cells
    """
    rowA,colA = posA if type(posA) == tuple else (posA.row, posA.col)
    rowB,colB = posB if type(posB) == tuple else (posB.row, posB.col)

    drow = rowA - rowB
    dcol = colA - colB
    dist = math.sqrt(dcol * dcol + drow * drow)
    return dist


def get_astar_dist(start, target, g: GameState):
    """
    Get the A* distance between two points on the grid.
    @param:
        - start: tuple, (row, col) of the starting point
        - target: tuple, (row, col) of the target point
        - g: GameState object
    @return: 
        - int, the A* distance between the two points
    """
    frontier = PriorityQueue()
    frontier.push(start, 0)
    expanded = []
    reached = {start: {"cost": get_distance(start, target), "parent": None}}

    while not frontier.empty():
        # Pop highest priority (smallest distance) from the frontier
        currentNode = frontier.pop()

        # If current node is target, retrace path
        if currentNode == target:
            retrace = currentNode
            path = []
            while retrace is not start:
                path.append(retrace)
                retrace = reached[retrace]["parent"]
            path.reverse()
            return len(path)

        # Add current, non-goal node to the expanded list
        expanded.append(currentNode)
        
        # Add neighboring nodes to the frontier
        neighbors = get_neighbors(g, currentNode)
        for neighbor in neighbors:
            # Get cumulative cost, g
            tentative_gScore = reached[currentNode]['cost'] + get_distance(currentNode, neighbor)
            
            if neighbor not in reached:
                reached[neighbor] = {"cost": tentative_gScore, "parent": currentNode}
                frontier.push(neighbor, (tentative_gScore + get_distance(neighbor, target)))
    
    return (0, 0)