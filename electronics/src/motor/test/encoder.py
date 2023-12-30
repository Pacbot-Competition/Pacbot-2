from gpiozero import DigitalInputDevice
from time import sleep

RightOutA = 7   # GPIO7
RightOutB = 8   # GPIO8
LeftOutA = 24   # GPIO24
LeftOutB = 23   # GPIO23


RightEncoderOutA = DigitalInputDevice(pin=RightOutA)
RightEncoderOutB = DigitalInputDevice(pin=RightOutB)
LeftEncoderOutA = DigitalInputDevice(pin=LeftOutA)
LeftEncoderOutB = DigitalInputDevice(pin=LeftOutB)

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


while True:
    print("Right Encoder OUTA={}, OUTB={}".format(RightEncoderOutA.value, RightEncoderOutB.value))
    print("Left Encoder  OUTA={}, OUTB={}".format(LeftEncoderOutA.value, LeftEncoderOutB.value))