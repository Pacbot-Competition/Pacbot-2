This folder contains *sample* Python code to function as a computer vision (CV) client for the Pacbot competition.

To run the sample bot client, simply run `python cvClient.py`. It may be necessary to also run `pip install -r requirements.txt`, to get important libraries installed the first time.

Other useful files:
* `decisionModule.py`: a sample decision module (policy) with an asynchronous loop and game state locking capabilities
* `connectionState.py`: an object which schedules localization messages to be sent to the game server
* `walls.py`: a binary representation of the maze walls (identical to `initWalls` in the server code)