import math

def get_distance(posA, posB):
    dx = posB[0] - posA[0]
    dy = posB[1] - posA[1]
    dist = math.sqrt(dx * dx + dy * dy)
    return dist
