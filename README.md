# basic-trading  
Testing Blockchain on Bluemix

Based on IBM Blockchain Marbles Tutorial 

## Warning
The params.chaincodeID.name property that is used in the examples below is likely out of date. Contact the author for the updated chaincodeID.
# Usage  
This section breaks chaincode operations into sections based on their type and their usage. To use these commands, edit the "ctorMsg" property of the JSON object that is sent to /chaincode. Arguments to functions are always passed in as a string array.
## Query  
The "method" property in the JSON object that is sent to /chaincode for operations in this section should be set to "query".
### Read a variable from the chaincode state  
Function name: "read"  
Arguments: 1  
1) The name of the variable to be read from the chaincode state  
  
Example: Read the variable "ece" from the chaincode state  
```javascript
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "read",
      "args": [
        "ece"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### See all of the available Query functions  
Function name: "query_functions"  
Arguments: 0  

Example:  
```javascript
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "query_functions",
      "args": []
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### See all of the available Invoke functions
Function name: "invoke_functions"  
Arguments: 0  

Example:  
```javascript
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "invoke_functions",
      "args": []
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### See all of the open trade orders
Function name: "open_trades"  
Arguments: 0  

Example:  
```javascript
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "open_trades",
      "args": []
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
## Invoke  
The "method" property in the JSON object that is sent to /chaincode for operations in this section should be set to "query".  
### Write a variable to the chaincode state  
Function name: "write"  
Arguments: 2  
1) Name of the variable  
2) Value of the variable  

Example: Write the value "485" to the variable "ece" to the chaincode state  
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "write",
      "args": [
        "ece",
        "485"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### Delete a variable from the chaincode state
Function name: "delete"  
Arguments: 1  
1) Name of the variable to delete

Example: Delete the variable "ece" from the chaincode state
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "delete",
      "args": [
        "ece"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### Create a new Energy asset  
Function name: "init_energy"  
Arguments: 4  
1) Unique identifier  
2) Amount of energy in this asset  
3) Price of energy in this asset  
4) Owner of the energy  

Example: Create a new energy asset called "asset1" that is owned by bob with energy amount 50 and price 25  
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "init_energy",
      "args": [
        "asset1",
        "50",
        "25",
        "bob"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### Set the owner of an energy asset
Function name: "set_owner"  
Arguments: 2  
1) Unique Indentifier of an energy asset  
2) New owner of that energy asset  

Example: Set the owner of asset1 to alice  
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "set_owner",
      "args": [
        "asset1",
        "alice"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```
### Open up a new trade order  
Function name: "open_trade"  
Arguments: 2  
1) Creator of the energy asset  
2) Energy asset to be traded  

Example: Bob creates an open trade order for his asset called "asset1"  
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "open_trade",
      "args": [
        "bob",
        "asset1"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```

### Fulfill an open trade order
Function name: "perform_trade"  
Arguments: 2  
1) Energy asset that will be bought  
2) New owner of the energy asset (purchaser)    

Example: alice purchases the asset "asset1" that bob listed previously  
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "perform_trade",
      "args": [
        "asset1",
        "alice"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```

### Remove an open trade order
Function name: "remove_trade"  
Arguments: 2  
1) Creator of the open trade agreement  
2) Unique Identifier of energy asset in open trade order  

Example: Remove the bob's open trade order for asset1  
```javascript
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "2553575989126bf89371ff4a63c40221f72d0f141ffdfde3ba196fde5df53621f1295ce19dbcc92d68dc5c67235e056b1eb52e9bdde9e03c8e799f22f8439910"
    },
    "ctorMsg": {
      "function": "remove_trade",
      "args": [
        "bob",
        "asset1"
      ]
    },
    "secureContext": "user_type1_1"
  },
  "id": 0
}
```

# The Energy Asset
```javascript
{
  "id": "asset1"  // Unique Identifier
  "amount" "50"   // Amount of energy
  "price": "25"   // Selling price
  "owner": "bob"  // Person who owns the energy
}
```
# The Open Trade order
```javascript
{
  "owner": "bob"           // Owner of the energy that initiates the trade
  "timestamp": <timestamp> // UTC Timestamp of when the offer was created
  "id": "asset1"           // Id of the asset used to create the trade
  "amount": "50"           // Amount of energy for trade
  "price": "25"            // Minimum price willing to accept
}
```
