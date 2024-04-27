from gameState import GameState
import json

# # Load dtDict.json
# with open('dtDict.json', 'r') as f:
#     dt_dict = json.load(f)

# # Load distTable.json
# with open('distTable.json', 'r') as f:
#     dist_table = json.load(f)

class MazeGraph:
    def __init__(self, g: GameState):
        self.g = g
        self.intersections = set() # list of intersections
        self.graph = {}
        self.createGraph()
    
    def isIntersection(self, row: int, col: int) -> bool:
        above = self.g.wallAt(row, col+1)
        below = self.g.wallAt(row, col-1)
        left = self.g.wallAt(row-1, col)
        right = self.g.wallAt(row+1, col)
        # 4 way intersection
        if above and below and left and right:
            return True
        # 3 way intersection
        if above and below and left:
            return True
        if above and below and right:
            return True
        if above and left and right:
            return True
        if below and left and right:
            return True
        # 2 way intersection
        if above and right:
            return True
        if above and left:
            return True
        if below and right:
            return True
        if below and left:
            return True
        return False
    
    def createGraph(self):
        # create a graph (adjacency list representation) of all the intersections, and distances and values between 
        # graph = {
        #   (row_1, col_1): {(row_n, col_n): [distance, value], 
        #                   (row_k, col_k): [distance, value], ...}, ...
        # } 
        # find and create a list of all intersections
        for row in range(self.g.height):
            for col in range(self.g.width):
                if not self.g.wallAt(row, col):
                    if self.isIntersection(row, col):
                        self.intersections.add((row, col))
        
        for intersect in self.intersections:
            directions = [(0, 1), (0, -1), (-1, 0), (1, 0)]
            for i in directions:
                dist = 0
                row = intersect[0]
                col = intersect[1]
                while not self.g.wallAt(row, col): # keep moving in the specified direction until you hit a wall or another intersection
                    row += i[0] 
                    col += i[1]
                    dist += 1
                    if (row, col) in self.intersections:
                        if intersect not in self.graph:
                            self.graph[intersect] = {(row, col): [dist, 0]}
                        else:
                            self.graph[intersect][(row, col)] =  [dist, 0]
                        break

    def updateGhostValue(self, intersect1, intersect2, new_value):
        self.graph[intersect1][intersect2][1] = new_value
        self.graph[intersect2][intersect1][1] = new_value