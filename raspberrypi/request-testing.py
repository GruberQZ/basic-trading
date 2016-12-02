import requests
import json

# Define destination URL
url = 'https://api.github.com/events'
# url = 'https://b5b1f30cd80c4041972890286eb7e5df-vp0.us.blockchain.ibm.com:5001/chaincode'

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

# Encode object to json string
json_data = json.dumps(reqObj)
#print(json_data)

r = requests.get(url)