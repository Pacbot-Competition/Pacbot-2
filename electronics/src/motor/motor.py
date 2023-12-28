from gpiozero import PhaseEnableMotor
from time import sleep

RightEnable = 18   # BEN, GPIO18
RightPhase = 19  # BPH, GPIO19
LeftEnable = 12    # AEN, GPIO12
LeftPhase = 6    # APH, GPIO6
RightSpeed = 1
LeftSpeed = 1

RightMotor = PhaseEnableMotor(phase=RightPhase, enable=RightEnable)
LeftMotor = PhaseEnableMotor(phase=LeftPhase, enable=LeftEnable)

# RightMotor.forward(speed=RightSpeed)
# LeftMotor.forward(speed=LeftSpeed)

while True:
    RightMotor.forward(speed=RightSpeed)
    LeftMotor.forward(speed=LeftSpeed)
    # RightMotor.backward()
    # LeftMotor.backward()
    # RightMotor.stop()
    # LeftMotor.stop()
    # sleep(1)
    # RightMotor.forward()
    # LeftMotor.forward()
    # sleep(1)