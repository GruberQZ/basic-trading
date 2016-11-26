# basic-trading
Testing Blockchain

Based on IBM Blockchain Marbles Tutorial

## Usage
This section breaks chaincode operations into sections based on their type and their usage. To use these commands, edit the "ctorMsg" property of the JSON object that is sent to /chaincode. Arguments to functions are always passed in as a string array.

### Query
The "method" property in the JSON object that is sent to /chaincode for operations in this section should be set to "query".

#### Read a variable from the chaincode state
Function name: "read"
Arguments: 1
1) The name of the variable to be read from the chaincode state

```javascript
Example:
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
