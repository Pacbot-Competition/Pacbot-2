from gpiozero import Robot, DigitalInputDevice
from time import sleep

# use this code to get the number of tick's per sample
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

SAMPLETIME = 1

# motor pins
RightEnable = 18   # BEN, GPIO18
RightPhase = 19  # BPH, GPIO19
LeftEnable = 12    # AEN, GPIO12
LeftPhase = 6    # APH, GPIO6

# magnetic encoder pins
RightOutA = 7   # GPIO7
RightOutB = 8   # GPIO8
LeftOutA = 24   # GPIO24
LeftOutB = 23   # GPIO23

r = Robot(left=(LeftPhase, LeftEnable), right=(RightPhase, RightEnable))
e1 = QuadratureEncoder(RightOutA, RightOutB)
e2 = QuadratureEncoder(LeftOutA, LeftOutB)

r.value = (1,1)

while True:
    print("e1 {} e2 {}".format(e1.value, e2.value))
    sleep(SAMPLETIME)