# basic-trading  
Testing Blockchain  

Based on IBM Blockchain Marbles Tutorial  
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

### Open up a new trade order

### Fulfill an open trade order

### Remove an open trade order

# The Energy Asset
```javascript
{
  "id": "asset1"  // Unique Identifier
  "amount" "50"   // Amount of energy
  "price": "25"   // Selling price
  "owner": "bob"  // Person who owns the energy
}
```
