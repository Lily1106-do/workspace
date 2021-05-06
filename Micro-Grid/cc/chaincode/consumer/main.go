package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

var numID="numID"

// The length of the ID to be filled
var lenID = 6

var chName = "myc"

// chaincode for user
type UserInfo struct {
	ID     string  `json:"ID"`
	Freeze bool    `json:"Freeze"`
	Asset  float64 `json:"Asset"`
	Elec   float64 `json:"Elec"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var s []byte
	result, err := stub.GetState(numID)
	if err != nil {
		return 	shim.Error(err.Error())
	}
	if strings.EqualFold(string(s), string(result)) {
		err := stub.PutState(numID, []byte("1"))
		println("write to ledger")
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to write to ledger:" + err.Error() + "\"}"
			return shim.Error(jsonResp)
		}
	}
	fmt.Print("chaincode init success")
	return shim.Success([]byte("init success"))
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	if fn == "addUser" {
		return t.addUser(stub, args)
	} else if fn == "findUser" {
		return t.findUser(stub, args)
	} else if fn == "getAsset" {
		return t.getAsset(stub, args)
	} else if fn == "getElec" {
		return t.getElec(stub, args)
	} else if fn == "getSellers" {
		return t.getSellers(stub, args)
	} else if fn == "freezeAccount" {
		return t.freezeAccount(stub, args)
	} else if fn == "setAsset" {
		return t.setAsset(stub, args)
	} else if fn == "setElec" {
		return t.setElec(stub, args)
	}

	// error function name
	return shim.Error("Invalid invoke function name:" + fn)
}

// Use ID create user account
func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect arguments.Excepting 0")
	}
	idBytes, _ := stub.GetState(numID)
	id, _ := strconv.Atoi(string(idBytes))
	producerJson := UserInfo{fmt.Sprintf("%0" + strconv.Itoa(lenID) + "d", id), false, 0, 0}
	producerJsonAsBytes, err := json.Marshal(producerJson)
	if err != nil {
		shim.Error(err.Error())
	}
	err = stub.PutState(string(idBytes), producerJsonAsBytes)
	if err != nil {
		shim.Error(err.Error())
	}
	id = id + 1
	err = stub.PutState(numID, []byte(strconv.Itoa(id)))
	return shim.Success(producerJsonAsBytes)
}

// Use ID create user account
func (t *SimpleChaincode) findUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments.")
	}
	userID := args[0]
	userValBytes, _ := stub.GetState(userID)
	if userValBytes == nil {
		shim.Error("ID does not exeist:" + userID + ".")
	}
	return shim.Success(userValBytes)
}

// get asset by ID
func (t *SimpleChaincode) getAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments.")
	}
	userValBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var userVal	UserInfo
	_ = json.Unmarshal(userValBytes, &userVal)
	return shim.Success([]byte(strconv.FormatFloat(userVal.Asset, 'f', 30, 64)))
}

// get producer by ID
func (t *SimpleChaincode) getElec(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments.")
	}
	userValBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var userVal	UserInfo
	_ = json.Unmarshal(userValBytes, &userVal)
	return shim.Success([]byte(strconv.FormatFloat(userVal.Elec, 'f', 30, 64)))
}

//// recharge token
//func (t *SimpleChaincode) rechargeToken(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 2 {
//		return shim.Error("Incorrect arguments.")
//	}
//	producerID := args[0]
//	amount, _ := strconv.ParseFloat(args[1], 64)
//	producerBytes, err := stub.GetState(producerID)
//	if err != nil {
//		shim.Error(err.Error())
//	}
//	var user UserInfo
//	err = json.Unmarshal(producerBytes, &user)
//	if amount > user.Asset {
//		return shim.Error("Failed to unmarshal")
//	}
//	user.Asset += amount
//	return doOtherCC(stub, "", "tokencc", "Transfer", "", args[1])
//}

// find all seller
func (t *SimpleChaincode) getSellers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect arguments.")
	}
	return doOtherCC(stub, chName, "procc", "getAllSeller")
}

// freeze or unfreeze a account by ID
func (t *SimpleChaincode) freezeAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.")
	}
	userValBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var userVal	UserInfo
	_ = json.Unmarshal(userValBytes, &userVal)
	freezeBool, err := strconv.ParseBool(args[1])
	if err != nil {
		shim.Error(err.Error())
	}
	userVal.Freeze = freezeBool
	userJsonBytes, _ := json.Marshal(userVal)
	err = stub.PutState(userVal.ID, userJsonBytes)
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// set asset by ID
func (t *SimpleChaincode) setAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.")
	}
	userValBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var userVal	UserInfo
	_ = json.Unmarshal(userValBytes, &userVal)
	userVal.Asset, err = strconv.ParseFloat(args[1], 64)
	userJsonBytes, _ := json.Marshal(userVal)
	err = stub.PutState(userVal.ID, userJsonBytes)
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)}

// set producer by ID
func (t *SimpleChaincode) setElec(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.")
	}
	userValBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var userVal	UserInfo
	_ = json.Unmarshal(userValBytes, &userVal)
	userVal.Elec, err = strconv.ParseFloat(args[1], 64)
	userJsonBytes, _ := json.Marshal(userVal)
	err = stub.PutState(userVal.ID, userJsonBytes)
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}


func doOtherCC(stub shim.ChaincodeStubInterface, channelName string, ccName string, params ...string) pb.Response {
	var paramStr []string
	for _, val := range params {
		paramStr = append(paramStr, val)
	}
	argsBytess := make([][]byte, len(params))
	for i, arg := range params {
		argsBytess[i] = []byte(arg)
	}
	pResponse := stub.InvokeChaincode(ccName, argsBytess, channelName)
	if pResponse.Status != 200{
		jsonResp := "{\"Error\":\"Failed to invoke from ChainCode:" + pResponse.String() + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(pResponse.Payload)
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleChaincode)); err != nil {
		fmt.Printf("Error starting UserInfo chaincode: %s", err)
	}
}

