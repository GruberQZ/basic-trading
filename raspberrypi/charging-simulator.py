import RPi.GPIO as GPIO
import time
import requests
import json
import sys

# Define destination URL
url = 'https://f49d9a7b5a8b4f62842153fc4019cb38-vp0.us.blockchain.ibm.com:5001/chaincode'

# Build the request object to POST
reqObj = {}
reqObj["jsonrpc"] = "2.0"
reqObj["method"] = "query"
reqObj["params"] = {}
reqObj["params"]["type"] = 1
reqObj["params"]["chaincodeID"] = {}
reqObj["params"]["chaincodeID"]["name"] = "53bef9940340843528122bc4270f75c0c7dd29d3791abaf91127c2404ddaa04843e694328fd541221199787b9086a8f098017206fed17620186e15bececd0b14"
reqObj["params"]["ctorMsg"] = {}
reqObj["params"]["ctorMsg"]["function"] = "get_charging_trade"
reqObj["params"]["ctorMsg"]["args"] = []
reqObj["params"]["secureContext"] = "user_type1_1"
reqObj["id"] = 0

# Build the complete charging object
completeObj = {}
completeObj["jsonrpc"] = "2.0"
completeObj["method"] = "invoke"
completeObj["params"] = {}
completeObj["params"]["type"] = 1
completeObj["params"]["chaincodeID"] = {}
completeObj["params"]["chaincodeID"]["name"] = "53bef9940340843528122bc4270f75c0c7dd29d3791abaf91127c2404ddaa04843e694328fd541221199787b9086a8f098017206fed17620186e15bececd0b14"
completeObj["params"]["ctorMsg"] = {}
completeObj["params"]["ctorMsg"]["function"] = "complete_charging_trade"
completeObj["params"]["ctorMsg"]["args"] = []
completeObj["params"]["secureContext"] = "user_type1_1"
completeObj["id"] = 0

# Encode object to json string
json_data = json.dumps(reqObj)

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
all_leds = [red,yellow1,yellow2,yellow3,yellow4,green]
done_charging_leds = [yellow1,yellow2,yellow3,yellow4,green]
GPIO.setup(all_leds, GPIO.OUT)

def chargeTheCar(seconds):
    # Charging sequence
    GPIO.output(red, GPIO.LOW)
    GPIO.output(yellow1, GPIO.HIGH)
    time.sleep(seconds/5)
    GPIO.output(yellow2, GPIO.HIGH)
    time.sleep(seconds/5)
    GPIO.output(yellow3, GPIO.HIGH)
    time.sleep(seconds/5)
    GPIO.output(yellow4, GPIO.HIGH)
    time.sleep(seconds/5)
    GPIO.output(green, GPIO.HIGH)
    time.sleep(seconds/5)
    # Complete Sequence
    for i in range(4):
        # Turn off everything for .5 second
        GPIO.output(all_leds, GPIO.LOW)
        time.sleep(.25)
        # Turn on everything for .5 seconds
        GPIO.output(done_charging_leds, GPIO.HIGH)
        time.sleep(.25)
    # Turn off everything
    GPIO.output(all_leds, GPIO.LOW)

while True:
    # Turn on red "standby" LED
    GPIO.output(red, GPIO.HIGH)
    time.sleep(8)
    # Request data from Bluemix
    r = requests.post(url, data=json_data)
    # Turn response into an object
    res = json.loads(r.text)
    # Get result attribute
    result = res.get('result')
    # Get the message if it exists
    try:
        message = result.get('message')
        # Turn the message into an object
        messageObj = json.loads(message)
        print(messageObj)
        print("Amount to charge: " + str(messageObj.get('amount')))
        print("Timestamp: " + str(messageObj.get('timestamp')))
        print("ID: " + str(messageObj.get('id')))
        # Complete charging
        chargeTheCar(int(messageObj.get('amount'))/10)
        # Update object params
        completeObj["params"]["ctorMsg"]["args"] = [str(messageObj.get('timestamp')),str(messageObj.get('id'))]
        # Send update to blockchain
        postData = json.dumps(completeObj)
        r = requests.post(url, data=postData)
        print(r.text)
    except:
        print("Message not found in return object")

# Clean up GPIO on exit
GPIO.cleanup()
