package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// TODO: ^ Add "strings" and "time" back to import

// Needed for function pointer (t *SimpleChaincode)
// Leave empty
type SimpleChaincode struct {
}

var energyIdStr = "_energyindex"  // name for the key/value that will store a list of all known energy
var openTradesStr = "_opentrades" // name for the key/value that will store a list of open trades

type Energy struct {
	Owner  string `json:"owner"`  // Person who owns the energy
	Amount int    `json:"amount"` // Amount of energy
	Price  int    `json:"price"`  // Selling price
	Id     string `json:"id"`     // Unique Identifier
}

type AnOpenTrade struct {
	Owner     string `json:"owner"`     // Owner of the energy that initiates the trade
	Timestamp int64  `json:"timestamp"` // UTC Timestamp of when the offer was created
	Want      int    `json:"want"`      // Minimum amount of energy desired
	Willing   int    `json:"willing"`   // Maximum price willing to spend
}

type AllTrades struct {
	OpenTrades []AnOpenTrade `json:"open_trades"`
}

// Main function - Runs on start
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init - reset the state of the chaincode
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	// Check the number of args passed in
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1: Initial value")
	}

	// Get initial value
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	// Use test var ece because reasons
	err = stub.PutState("ece", []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	// Use JSON.Marshal to get a JSON encoding of an empty string array
	// Use that encoding to clear out the list of energy assets
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)
	err = stub.PutState(energyIdStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	// Clear the trade struct by making a new AllTrades struct (empty by default)
	// and assigning it to openTradesStr
	var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	// Successful init return
	return nil, nil
}

// Run function - entry point for invocations
// Probably unnecessary, but it can't hurt to have
// Just pass arguments along to Invoke()
func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Print debug message
	fmt.Println("Run() is running: " + function)
	// Pass arguments to Invoke()
	return t.Invoke(stub, function, args)
}

// Invoke function - entry point for invocations
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Print debug message
	fmt.Println("Invoke() is running: " + function)

	// Handle all of the different function calls
	// Initialize the state of the chaincode (RESET)
	if function == "init" {
		// Call the Init() function
		return t.Init(stub, "init", args)
	} else if function == "delete" { // Delete an entity from
		// Call the Delete() function
		res, err := t.Delete(stub, args)
		// Ensure all open trades are still valid after deletion
		cleanTrades(stub)
		// Return result from Delete()
		return res, err
	} else if function == "write" { // Write a value to the chaincode state
		// Pass arguments along to the Write() function
		return t.Write(stub, args)
	} else if function == "init_energy" { // Create a new energy block
		// Pass arguments along to the init_energy function
		return t.init_energy(stub, args)
	} else if function == "set_owner" { // Transfer ownership of energy block
		// Call the set_owner() function
		res, err := t.set_owner(stub, args)
		// Make sure open trades are still valid after ownership changes
		cleanTrades(stub)
		// Return result from set_owner()
		return res, err
	} else if function == "open_trade" { // Create a new trade order
		// Pass arguments along to the open_trade function
		return t.open_trade(stub, args)
	} else if function == "perform_trade" { // Fulfill an open trade order
		// Pass arguments along to the perform_trade function
		res, err := t.perform_trade(stub, args)
		// Make sure open trades are still valid after trade resolves
		cleanTrades(stub)
		// Return result from perform_trade()
		return res, err
	} else if function == "remove_trade" { // Cancel an open trade order
		// Pass arguments along to the remove_trade function
		return t.remove_trade(stub, args)
	}

	// Print error message if function not found
	fmt.Println("Invoke() did not find function: " + function)

	// Return error
	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query function
// Handles all query type functions
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Debug message
	fmt.Println("Query() is running: " + function)

	// Handle the different types of query functions
	if function == "read" {
		return t.read(stub, args)
	}

	// Print message if query function not found
	fmt.Println("Query() did not find function: " + function)

	// Return on error
	return nil, errors.New("Received unknown query function: " + function)
}

// Read function
// Read a variable from chaincode state
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	// Check to make sure number of arguments is correct
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1: name of variable to query")
	}

	// Get the variable from the chaincode state
	name = args[0]
	valAsBytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	// Successful return
	return valAsBytes, nil
}

// Delete function
// Remove a key/value pair from the chaincode state
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Check number of arguments passed in
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1: key/value pair to delete")
	}

	// Remove key from the chaincode state
	name := args[0]
	err := stub.DelState(name)
	if err != nil {
		return nil, errors.New("Failed to delete state: " + name)
	}

	// Get the energy index
	energyAsBytes, err := stub.GetState(energyIdStr)
	if err != nil {
		return nil, errors.New("Failed to get energy index")
	}

	// Turn the energy index into a string array
	var energyIndex []string
	json.Unmarshal(energyAsBytes, &energyIndex)

	// Iterate through the energy index looking for the energy asset
	for i, val := range energyIndex {
		// Debug message
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		// Determine if this is the correct energy asset
		if val == name {
			fmt.Println("Found energy asset to remove")
			// Reconstruct the energyIndex to include everything except asset to be deleted
			energyIndex = append(energyIndex[:i], energyIndex[i+1:]...)
			// Debug: Print out all of the assets in the energyIndex
			fmt.Println("New state of energy assets:")
			for x := range energyIndex {
				fmt.Println(string(x) + " - " + energyIndex[x])
			}
			// Found asset to remove, break
			break
		}
	}

	// Turn energyIndex back into JSON for chaincode state
	jsonAsBytes, _ := json.Marshal(energyIndex)
	err = stub.PutState(energyIdStr, jsonAsBytes)
	// Successful exit
	return nil, nil
}

// Write function
// Write variables into the chaincode state
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string
	var err error

	// Debug message
	fmt.Println("Running Write()")

	// Rename for clarity
	name = args[0]
	value = args[1]
	// Write variable into the chaincode state
	err = stub.PutState(name, []byte(value))
	if err != nil {
		return nil, err
	}

	// Successful exit
	return nil, nil
}
