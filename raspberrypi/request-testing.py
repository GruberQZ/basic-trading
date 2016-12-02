import requests
import json
import sys
import time

# Define destination URL
#url = 'https://api.github.com/events'
url = 'https://b5b1f30cd80c4041972890286eb7e5df-vp0.us.blockchain.ibm.com:5001/chaincode'

# Build the request object to POST
reqObj = {}
reqObj["jsonrpc"] = "2.0"
reqObj["method"] = "query"
reqObj["params"] = {}
reqObj["params"]["type"] = 1
reqObj["params"]["chaincodeID"] = {}
reqObj["params"]["chaincodeID"]["name"] = "12decb80905604f0f7ea238c51e8cd4c7496640865ba48b8798d8c35eb9be86b206c7119c925aa115513941a5c10dc9af34e6530ced4e673db2376553aa156ce"
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
completeObj["params"]["chaincodeID"]["name"] = "12decb80905604f0f7ea238c51e8cd4c7496640865ba48b8798d8c35eb9be86b206c7119c925aa115513941a5c10dc9af34e6530ced4e673db2376553aa156ce"
completeObj["params"]["ctorMsg"] = {}
completeObj["params"]["ctorMsg"]["function"] = "complete_charging_trade"
completeObj["params"]["ctorMsg"]["args"] = []
completeObj["params"]["secureContext"] = "user_type1_1"
completeObj["id"] = 0

# Encode object to json string
json_data = json.dumps(reqObj)
#print(json_data)

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
    time.sleep(4)
    # Complete charging
    # Update object params
    completeObj["params"]["ctorMsg"]["args"] = [str(messageObj.get('timestamp')),str(messageObj.get('id'))]
    # Send update to blockchain
    postData = json.dumps(completeObj)
    r = requests.post(url, data=postData)
    print(r.text)
except:
    print("Message not found in return object")
    sys.exit()

print(result)











