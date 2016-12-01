import RPi.GPIO as GPIO
import time

# Turn off warnings because this is the only script altering GPIO
GPIO.setwarnings(False)
# Set pin numbering system
GPIO.setmode(GPIO.BCM)

# Define channels for LEDs
red = 17
yellow1 = 27
yellow2 = 22
yellow3 = 16
yellow4 = 20
green = 21

# Set up inputs
chan_list = [red,yellow1,yellow2,yellow3,yellow4,green]
GPIO.setup(chan_list, GPIO.OUT)

# Turn on all LEDs, wait a second, turn them all off again
GPIO.output(chan_list, GPIO.HIGH)
time.sleep(1)
GPIO.output(chan_list, GPIO.LOW)

# Clean up GPIO on exit
GPIO.cleanup()
