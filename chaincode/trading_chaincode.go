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

var energyIndexStr = "_energyindex"      // name for the key/value that will store a list of all known energy
var openTradesStr = "_opentrades"        // name for the key/value that will store a list of open trades
var chargingTradesStr = "_waitingtrades" // name for the key/value that will store a list of trades waiting for charger

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

type ChargingTrades struct {
	WaitingTrades []AnOpenTrade `json:"waiting_trades"`
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
	var retStr string

	// Check the number of args passed in
	if len(args) != 1 {
		retStr = "Incorrect number of arguments. Expecting 1: Initial value"
		return []byte(retStr), errors.New("Incorrect number of arguments. Expecting 1: Initial value")
	}

	// Get initial value
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		retStr = "Incorrect number of arguments. Expecting 1: Initial value"
		return []byte(retStr), errors.New("Expecting integer value for asset holding")
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

	// Clear the trade struct by making a new ChargingTrades struct (empty by default)
	// and assigning it to chargingTradesStr
	var trades2 ChargingTrades
	jsonAsBytes, _ = json.Marshal(trades2)
	err = stub.PutState(chargingTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	// Successful init return
	retStr = "Chaincode state initialized successfully."
	return []byte(retStr), nil
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
	} else if function == "set_price" {
		// Call the set_owner() function
		res, err := t.set_price(stub, args)
		// Make sure open trades are still valid after price changes
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
	} else if function == "complete_charging_trade" { // Get the next trade to be charged
		// Pass arguments along to the get_charging_trade function
		return t.complete_charging_trade(stub, args)
	}

	// Print error message if function not found
	fmt.Println("Invoke() did not find function: " + function)

	// Return error
	return []byte("Invoke() did not find function: " + function), errors.New("Received unknown function invocation: " + function)
}

// Query function
// Handles all query type functions
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Debug message
	fmt.Println("Query() is running: " + function)

	// Handle the different types of query functions
	if function == "read" {
		return t.read(stub, args)
	} else if function == "query_functions" {
		return t.query_functions(stub)
	} else if function == "invoke_functions" {
		return t.invoke_functions(stub)
	} else if function == "open_trades" {
		return t.open_trades(stub)
	} else if function == "view_my_assets" {
		return t.view_my_assets(stub, args)
	} else if function == "get_charging_trade" {
		return t.get_charging_trade(stub, args)
	}

	// Print message if query function not found
	fmt.Println("Query() did not find function: " + function)

	// Return an error
	return []byte("Query() did not find function: " + function), errors.New("Received unknown query function: " + function)
}

// Read function
// Read a variable from chaincode state
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error
	var retStr string

	// Check to make sure number of arguments is correct
	if len(args) != 1 {
		retStr = "Incorrect number of arguments. Expecting 1: name of variable to query"
		return []byte(retStr), errors.New(retStr)
	}

	// Get the variable from the chaincode state
	name = args[0]
	valAsBytes, err := stub.GetState(name)
	if err != nil {
		retStr = "Could not get state for variable " + name
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return []byte(retStr), errors.New(jsonResp)
	}

	// Return message if variable doesn't exist
	// Variable does not exist if byte array has length 0
	if len(valAsBytes) == 0 {
		return []byte("Variable \"" + name + "\" does not exist"), nil
	}

	// Successful return
	return valAsBytes, nil
}

// Query functions
// Return a list of all available query function names
func (t *SimpleChaincode) query_functions(stub shim.ChaincodeStubInterface) ([]byte, error) {
	retStr := "read, query_functions, invoke_functions, open_trades"
	return []byte(retStr), nil
}

// Invoke functions
// Return a list of all available invoke function names
func (t *SimpleChaincode) invoke_functions(stub shim.ChaincodeStubInterface) ([]byte, error) {
	retStr := "init, delete, write, init_energy, set_owner, set_price, open_trade, perform_trade, remove_trade"
	return []byte(retStr), nil
}

// Query Open Trades
// View all open trades
func (t *SimpleChaincode) open_trades(stub shim.ChaincodeStubInterface) ([]byte, error) {
	// Get the open trade array from the chaincode state
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	return tradesAsBytes, nil
}

