# Asyncio (for concurrency)
import asyncio

# Import connection state object
from connectionState import ConnectionState

# Import the wall array
from walls import wallArr

# OpenCV
import cv2

# ArUco
from cv2 import aruco

# Numpy
import numpy as np

# Plt
import matplotlib.pyplot as plt

# Typing
from typing import Any

# Typedef
MatLike = cv2.typing.MatLike
IntArray = np.ndarray[Any, np.dtype[np.intp]]

class CameraModule:
	'''
	Sample implementation of a decision module for computer vision
	for Pacbot, using asyncio.
	'''

	def __init__(self, state: ConnectionState) -> None:
		'''
		Construct a new decision module object
		'''

		# Game state object to store the game information
		self.state = state

		# A dictionary of 4x4 ArUco markers
		self.dictionary = aruco.getPredefinedDictionary(cv2.aruco.DICT_4X4_250)

		# Instantiate a new ArUco detector
		self.detector = aruco.ArucoDetector(self.dictionary, aruco.DetectorParameters())

	async def decisionLoop(self) -> None:
		'''
		Decision loop for CV
		'''

		# Receive values as long as we have access
		while self.state.isConnected():

			# Get a frame
			img = self.capture()

			# Process the frame
			pacman_row, pacman_col = self.localize(img)

			# Write back to the server, as a test (move right)
			self.state.send(pacman_row, pacman_col)

			# Free up the event loop
			await asyncio.sleep(0)

	def capture(self) -> MatLike:
		'''
		Capture an image
		'''

		cap = cv2.VideoCapture(1)
		_, img = cap.read()
		img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
		return img

	def wallAt(self, row: int, col: int) -> bool:
		'''
		Helper function to check if a wall is at a given location
		'''

		# Check if the position is off the grid, and return true if so
		if (row < 0 or row >= 31) or (col < 0 or col >= 28):
			return True

		# Return whether there is a wall at the location
		return bool((wallArr[row] >> col) & 1)

	def localize(self, img: MatLike, warp: bool = False, annotate: bool = False) -> tuple[int, int]:

		# Detect markers
		corners, ids, _ = self.detector.detectMarkers(img)

		# Array of ids with centroids
		ids_centroids: list[tuple[int, IntArray]] = []

		# Variable for whether Pacman was found in frame
		foundPacman = False

		# Loop over the ids
		for j in range(len(ids)):

			# Find this id
			id = ids[j, 0]

			# If the id is invalid, skip it
			if id > 6:
				continue

			# If the id is 0, Pacman has been found
			if id == 0:
				foundPacman = True

			# Find the coordinates of this centroid
			centroid = np.array([
				int(corners[j][0][:, 0].mean()),
				int(corners[j][0][:, 1].mean())
			])

			# Put these together as a pair
			pair = (id, centroid)

			# Find the coordinates of each centroid
			ids_centroids.append(pair)

		# Assert that Pacman was found
		assert foundPacman, "ERR: Pacman not found"

		# Sort the centroids
		ids_centroids.sort()

		# Get the sorted ids
		ids, centroids = list(zip(*ids_centroids))

		# Determine if the region is the top half or the bottom half
		topHalf = (ids == (0, 1, 2, 3, 4))
		bottomHalf = (ids == (0, 3, 4, 5, 6))

		# Assert that we're either in the top half or bottom half
		assert (topHalf or bottomHalf), "ERR: The image is neither the top or bottom half"

		# Assert that we're not including both halves (because it's too 'zoomed out')
		assert (not (topHalf and bottomHalf)), "ERR: The image includes the whole maze"

		# Dimensions
		width = 28
		height = 16 if topHalf else 15

		# Put the four corner centroids in an array
		four_corners = np.array(centroids[1:5]).astype('float32')

		# Create an array describing the final locations of those points
		result = 100 * np.array([
			[0, 0],
			[width, 0],
			[0, height],
			[width, height]
		]).astype('float32')

		# Calculate the perspective matriix
		matrix = cv2.getPerspectiveTransform(four_corners, result)

		# Warp due to the perspective change
		if warp:
			warped = cv2.warpPerspective(img, matrix, (width * 100, height * 100))
			plt.imshow(warped, cmap='gray')                                          # type: ignore

		# Calculate the inverse perspective matrix
		inverse = np.linalg.inv(matrix)                                              # type: ignore

		# Offsets
		offset = 0 if topHalf else 16

		# Show the 'dots' on the maze
		if annotate:
			plt.imshow(img, cmap='gray')                                             # type: ignore
			for transformed_row in range(0, height):
				for transformed_col in range(0, width):
					vector = inverse @ np.array([                                    # type: ignore
						transformed_col * 100 + 50, transformed_row * 100 + 50, 1
					])
					if self.wallAt(transformed_row + offset, transformed_col):
						plt.plot([vector[0]/vector[2]], [vector[1]/vector[2]], "m.") # type: ignore
					else:
						plt.plot([vector[0]/vector[2]], [vector[1]/vector[2]], "c.") # type: ignore

		# Figure out where Pacman is
		vector = matrix @ np.array([centroids[0][0], centroids[0][1], 1])

		# Figure out the transformed centroid of Pacman
		pacman_transformed_rowf = vector[1]/vector[2]/100.0 - 0.5
		pacman_transformed_colf = vector[0]/vector[2]/100.0 - 0.5

		# Round to the nearest transformed row and column
		pacman_transformed_rowr = round(pacman_transformed_rowf)
		pacman_transformed_colr = round(pacman_transformed_colf)
		print(pacman_transformed_rowr + offset, pacman_transformed_colr, end=' -> ')

		# Loop over a 3x3 square focused on the spot
		neighbors: list[tuple[float, tuple[int, int]]] = []
		for transformed_row in range(pacman_transformed_rowr - 1, pacman_transformed_rowr + 2):
			for transformed_col in range(pacman_transformed_colr - 1, pacman_transformed_colr + 2):
				if not self.wallAt(transformed_row + offset, transformed_col):
					distSq = (transformed_row - pacman_transformed_rowf) * \
								(transformed_row - pacman_transformed_rowf) + \
							(transformed_col - pacman_transformed_colf) * \
								(transformed_col - pacman_transformed_colf)
					neighbors.append((distSq, (transformed_row + offset, transformed_col)))

		assert len(neighbors), "ERR: Pacbot was found to be in a wall"

		pacman_transformed_row, pacman_transformed_col = min(neighbors)[1]
		print(pacman_transformed_row, pacman_transformed_col)
		if annotate:
			vector = inverse @ np.array([                                            # type: ignore
				pacman_transformed_col * 100 + 50, (pacman_transformed_row - offset) * 100 + 50, 1
			])
			plt.plot([vector[0]/vector[2]], [vector[1]/vector[2]], 'y*')                     # type: ignore

		return pacman_transformed_row, pacman_transformed_col

