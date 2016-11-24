package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Needed for function pointer (t *SimpleChaincode)
// Leave empty
type SimpleChaincode struct {
}

var energyIndexStr = "_energyindex" // name for the key/value that will store a list of all known energy
var openTradesStr = "_opentrades"   // name for the key/value that will store a list of open trades

type Energy struct {
	Id     string `json:"id"`     // Unique Identifier
	Amount int    `json:"amount"` // Amount of energy
	Price  int    `json:"price"`  // Selling price
	Owner  string `json:"owner"`  // Person who owns the energy
}

type AnOpenTrade struct {
	Owner     string `json:"owner"`     // Owner of the energy that initiates the trade
	Timestamp int64  `json:"timestamp"` // UTC Timestamp of when the offer was created
	Id        string `json:"id"`        // Id of the asset used to create the trade
	Amount    int    `json:"amount"`    // Amount of energy for trade
	Price     int    `json:"price"`     // Minimum price willing to accept
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
	err = stub.PutState(energyIndexStr, jsonAsBytes)
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
	energyAsBytes, err := stub.GetState(energyIndexStr)
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
	err = stub.PutState(energyIndexStr, jsonAsBytes)
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

// Initialize new energy asset
// Create a new energy asset and store it in the chaincode state
func (t *SimpleChaincode) init_energy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	// Arguments passed in the following order:
	// 0 --> "asset1" == Unique Identifier
	// 1 --> "50" == Amount of energy
	// 2 --> "25" == Price of energy
	// 3 --> "bob" == Owner of the energy

	// Check the number of arguments passed in
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4.")
	}

	// Check for valid input
	fmt.Println("Creating a new energy asset")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}

	// Rename and convert variables
	id := strings.ToLower(args[0])
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}
	price, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	owner := strings.ToLower(args[3])

	// Check to see if energy asset with this id already exists
	assetAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New("Failed to get energy asset id")
	}
	res := Energy{}
	json.Unmarshal(assetAsBytes, &res)
	// If there exists an energy asset with the same name, error
	if res.Id == id {
		fmt.Println("An energy asset with this id already exists: " + id)
		fmt.Println(res)
		return nil, errors.New("An energy asset with this id already exists")
	}

	// Build the JSON string representation of the new energy asset
	str := `{"id": "` + id + `", "amount": ` + strconv.Itoa(amount) + `, "price": ` + strconv.Itoa(price) + `, "owner": "` + owner + `"}`
	// Store this energy asset with the id as the key
	err = stub.PutState(id, []byte(str))

	// Add this energy asset to the energy index
	energyAsBytes, err := stub.GetState(energyIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get the energy index")
	}

	// Convert the energy index into an array of strings
	var energyIndex []string
	json.Unmarshal(energyAsBytes, &energyIndex)

	// Append the energy asset to the energy index
	energyIndex = append(energyIndex, id)
	// Debug message
	fmt.Println("Current energy index: ", energyIndex)
	// Re-encode the new energy index and write it back to the chaincode state
	jsonAsBytes, _ := json.Marshal(energyIndex)
	err = stub.PutState(energyIndexStr, jsonAsBytes)

	// Debug message & successful return
	fmt.Println("End initialize energy asset")
	return nil, nil
}

// set_owner function
// Set the owner of an energy asset
func (t *SimpleChaincode) set_owner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	// Arguments passed in the following order:
	// 0 --> "asset1" == Unique Identifier
	// 1 --> "alice" == New owner of this asset
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2: asset identifier and new owner")
	}

	id := args[0]
	owner := args[1]
	fmt.Println("Starting set owner")
	fmt.Println("Setting owner of " + id + " to " + owner)
	energyAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New("Failed to get energy asset")
	}

	// Get the asset and change the owner
	res := Energy{}
	json.Unmarshal(energyAsBytes, &res)
	res.Owner = owner

	// Rewrite the energy asset into the chaincode state
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(id, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	// Successful exit
	fmt.Println("Done setting owner")
	return nil, nil
}