// View My Assets
// View a list of all assets owned by an individual
func (t *SimpleChaincode) view_my_assets(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var retStr string

	// Get parameters
	if len(args) != 1 {
		retStr = "Incorrect number of arguments. Expecting 1: Owner's name"
		return []byte(retStr), errors.New(retStr)
	}
	owner := args[0]
	var ownerAssets []Energy

	// Get the energy index
	energyAsBytes, err := stub.GetState(energyIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get energy index")
	}
	// Turn the energy index into a string array
	var energyIndex []string
	json.Unmarshal(energyAsBytes, &energyIndex)

	// Iterate through the energy index
	for i, id := range energyIndex {
		// Debug message
		fmt.Println(strconv.Itoa(i) + ": looking at asset with id " + id)
		// Get the asset from the chaincode state
		assetAsBytes, err := stub.GetState(id)
		if err != nil {
			return nil, errors.New("Failed to get energy asset id")
		}
		res := Energy{}
		json.Unmarshal(assetAsBytes, &res)
		// If the asset is owned by the owner requesting the list,
		// add the asset to the ownerAssets array
		if res.Owner == owner {
			ownerAssets = append(ownerAssets, res)
		}
	}

	// Return the completed list to the requester
	// Send message if owner does not own anything
	if len(ownerAssets) == 0 {
		retStr = owner + ", it appears you do not own any assets!"
		return []byte(retStr), errors.New(retStr)
	}
	// Prepare ownerAssets for return
	jsonAsBytes, _ := json.Marshal(ownerAssets)
	return jsonAsBytes, nil
}

// Delete function
// Remove a key/value pair from the chaincode state
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var retStr string
	// Check number of arguments passed in
	if len(args) != 1 {
		retStr = "Incorrect number of arguments. Expecting 1: key/value pair to delete"
		return []byte(retStr), errors.New(retStr)
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
	retStr = "Variable/Asset [" + name + "] deleted successfully."
	return []byte(retStr), nil
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
	retStr := "Successfully set value of [" + name + "] to " + value
	return []byte(retStr), nil
}

// Initialize new energy asset
// Create a new energy asset and store it in the chaincode state
func (t *SimpleChaincode) init_energy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var retStr string

	// Arguments passed in the following order:
	// 0 --> "asset1" == Unique Identifier
	// 1 --> "50" == Amount of energy
	// 2 --> "25" == Price of energy
	// 3 --> "bob" == Owner of the energy

	// Check the number of arguments passed in
	if len(args) != 4 {
		retStr = "Incorrect number of arguments. Expecting 4."
		return []byte(retStr), errors.New(retStr)
	}

	// Check for valid input
	fmt.Println("Creating a new energy asset")
	if len(args[0]) <= 0 {
		retStr = "1st argument must be a non-empty string"
		return []byte(retStr), errors.New(retStr)
	}
	if len(args[1]) <= 0 {
		retStr = "2nd argument must be a non-empty string"
		return []byte(retStr), errors.New(retStr)
	}
	if len(args[2]) <= 0 {
		retStr = "3rd argument must be a non-empty string"
		return []byte(retStr), errors.New(retStr)
	}
	if len(args[3]) <= 0 {
		retStr = "4th argument must be a non-empty string"
		return []byte(retStr), errors.New(retStr)
	}

	// Rename and convert variables
	id := strings.ToLower(args[0])
	amount, err := strconv.Atoi(args[1])
	if err != nil {
		retStr = "2nd argument must be a numeric string"
		return []byte(retStr), errors.New(retStr)
	}
	price, err := strconv.Atoi(args[2])
	if err != nil {
		retStr = "3rd argument must be a numeric string"
		return []byte(retStr), errors.New(retStr)
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
		retStr = "An energy asset with id [" + id + "] already exists"
		return []byte(retStr), errors.New(retStr)
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
	retStr = "Asset [" + id + "] created successfully."
	return []byte(retStr), nil
}

