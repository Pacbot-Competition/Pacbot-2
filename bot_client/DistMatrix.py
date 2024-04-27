
import json
from json import JSONEncoder
import numpy as np

from utils import get_astar_dist, get_walkable_tiles


class NumpyArrayEncoder(JSONEncoder):
    """
    Encode numpy arrays to JSON
    """
    def default(self, obj):
        if isinstance(obj, np.ndarray):
            return obj.tolist()
        return JSONEncoder.default(self, obj)


def createDistTable(g):
    """
    Creates a table of distances between all walkable tiles in the game.
    Includes ghost spawn tiles.
    @param: 
        - g, GameState object
    @return: 
        - json representation of np.ndarray
    """
    walkable_tiles = get_walkable_tiles(g)
    distTable = np.zeros((len(walkable_tiles), len(walkable_tiles)))

    ghost_spawn_tiles = set()
    for i in range(13, 16):
            for j in range(11, 17):
                ghost_spawn_tiles.add((i, j))
    
    my_tiles = walkable_tiles.union(ghost_spawn_tiles)

    # Create distance matrix
    for tile1 in my_tiles:
        for tile2 in my_tiles:
            dist = get_astar_dist(tile1, tile2, g)
            distTable[tile1, tile2] = dist

    # Write to json file
    json_distTable = {"distTable": distTable}
    with open('static/distTable.json', 'w') as f:
        json.dump(json_distTable, f, cls=NumpyArrayEncoder)

    print("done")


def loadDistTable():
    """
    Load the distance table from the JSON file.
    @return: np.ndarray, the distance table
    """
    with open('static/distTable.json', 'r') as f:
        data = json.load(f)
    return data["distTable"]


def createDistTableDict(g):
    """"
    Create a dictionary of {tile tuples: indices} for the distance table.
    @param: 
        - g, GameState object
    @return:
        - json representation of dict
    """
    # Get ghost spawn tiles
    ghost_spawn_tiles = set()
    for i in range(13, 16):
            for j in range(11, 17):
                ghost_spawn_tiles.add((i, j))

    # All walkable tiles
    my_tiles = get_walkable_tiles(g).union(ghost_spawn_tiles)

    # Create dictionary
    dtDict = {}
    for idx, tile in enumerate(my_tiles):
        tuple_str = ','.join(map(str, tile))
        dtDict[tuple_str] = idx

    with open('static/dtDict.json', 'w') as f:
        json.dump(dtDict, f)


def loadDistTableDict():
    """
    Load the distance table dictionary from the JSON file.
    @return: 
        - dict, the distance table dictionary
    """
    with open('static/dtDict.json', 'r') as f:
        data = json.load(f)

        # Convert string keys back to tuples
        converted_data = {}
        for key, value in data.items():
            # Split the string key by comma and convert elements to integers
            key_tuple = tuple(map(int, key.split(',')))
            converted_data[key_tuple] = value
    
    return converted_data