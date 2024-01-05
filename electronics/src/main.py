# run low-level code 
# including motor, magnetic encoder, motor PID control, and IR sensor
from motor import PID
from sensor import IRsensor

# movement commands definition
BRAKE       = 0
FORWARD     = 1
BACKWORD    = 2
RIGHT       = 3
LEFT        = 4

SAMPLETIME = 50 # sampling time to get distance info

def setup() -> None:
    # initialize robot
    # initialize PID
    # initialize IR sensor
    pass

def loop(command: int) -> None:
    # run motor PID and get sensor data
    pass

# interface with other modules
def getCommand() -> int:
    # get command from other thread to control the robot
    command = BRAKE

    return command

def getDistance() -> list:
    # get distance from 4 IR sensors and send them to other modules
    results = []

    return results

# run the robot
setup()
while True:
    command = getCommand()
    getDistance(loop(command), SAMPLETIME)
