from gpiozero import DigitalInputDevice, PhaseEnableRobot
from time import sleep

class QuadratureEncoder(object):
    """
    A simple quadrature encoder class

    Note - this class does not determine direction
    """
    def __init__(self, pin_a, pin_b):
        self._value = 0

        encoder_a = DigitalInputDevice(pin_a)
        encoder_a.when_activated = self._increment
        encoder_a.when_deactivated = self._increment

        encoder_b = DigitalInputDevice(pin_b)
        encoder_b.when_activated = self._increment
        encoder_b.when_deactivated = self._increment
        
    def reset(self):
        self._value = 0

    def _increment(self):
        self._value += 1

    @property
    def value(self):
        return self._value

def clamp(value):
    return max(min(1, value), 0)

# PID parameters
SAMPLETIME = 0.5
TARGET = 20         # about 75% of the encoder's tick per sample
KP = 0.02           # start at 1 divided by the encoder's tick per sample
KD = 0.01           # start at 0.5*KP
KI = 0.005          # start at 0.5*KD

# motor pins
RightEnable = 18   # BEN, GPIO18
RightPhase = 19  # BPH, GPIO19
LeftEnable = 12    # AEN, GPIO12
LeftPhase = 6    # APH, GPIO6

robot = PhaseEnableRobot(left=(LeftPhase, LeftEnable), right=(RightPhase, RightEnable))

# magnetic encoder pins
RightOutA = 7   # GPIO7
RightOutB = 8   # GPIO8
LeftOutA = 24   # GPIO24
LeftOutB = 23   # GPIO23

rightEncoder = QuadratureEncoder(RightOutA, RightOutB)
leftEncoder = QuadratureEncoder(LeftOutA, LeftOutB)

rightMotorSpeed = 1
leftMotorSpeed = 1
robot.value = (rightMotorSpeed, leftMotorSpeed)

rightEncoderPrevError = 0
leftEncoderPrevError = 0

rightEncoderSumError = 0
leftEncoderSumError = 0

while True:
    # make left and right motor run at the same speed
    # i.e. make robot go in straight line
    rightEncoderError = TARGET - rightEncoder.value
    leftEncoderError = TARGET - leftEncoder.value

    rightMotorSpeed += (rightEncoderError * KP) + (rightEncoderPrevError * KD) + (rightEncoderSumError * KI)
    leftMotorSpeed += (leftEncoderError * KP)  + (rightEncoderPrevError * KD) + (leftEncoderSumError * KI)

    rightMotorSpeed = clamp(rightMotorSpeed)
    leftMotorSpeed = clamp(leftMotorSpeed)
    robot.value = (rightMotorSpeed, leftMotorSpeed)

    print("right encoder {} left encoder {}".format(rightEncoder.value, leftEncoder.value))
    print("right motor {} left motor {}".format(rightMotorSpeed, leftMotorSpeed))

    rightEncoder.reset()
    leftEncoder.reset()

    sleep(SAMPLETIME)

    rightEncoderPrevError = rightEncoderError
    leftEncoderPrevError = leftEncoderError

    rightEncoderSumError += rightEncoderError
    leftEncoderSumError += leftEncoderError