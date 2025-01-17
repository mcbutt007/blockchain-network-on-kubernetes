/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.
import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func intToByteArray(i int) []byte {
	s := strconv.Itoa(i)
	return []byte(s)
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "invoke":
		// Make payment of X units from A to B
		return t.transfer(stub, args)
	case "add": // modified code by mancitiss
		// Add an entity to its state
		return t.addKey(stub, args)
	case "addVirus":
		// Add an entity to its state
		return t.addVirus(stub, args)
	case "delete":
		// Deletes an entity from its state
		return t.deleteKey(stub, args)
	case "deleteVirus":
		// Deletes an entity from its state
		return t.deleteVirus(stub, args)
	case "query":
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	default:
		return shim.Error("Invalid invoke function name. Expecting \"transfer\" \"add\" \"addVirus\" \"delete\" \"deleteVirus\" \"query\"")
	}

	return shim.Error("Invalid invoke function name. Expecting \"transfer\" \"add\" \"addVirus\" \"delete\" \"deleteVirus\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}
	bytes := intToByteArray(0)
	return shim.Success(bytes)
}

// Add an entity to state
func (t *SimpleChaincode) addKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 || len(args) > 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1 or 2")
	}

	key := args[0]
	amount := 0

	if len(args) == 2 {
		amount, _ = strconv.Atoi(args[1])
	}

	// Check if the key already exists in the ledger
	if value, err := stub.GetState(key); err != nil || value != nil {
		return shim.Error("Key already exists")
	}

	// Key doesn't exist, add it to the ledger
	err := stub.PutState(key, []byte(strconv.Itoa(amount)))
	if err != nil {
		return shim.Error(err.Error())
	}
	bytes := intToByteArray(0)
	return shim.Success(bytes)
}

// Add a virus to state
func (t *SimpleChaincode) addVirus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 || len(args) > 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1 or 2")
	}

	creator := args[0]
	signature := args[1]

	// Check if the key already exists in the ledger
	if value, err := stub.GetState(signature); err != nil || value != nil {
		return shim.Error("Key already exists")
	}

	// Key doesn't exist, add it to the ledger
	err := stub.PutState(signature, []byte(creator))
	if err != nil {
		return shim.Error(err.Error())
	}

	// check if creator exist
	if value, err := stub.GetState(creator); err != nil || value == nil {
		// creator doesn't exist, add it to the ledger
		err := stub.PutState(creator, []byte("1"))
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	// increment creator's virus count
	creatorVirusCount, _ := stub.GetState(creator)
	creatorVirusCountInt, _ := strconv.Atoi(string(creatorVirusCount))
	creatorVirusCountInt++
	err = stub.PutState(creator, []byte(strconv.Itoa(creatorVirusCountInt)))
	if err != nil {
		return shim.Error(err.Error())
	}
	bytes := intToByteArray(0)
	return shim.Success(bytes)
}

// Deletes a virus from state
func (t *SimpleChaincode) deleteVirus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting signature of the virus to delete")
	}

	signature := args[0]

	// get creator
	creator, err := stub.GetState(signature)
	if err != nil {
		return shim.Error("Failed to get virus state")
	}

	// decrement creator's virus count
	creatorVirusCount, _ := stub.GetState(string(creator))
	creatorVirusCountInt, _ := strconv.Atoi(string(creatorVirusCount))
	creatorVirusCountInt--
	err = stub.PutState(string(creator), []byte(strconv.Itoa(creatorVirusCountInt)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Delete the key from the state in ledger
	err = stub.DelState(signature)
	if err != nil {
		return shim.Error("Failed to delete virus state")
	}
	bytes := intToByteArray(0)
	return shim.Success(bytes)
}

// Deletes an entity from state
func (t *SimpleChaincode) deleteKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	bytes := intToByteArray(0)
	return shim.Success(bytes)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
