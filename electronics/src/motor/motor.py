from gpiozero import Motor
from time import sleep

forward = 4
backward = 14

motor = Motor(orward=forward, backward=backward)

while True:
    motor.forward()