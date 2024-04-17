This folder contains *sample* Python code to function as a computer vision (CV) client for the Pacbot competition.

**WARNING: This CV client is experimental and currently doesn't work (you might notice it generates random numbers as coordinates). However, there's a tutorial that you can refer to, which I created to demonstrate how it WILL work. Once we confirm it works as expected on our (Princeton's) maze this weekend, I'll push the final changes, but you're welcome to implement your own version in the meantime.**

To run the sample bot client, simply run `python cvClient.py`. It may be necessary to also run `pip install -r requirements.txt`, to get important libraries installed the first time.

Other useful files:
* `cameraModule.py`: a sample camera (localization and mapping, aka SLAM) module with an asynchronous loop
* `connectionState.py`: an object which schedules localization messages to be sent to the game server
* `walls.py`: a binary representation of the maze walls (identical to `initWalls` in the server code)