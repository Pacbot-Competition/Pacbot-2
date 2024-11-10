import numpy as np
from gameState import GhostColors, GameState  # Import relevant elements
from debugServer import DebugServer
from utils import get_walkable_tiles, get_distance
from DistMatrix import loadDistTable, loadDistTableDict

class cellAvoidanceMap:
    def __init__(self, g: GameState):
        """
        Creates instance of cellAvoidanceMap for a given GameState.
        Tunable parameters: ghost_proximity, pellet_boost, superPellet_boost.
        """
        self.avoidance_map = {}
        self.g = g
        self.pellet_boost = 50
        self.superPellet_boost = 200
        self.num_pellets = g.numPellets()

        # Distance table and dictionary for cached distances (if available)
        self.distTable = loadDistTable()
        self.dtDict = loadDistTableDict()

        self.updateMap(self.g)

    def calculate_boosts(self, tile):
        """
        Calculate the boost for normal pellets and super pellets based on the remaining pellets.
        """
        base_boost = self.pellet_boost

        # Super pellet boost is ignored if there are more than 150 pellets
        super_pellet_boost = self.superPellet_boost if self.num_pellets <= 150 else 0

        # Aggressively boost pellet collection when few pellets remain
        if self.num_pellets < 5:
            pellet_boost = base_boost * 25
        elif self.num_pellets < 10:
            pellet_boost = base_boost * 10
        elif self.num_pellets < 20:
            pellet_boost = base_boost * 7
        elif self.num_pellets < 50:
            pellet_boost = base_boost * 4
        elif self.num_pellets < 100:
            pellet_boost = base_boost * 2
        else:
            pellet_boost = base_boost

        # Check if this tile has a pellet or super pellet and return appropriate boost
        if self.g.superPelletAt(tile[0], tile[1]):
            return super_pellet_boost
        elif self.g.pelletAt(tile[0], tile[1]):
            return pellet_boost
        return 0

    def calculate_ghost_proximity(self, tile, ghost):
        """
        Calculate the proximity influence for a ghost at a given tile based on distance and ghost color.
        """
        dist = get_distance(tile, (ghost.location.row, ghost.location.col))
        fright_modifier = -1 if ghost.isFrightened() else 1

        # Use a smaller avoidance distance if fewer than 10 pellets remain
        threshold_dist = 4 if self.num_pellets < 10 else 8

        # Adjust avoidance based on ghost color
        if dist < threshold_dist:
            avoidance_weight = 1.0
            if ghost.color == GhostColors.RED:  # Blinky: aggressive
                avoidance_weight = 1.5
            elif ghost.color == GhostColors.PINK:  # Pinky: ambush
                avoidance_weight = 1.2
            elif ghost.color == GhostColors.CYAN:  # Inky: unpredictable
                avoidance_weight = 1.1
            elif ghost.color == GhostColors.ORANGE:  # Clyde: less aggressive if close
                avoidance_weight = 1.3 if dist > 8 else 0.8

            return (1 / dist) * 250 * fright_modifier * avoidance_weight if dist > 0 else 1000 * fright_modifier
        return 0

    def updateMap(self, g: GameState):
        """
        Update the avoidance map based on ghost proximity and pellet boosts.
        Strategy:
        1. Before 150 pellets, ignore super pellets by setting their boost to 0.
        2. After 150 pellets, use super pellets as a way to escape if threatened by nearby ghosts.
        """
        self.g = g
        self.avoidance_map = {}
        self.ghosts = self.g.ghosts  # Get ghost data

        for tile in get_walkable_tiles(g):
            ghost_proximity = sum(self.calculate_ghost_proximity(tile, ghost) for ghost in self.ghosts)
            pellet_boost = self.calculate_boosts(tile)
            final_score = ghost_proximity - pellet_boost

            # If ghosts are far away, prioritize pellet collection more aggressively
            if all(get_distance(tile, (ghost.location.row, ghost.location.col)) > 16 for ghost in self.ghosts):
                final_score -= pellet_boost * 1.5  # Aggressively collect pellets

            self.avoidance_map[tile] = final_score

    def show_map(self):
        """
        Display the avoidance map on the debug server with color coding for ghost proximity and pellet boost.
        """
        new_cell_colors = []
        for cell, score in self.avoidance_map.items():
            score = min(max(-255, score), 255)  # Clamp score to avoid overflow
            color = (max(0, score), max(0, -score), 0)  # RGB: (Red, Green, Blue)
            new_cell_colors.append((cell, color))

        DebugServer.instance.set_cell_colors(new_cell_colors)
