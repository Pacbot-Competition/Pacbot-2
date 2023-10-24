This folder contains *sample* Python code to function as a high-level client for the Pacbot competition.

Teams are encouraged (and expected) to modify client code to fit their navigation algorithms and robot communication protocols.

To run the sample bot client, simply run `python pacbotClient.py`. It may be necessary to also run
`pip install -r requirements.txt`, to get important libraries installed the first time.

Other useful files:
* `decisionModule.py`: a sample decision module (policy) with an asynchronous loop and game state locking capabilities
* `gameState.py`: a game state object which parses serialized data and offers simple methods to interact with and predict the game state
* `walls.py`: a binary representation of the maze walls (identical to `initWalls` in the server code)