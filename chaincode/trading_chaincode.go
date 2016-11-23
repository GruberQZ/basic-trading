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
	// if function == "init" {
	//
	// }
	return nil, errors.New("Received unknown function invocation: " + function)
}