// set_owner function
// Set the owner of an energy asset
func (t *SimpleChaincode) set_owner(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var retStr string
	var err error

	// Arguments passed in the following order:
	// 0 --> "asset1" == Unique Identifier
	// 1 --> "alice" == New owner of this asset
	if len(args) != 2 {
		retStr = "Incorrect number of arguments. Expecting 2: asset identifier and new owner"
		return []byte(retStr), errors.New(retStr)
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
	retStr = "Asset [" + id + "] is now owned by " + owner
	return []byte(retStr), nil
}

// Set price function
// Set the price attribute of an energy asset
func (t *SimpleChaincode) set_price(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var retStr string

	// Check arguments
	if len(args) != 3 {
		retStr = "Incorrect number of arguments. Expecting 3."
		return []byte(retStr), errors.New(retStr)
	}
	owner := args[0]
	id := args[1]
	newPrice, err := strconv.Atoi(args[2])
	if err != nil {
		retStr = "Price (3rd argument) must be a numeric string"
		return []byte(retStr), errors.New(retStr)
	}
	if newPrice < 0 {
		retStr = "Price (3rd argument) must not be less than 0"
		return []byte(retStr), errors.New(retStr)
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
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + id)
		// Determine if this is the correct energy asset
		if val == id {
			fmt.Println("Found energy asset to edit")
			// Get the asset from the chiancode state
			assetAsBytes, err := stub.GetState(id)
			if err != nil {
				return nil, errors.New("Failed to get energy asset")
			}
			res := Energy{}
			json.Unmarshal(assetAsBytes, &res)
			// Verify that the energy asset is owned by the person requesting the price adjustment
			if res.Owner != owner {
				retStr = "Error: " + owner + " does not own asset [" + id + "]"
				return []byte(retStr), errors.New(retStr)
			}
			// Change the price of the asset and write it back into the chaincode state
			res.Price = newPrice
			jsonAsBytes, _ := json.Marshal(res)
			err = stub.PutState(id, jsonAsBytes)
			if err != nil {
				return nil, err
			}
			// Found asset to remove, return
			retStr = "Price of " + id + " changed to " + args[2]
			return []byte(retStr), nil
		}
	}

	// Unsuccessful return
	retStr = "Price could not be set because asset [" + id + "] does not exist"
	return []byte(retStr), errors.New(retStr)
}

// open_trade function
// Create an open trade for an energy asset you have
func (t *SimpleChaincode) open_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var retStr string

	fmt.Println("Starting open_trade")
	// Arguments passed in the following order:
	// 0 --> "bob" == Owner
	// 1 --> "asset1" == Unique Identifier

	if len(args) != 2 {
		retStr = "Incorrect number of arguments. Expecting 4."
		return []byte(retStr), errors.New(retStr)
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
		retStr = "Invalid trade opening: " + owner + " does not own the asset " + id
		return []byte(retStr), errors.New(retStr)
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
			retStr = "Invalid trade opening: Asset for trade cannot be part of existing open trade"
			return []byte(retStr), errors.New(retStr)
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
	retStr = owner + " successfully opened trade for asset [" + id + "]."
	return []byte(retStr), nil
}

// Perform trade function
// Close an open trade and move ownership to the buyer
func (t *SimpleChaincode) perform_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var retStr string

	// Arguments are passed in the following order:
	// 0 --> "asset1" == Unique Id of the energy asset to trade
	// 1 --> "alice" == New owner
	if len(args) != 2 {
		retStr = "Incorrect number of arguments. Expecting 2."
		return []byte(retStr), errors.New(retStr)
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
				retStr = newOwner + " cannot accept their own trade order"
				return []byte(retStr), errors.New(retStr)
			}

			// Change the owner of the asset
			t.set_owner(stub, []string{id, newOwner})

			// Add the trade to the waiting trade list if it has an energy amount higher than 0
			if asset.Amount > 0 {
				// Get the WaitingTrades index
				var cTradesAsBytes []byte
				cTradesAsBytes, err = stub.GetState(chargingTradesStr)
				if err != nil {
					return nil, errors.New("Failed to get charging trades")
				}
				var chargingTrades ChargingTrades
				json.Unmarshal(cTradesAsBytes, &chargingTrades)
				// Add the trade to the waiting trades
				chargingTrades.WaitingTrades = append(chargingTrades.WaitingTrades, trades.OpenTrades[i])
				// Put the waiting trades back into the chaincode state
				cTradesAsBytes, err = json.Marshal(chargingTrades)
				err = stub.PutState(chargingTradesStr, cTradesAsBytes)
				if err != nil {
					return nil, err
				}
			}

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
	retStr = "Trade complete: " + newOwner + " now owns asset [" + id + "]."
	return []byte(retStr), nil
}

// Remove Open trade
// Close an open trade with no change of ownership taking place
func (t *SimpleChaincode) remove_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var retStr string

	// Only argument needed is the unique ID of the asset that should no longer be eligible for trade
	if len(args) != 2 {
		retStr = "Incorrect number of arguments. Expecting 2."
		return []byte(retStr), errors.New(retStr)
	}
	owner := args[0]
	id := args[1]

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
			fmt.Println("Found trade to remove, checking owner")
			// Verify owner is the one removing the trade
			if trades.OpenTrades[i].Owner != owner {
				retStr = "Error removing trade: Only trade order creator can remove this trade order."
				return []byte(retStr), errors.New(retStr)
			}
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
	retStr = "Open trade for asset [" + id + "] has been removed."
	return []byte(retStr), nil
}

