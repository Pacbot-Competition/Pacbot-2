from gpiozero import PhaseEnableRobot, DigitalInputDevice
from time import sleep

# motor pins
RightEnable = 18   # BEN, GPIO18
RightPhase = 19  # BPH, GPIO19
LeftEnable = 12    # AEN, GPIO12
LeftPhase = 6    # APH, GPIO6
speed = 0.5

robot = PhaseEnableRobot(left=(LeftPhase, LeftEnable), right=(RightPhase, RightEnable))

# magnetic encoder pins
RightOutA = 7   # GPIO7
RightOutB = 8   # GPIO8
LeftOutA = 24   # GPIO24
LeftOutB = 23   # GPIO23

# Direction
# right motor forward: 
# OUTA OUTB
#  0    1
#  0    0
#  1    0
#  1    1
# right motor backward: 
# OUTA OUTB
#  0    1
#  1    1
#  1    0
#  0    0
# left motor forward: 
# OUTA OUTB
#  0    1
#  1    1
#  1    0
#  0    0
# left motor backward: 
# OUTA OUTB
#  0    1
#  0    0
#  1    0
#  1    1

RightEncoderOutA = DigitalInputDevice(pin=RightOutA)
RightEncoderOutB = DigitalInputDevice(pin=RightOutB)
LeftEncoderOutA = DigitalInputDevice(pin=LeftOutA)
LeftEncoderOutB = DigitalInputDevice(pin=LeftOutB)

while True:
    robot.backward(speed=speed)
    print("Left Encoder  OUTA={}, OUTB={}".format(LeftEncoderOutA.value, LeftEncoderOutB.value))