
import bitstruct.c as bitstruct
from gameState import Directions, Location, GameState


class Node:
    def __init__(self, loc: Location, dist: int):
        self.loc = loc
        self.dist = dist

def comp_locations(locA: Location, locB: Location) -> bool:
    return locA.row == locB.row and locA.col == locB.col

def comp_location_pair(locA1: Location, locA2: Location, locB1: Location, locB2: Location) -> bool:
    return (comp_locations(locA1, locB1) and comp_locations(locA2, locB2)) or \
            (comp_locations(locA1, locB2) and comp_locations(locA2, locB1))

def getKey(loc1: Location, loc2: Location) -> int:
    r1 = loc1.row
    c1 = loc1.col
    r2 = loc2.row
    c2 = loc2.col

    if ((r1, c1) > (r2, c2)):
        r1, c1, r2, c2 = r2, c2, r1, c1
    return int.from_bytes(bitstruct.pack('u5u5u5u5', r1, c1, r2, c2), "big")

# use BFS to get dist btwn loc
def getDistance(loc: Location, state: GameState, dist_dict, count):

    # BFS queue
    firstNode = Node(loc, 0)
    queue = [firstNode]
    visited = set()
    if str(firstNode.loc) + " " + str(firstNode.loc.row) not in visited:
        count += 1
        # dist_dict[getKey(loc, loc)] = (0, loc, loc)
        dist_dict[getKey(loc, loc)] = 0
        visited.add(str(firstNode.loc))
    

    while queue:
        # pop from queue
        currNode = queue.pop(0)

        # Loop over the directions
        for direction in Directions:

            # If the direction is none, skip it
            if direction == Directions.NONE:
                continue

            # get next location (deep copy)
            nextLoc = Location(state)
            nextNode = Node(nextLoc, currNode.dist + 1)
            nextNode.loc.col = currNode.loc.col
            nextNode.loc.row = currNode.loc.row
            nextNode.loc.setDirection(direction)
            valid = nextNode.loc.advance()

            # avoid same node twice
            # check this is a valid move
            if str(nextNode.loc) not in visited and valid:
                count += 1
                # dist_dict[getKey(loc, nextNode.loc)] = (nextNode.dist, loc, nextNode.loc)
                dist_dict[getKey(loc, nextNode.loc)] = nextNode.dist
                queue.append(nextNode)
                visited.add(str(nextNode.loc))
    return count


def main():

    state: GameState = GameState()

    dist_dict = {}
    count = 0
    wall_count = 0


    for col in range(28):
        for row in range(31):
            if state.wallAt(col=col, row=row):
                wall_count += 1
                continue
            loc: Location = Location(state)
            loc.col = col
            loc.row = row
            count = getDistance(loc, state, dist_dict, count)
    print(dist_dict)

        
    

# main()