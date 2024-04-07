import json
from json import JSONEncoder
import numpy as np

from gameState import GameState
from debugServer import DebugServer
from utils import get_walkable_tiles, get_distance
from DistMatrix import loadDistTable, loadDistTableDict


class cellAvoidanceMap:
    def __init__(self, g: GameState):
        """
        Creates instance of cellAvoidanceMap for a given GameState.
        Tunable parameters: ghost_proximity, pellet_boost
        """
        self.avoidance_map = {}
        self.g = g
        self.pellet_boost = 50
        self.superPellet_boost = 200

        self.distTable = loadDistTable()
        self.dtDict = loadDistTableDict()
        
        self.updateMap(self.g)
        
    
    def updateMap(self, g: GameState):
        """
        Reset map and ghost values
        @param: 
            - g, GameState object
        """
        self.g = g
        self.avoidance_map = {}

        print(self.dtDict)
        print(type(self.dtDict))

        # for ghost in self.g.ghosts:
        #     self.ghosts[ghost.color] = (ghost.location.row, ghost.location.col)
        
        self.ghost_positions = list(map(lambda ghost: (ghost.location.row, ghost.location.col), g.ghosts))

        for tile in get_walkable_tiles(g):
            ghost_proximity = 0
            for ghost_pos in self.ghost_positions:
                # TODO: Account for ghost color (i.e. avoid red ghost more than pink?)
                try:
                    tile_idx = self.dtDict[tile]
                    ghost_idx = self.dtDict[ghost_pos]
                    dist = self.distTable[tile_idx][ghost_idx]
                except IndexError:
                    dist = get_distance(tile, (ghost_pos[0], ghost_pos[1]))
                except KeyError:
                    dist = get_distance(tile, (ghost_pos[0], ghost_pos[1]))
                dist = get_distance(tile, ghost_pos)
                
                #dist = get_astar_dist(tile, ghost_pos, self.g)
                if dist == 0 or dist is None:
                    ghost_proximity += 1000  # Tunable
                else:
                    ghost_proximity += 1 / dist * 500  # Tunable

            # TODO: Maybe account for distance to nearby pellets?
            pellet_boost = 0
            if self.g.pelletAt(tile[0], tile[1]):
                pellet_boost = self.pellet_boost
            if self.g.superPelletAt(tile[0], tile[1]):
                pellet_boost = self.superPellet_boost

            self.avoidance_map[tile] = ghost_proximity - pellet_boost


# Original cell_avoidance_map code
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