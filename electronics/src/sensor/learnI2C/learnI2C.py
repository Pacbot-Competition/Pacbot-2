# RPI Master for Teensy slave
# connect RPI to Teensy via I2C
# reference: dronebotworkshop.com/i2c-arduino-raspberry-pi/#I2C_Test_Code

from smbus import SMBus

addr = 0x8 # bus address
bus = SMBus(1) # indicates /dev/ic2-1

numb = 1

print ("Enter 1 for ON or 0 for OFF")
while numb == 1:

	ledstate = input(">>>>   ")

	if ledstate == "1":
		bus.write_byte(addr, 0x1) # switch it on
	elif ledstate == "0":
		bus.write_byte(addr, 0x0) # switch it on
	else:
		numb = 0