import math
import os
from gameState import GameState
from debugServer import DebugServer
import json
from json import JSONEncoder
import numpy as np

from utils import PriorityQueue, get_distance, get_neighbors
from AvoidanceMap import cellAvoidanceMap, show_cell_avoidance_map


def estimate_heuristic(node_pos, target_pos, cell_avoidance_map):
    """Heuristic is based on Euclidean distance between node and target, plus cell avoidance map value."""
    if (node_pos == (32, 32)):
        return 0
    return get_distance(node_pos, target_pos) + (cell_avoidance_map[node_pos] if cell_avoidance_map is not None else 0)


def find_path(start, target, g: GameState):
    """
    Current Pac-Man policy: A* search to find path from start to target using cell avoidance map.
    @param:
        - start: tuple, (row, col) of the starting point
        - target: tuple, (row, col) of the target point
        - g: GameState object
    @return:
        - list, the path from start to target
    """
    map_class = cellAvoidanceMap(g)
    map_class.updateMap(g)
    cell_avoidance_map = map_class.avoidance_map

    #cell_avoidance_map = build_cell_avoidance_map(g)
    show_cell_avoidance_map(cell_avoidance_map)

    print(f'start: {start}, target: {target}')

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
            g_score = reached[currentNode]['cost'] + get_distance(currentNode, neighbor)
            
            if neighbor not in reached:
                h_score = estimate_heuristic(neighbor, target, cell_avoidance_map)
                reached[neighbor] = {"cost": g_score, "parent": currentNode}
                frontier.push(neighbor, (g_score + h_score))
    
    return (0, 0)


if __name__ == '__main__':
    g = GameState()

    start = (1,1)
    target = (6,6)

    path = find_path(start, target, g)
    #print(path)