// Get charging trades
// Return the next trade order that is waiting for its turn to transfer energy
func (t *SimpleChaincode) get_charging_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// Don't need any arguments to retrieve next trade to charge
	// Get the WaitingTrades index
	tradesAsBytes, err := stub.GetState(chargingTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get open trades")
	}
	var trades ChargingTrades
	json.Unmarshal(tradesAsBytes, &trades)

	// If there are any trades to get, return the first one
	if len(trades.WaitingTrades) > 0 {
		// Convert the first trade to []byte and return it
		var trade AnOpenTrade
		trade = trades.WaitingTrades[0]
		jsonAsBytes, _ := json.Marshal(trade)
		return jsonAsBytes, nil
	}

	// Return nothing if no waiting trades
	return nil, nil
}

func (t *SimpleChaincode) complete_charging_trade(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var retStr string
	// Trade details of the charge that has been completed should be passed in
	// Expecting 2: timestamp of trade and asset that was traded
	// Sufficient because same asset cannot be traded at the exact same time
	timestamp, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		retStr = "Something wrong with timestamp parameter: " + args[0]
		return []byte(retStr), errors.New(retStr)
	}
	id := args[1]

	// Get the WaitingTrades index
	tradesAsBytes, err := stub.GetState(chargingTradesStr)
	if err != nil {
		return nil, errors.New("Failed to get open trades")
	}
	var trades ChargingTrades
	json.Unmarshal(tradesAsBytes, &trades)

	// Try to find the matching trade in the waiting structure
	for i := range trades.WaitingTrades {
		var currentTrade AnOpenTrade
		currentTrade = trades.WaitingTrades[i]
		if currentTrade.Id == id && currentTrade.Timestamp == timestamp {
			// Found a match, remove it from the array
			trades.WaitingTrades = append(trades.WaitingTrades[:i], trades.WaitingTrades[i+1:]...)
			// Rewrite waiting trades to the chaincode state
			jsonAsBytes, _ := json.Marshal(trades)
			err = stub.PutState(chargingTradesStr, jsonAsBytes)
			if err != nil {
				return nil, err
			}
			// Return
			return []byte("Charging for trade of asset [" + id + "] at " + args[0] + " complete."), nil
		}
	}
	retStr = "Could not find asset [" + id + "] at " + args[0] + " in the charging trades queue."
	return []byte(retStr), errors.New(retStr)
}

// Clean up open trades
// Make sure all open trades are still possible/valid, and remove those that are not
func cleanTrades(stub shim.ChaincodeStubInterface) (err error) {
	// var didWork = false
	fmt.Println("Start cleaning trades")

	// Get the open trade struct
	tradesAsBytes, err := stub.GetState(openTradesStr)
	if err != nil {
		return errors.New("Failed to get open trades")
	}
	var trades AllTrades
	json.Unmarshal(tradesAsBytes, &trades)

	// Get the energy index
	energyAsBytes, err := stub.GetState(energyIndexStr)
	if err != nil {
		return errors.New("Failed to get energy index")
	}
	// Turn the energy index into a string array
	var energyIndex []string
	json.Unmarshal(energyAsBytes, &energyIndex)

	// Count the number of trades
	fmt.Println("# of open trades: " + strconv.Itoa(len(trades.OpenTrades)))
	// Iterate over all of the trades
	// Create a new list of all the trades that are valid
	var validTrades AllTrades
	for i := 0; i < len(trades.OpenTrades); i++ {
		// Look at every trade in the list
		currentTrade := trades.OpenTrades[i]
		// Determine if the asset still exists
		assetStillExists := false
		for j := range energyIndex {
			if currentTrade.Id == energyIndex[j] {
				assetStillExists = true
				break
			}
		}
		if assetStillExists == false {
			continue
		}
		// Determine if the asset is still owned by the person who initiated the trade
		// Get the state of the asset now that we know it exists
		assetAsBytes, assetErr := stub.GetState(currentTrade.Id)
		if assetErr != nil {
			return errors.New("Failed to get energy asset id")
		}
		res := Energy{}
		json.Unmarshal(assetAsBytes, &res)
		// If the owner of the asset is not the person who initiated trade, invalid trade
		if currentTrade.Owner != res.Owner {
			continue
		}
		// Set the price in the trade to the price of the asset (in case a price was changed)
		currentTrade.Amount = res.Price
		// Transaction is still valid, add it to the validTrades
		// Append the energy asset to the energy index
		validTrades.OpenTrades = append(validTrades.OpenTrades, currentTrade)
	}

	// Write the valid trades to the current state
	fmt.Println("Writing valid trades to chaincode state")
	jsonAsBytes, _ := json.Marshal(validTrades)
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return err
	}

	// Successful exit
	fmt.Println("End trade cleaning")
	return nil

}
