import cv2
from cv2 import aruco
import matplotlib.pyplot as plt
from ipywidgets import interact                                                      # type: ignore

# A dictionary of 4x4 ArUco markers
dictionary = aruco.getPredefinedDictionary(cv2.aruco.DICT_4X4_250)

# Method to generate an ArUco marker, then display it to the screen
# @interact(id = (0, 10))                                                              # type: ignore
def generateShow(id: int) -> None:
    marker_image = aruco.generateImageMarker(dictionary, id, 6)
    plt.imshow(marker_image, cmap='gray')                                            # type: ignore
    plt.show()                 

for id in range(7):
    marker_image = aruco.generateImageMarker(dictionary, id, 6)
    plt.imshow(marker_image, cmap='gray')        
    plt.savefig('aruco/aruco_' + str(id) + '.png')

"""
id 0 is the Pacman
id 1 is at 0,0
id 2 is at 0,28
id 3 is at 16,0
id 4 is at 16,28
id 5 is at 31,0
id 6 is at 31,28
"""