// open_trade function
// Create an open trade for an energy asset you have
func (t *SimpleChaincode) open_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	fmt.Println("Starting open_trade")
	// Arguments passed in the following order:
	// 0 --> "bob" == Owner
	// 1 --> "asset1" == Unique Identifier

	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4.")
	}
	owner := args[0]
	id := args[1]

	// Verify ownership of the asset before opening up a trade
	assetAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New("Failed to get energy asset id")
	}
	res := Energy{}
	json.Unmarshal(assetAsBytes, &res)
	if res.Owner != owner {
		return nil, errors.New("Invalid trade opening: " + owner + " does not own the asset " + id)
	}

	// Verify that the asset is not current part of an outstanding trade
	// Get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)

	// Search through open trades looking for asset
	for i := range trades.OpenTrades {
		if trades.OpenTrades[i].Id == id {
			return nil, errors.New("Invalid trade opening: Asset for trade cannot be part of existing open trade")
		}
	}

	// Ownership has been verified and asset is not part of an outstanding offer
	// Create a new trade offer
	open := AnOpenTrade{}
	open.Owner = owner
	open.Timestamp = time.Now().Unix() // [Use timestamp as unique identifier for trades]
	open.Amount = res.Amount
	open.Price = res.Price
	open.Id = id
	// Set a variable in the chaincode for debug
	jsonAsBytes, _ := json.Marshal(open)
	err = stub.PutState("_lastOpenedTrade", jsonAsBytes)

	// Append the new trade to the list of open trades
	trades.OpenTrades = append(trades.OpenTrades, open)
	fmt.Println("Appended new trade opening")
	jsonAsBytes, _ = json.Marshal(trades)
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("End open trade")
	return nil, nil
}

// Perform trade function
// Close an open trade and move ownership to the buyer
func (t *SimpleChaincode) perform_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	// Arguments are passed in the following order:
	// 0 --> "asset1" == Unique Id of the energy asset to trade
	// 1 --> "alice" == New owner
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}
	id := args[0]
	newOwner := args[1]

	fmt.Println("Start perform trade")

	// Get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)

	for i := range trades.OpenTrades {
		fmt.Println("Looking at " + trades.OpenTrades[i].Id + " for " + id)
		if trades.OpenTrades[i].Id == id {
			fmt.Println("Found the trade to perform")

			// Get the asset that will be traded
			assetAsBytes, err := stub.GetState(id)
			if err != nil {
				return nil, errors.New("Failed to get the asset")
			}
			asset := Energy{}
			json.Unmarshal(assetAsBytes, &asset)

			// Verify that the new owner is not the current owner
			if asset.Owner == newOwner {
				return nil, errors.New("New asset owner cannot be the current owner")
			}

			// Change the owner of the asset
			t.set_owner(stub, []string{id, newOwner})

			// Remove the trade from the list of open trades
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)
			jsonAsBytes, _ := json.Marshal(trades)
			err = stub.PutState(openTradesStr, jsonAsBytes)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	fmt.Println("End perform trade")
	return nil, nil
}

// Remove Open trade
// Close an open trade with no change of ownership taking place
func (t *SimpleChaincode) remove_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	// Only argument needed is the unique ID of the asset that should no longer be eligible for trade
	if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1.")
	}
	id := args[0]

	fmt.Println("Begin remove trade")

	// get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get open trades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)

	// Look for the trade in the list of open trades
	for i := range trades.OpenTrades {
		if trades.OpenTrades[i].Id == id {
			fmt.Println("Found trade to remove")
			// Remove this trade from the list
			trades.OpenTrades = append(trades.OpenTrades[:i], trades.OpenTrades[i+1:]...)
			jsonAsBytes, _ := json.Marshal(trades)
			// Rewrite the open orders to the chaincode state
			err = stub.PutState(openTradesStr, jsonAsBytes)
			if err != nil {
				return nil, err
			}
			// Done removing, break
			break
		}
	}

	// Successful return
	fmt.Println("End remove trade")
	return nil, nil
}

// Clean up open trades
// Make sure all open trades are still possible/valid, and remove those that are not
