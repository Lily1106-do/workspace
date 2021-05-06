package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer")

type SimpleChaincode struct {
}

var numID string
var chName = "channel"

// The length of the ID to be filled
var lenID = 4

// chaincode for user
type Producer struct {
	ID             string  `json:"ID"`
	Freeze         bool    `json:"Freeze"`
	TotalPower     float64 `json:"TotalPower"`
	Price          float64 `json:"Price"`
	Asset          float64 `json:"Asset"`
}

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
	if fn == "addProducer" {
		return t.addProducer(stub, args)
	} else if fn == "queryProducer" {
		return t.queryProducer(stub, args)
	} else if fn == "getTotalPower" {
		return t.getTotalPower(stub, args)
	} else if fn == "getPrice" {
		return t.getPrice(stub, args)
	} else if fn == "setPrice" {
		return t.setPrice(stub, args)
	} else if fn == "getAllSeller" {
		return t.getAllSeller(stub, args)
	} else if fn == "sellPower" {
		return t.sellPower(stub, args)
	} else if fn == "setPower" {
		return t.setPower(stub, args)
	}

	// error function name
	return shim.Error("Invalid invoke function name.")
}

// create user account
func (t *SimpleChaincode) addProducer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect arguments.Excepting 0")
	}
	idBytes, _ := stub.GetState(numID)
	id, _ := strconv.Atoi(string(idBytes))
	producerJson := Producer{fmt.Sprintf("%0" + strconv.Itoa(lenID) + "d", id), false, 0, 0, 0}
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

// Use ID find a producer
func (t *SimpleChaincode) queryProducer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments.")
	}
	id := args[0]
	userValBytes, _ := stub.GetState(id)
	if userValBytes == nil {
		shim.Error("ID does not exeist:" + id + ".")
	}
	return shim.Success(userValBytes)
}

// get total power by ID
func (t *SimpleChaincode) getTotalPower(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments.")
	}
	producerBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var producer Producer
	_ = json.Unmarshal(producerBytes, &producer)
	return shim.Success([]byte(strconv.FormatFloat(producer.TotalPower, 'f', 10, 64)))
}

// get price by ID
func (t *SimpleChaincode) getPrice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments.")
	}
	producerBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var producer Producer
	_ = json.Unmarshal(producerBytes, &producer)
	return shim.Success([]byte(strconv.FormatFloat(producer.Price, 'f', 10, 64)))
}

// set price by ID
func (t *SimpleChaincode) setPrice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.")
	}
	producerBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var producer Producer
	_ = json.Unmarshal(producerBytes, &producer)
	producer.Price, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		shim.Error(err.Error())
	}
	producerJsonBytes, _ := json.Marshal(producer)
	_ = stub.PutState(producer.ID, producerJsonBytes)
	return shim.Success(nil)
}

// get all available seller
func (t *SimpleChaincode) getAllSeller(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var buffer bytes.Buffer
	if len(args) != 0 {
		return shim.Error("Incorrect arguments.Excepting 0")
	}
	bArrayMemberHasWritten := false
	numOfID, _ := stub.GetState(numID)
	resultIterator, _ := stub.GetStateByRange(fmt.Sprintf("%0"+strconv.Itoa(lenID)+"d", 0), fmt.Sprintf("%0"+strconv.Itoa(lenID)+"d", string(numOfID)))
	buffer.WriteString("{\"allAvailableSellers\":[")
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to read iterator data" + err.Error() +"\"}"
			return shim.Error(jsonResp)
		}
		if bArrayMemberHasWritten == true {
			buffer.WriteString(",")
		}
		buffer.Write(queryResponse.Value)
		bArrayMemberHasWritten = true
	}
	buffer.WriteString("]}")
	return shim.Success(buffer.Bytes())
}

// sell power by IDs and amount
func (t *SimpleChaincode) sellPower(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.")
	}
	buyerID := args[0]
	sellerID := args[1]
	amount, _ := strconv.ParseFloat(args[1], 64)
	producerBytes, err := stub.GetState(sellerID)
	if err != nil {
		shim.Error(err.Error())
	}
	var producer Producer
	_ = json.Unmarshal(producerBytes, &producer)
	if amount > producer.TotalPower {
		return shim.Error("Insufficient power")
	}
	producer.TotalPower -= amount
	gains := amount * producer.Price
	pResponse := doOtherCC(stub, chName, "usercc", "getAsset", buyerID)
	if pResponse.Status != 200{
		jsonResp := "{\"Error\":\"Failed to invoke from ChainCode:" + pResponse.String() + "\"}"
		return shim.Error(jsonResp)
	}
	var user UserInfo
	err = json.Unmarshal(pResponse.Payload, & user)
	if user.Asset < gains {
		return shim.Error("Insufficient asset for user:" + user.ID)
	}
	pResponse = doOtherCC(stub, chName, "usercc", "setAsset", buyerID, strconv.FormatFloat(user.Asset - gains, 'f', 30, 64))
	if pResponse.Status != 200{
		jsonResp := "{\"Error\":\"Failed to invoke from ChainCode:" + pResponse.String() + "\"}"
		return shim.Error(jsonResp)
	}
	pResponse = doOtherCC(stub, chName, "usercc", "setElec", buyerID, strconv.FormatFloat(user.Elec + amount, 'f', 30, 64))
	if pResponse.Status != 200{
		jsonResp := "{\"Error\":\"Failed to invoke from ChainCode:" + pResponse.String() + "\"}"
		return shim.Error(jsonResp)
	}
	producer.Asset += gains
	producer.TotalPower -= amount
	producerJsonBytes, _ := json.Marshal(producer)
	_ = stub.PutState(sellerID, producerJsonBytes)
	response := args[0] +" sells "+ args[1] + " producer with price:" + strconv.Itoa(int(producer.Price)) + ", gains:" + strconv.Itoa(int(gains))
	return shim.Success([]byte(response))
}

// set price by ID
func (t *SimpleChaincode) setPower(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect arguments.")
	}
	producerBytes, err := stub.GetState(args[0])
	if err != nil {
		shim.Error(err.Error())
	}
	var producer Producer
	_ = json.Unmarshal(producerBytes, &producer)
	producer.Asset, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		shim.Error(err.Error())
	}
	producerJsonBytes, _ := json.Marshal(producer)
	_ = stub.PutState(producer.ID, producerJsonBytes)
	return shim.Success(nil)

}

//// get total power by ID
//func (t *SimpleChaincode) withdrawMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 2 {
//		return shim.Error("Incorrect arguments.")
//	}
//	producerID := args[0]
//	amount, _ := strconv.ParseFloat(args[1], 64)
//	producerBytes, err := stub.GetState(producerID)
//	if err != nil {
//		shim.Error(err.Error())
//	}
//	var producer Producer
//	_ = json.Unmarshal(producerBytes, &producer)
//	if amount > producer.Asset {
//		return shim.Error("Insufficient asset")
//	}
//	producer.Asset -= amount
//	method := `TransferFrom`
//	fromAccount := ""
//	toAccount := ""
//	params := []string{method, fromAccount, toAccount, args[1]}
//	argsBytess := make([][]byte, len(params))
//	for i, arg := range params {
//		argsBytess[i] = []byte(arg)
//	}
//	pResponse := stub.InvokeChaincode("tokencc",argsBytess,"")
//	if pResponse.Status != 200{
//		jsonResp := "{\"Error\":\"Failed to invoke from ChainCode:procc" + pResponse.String() + "\"}"
//		return shim.Error(jsonResp)
//	}
//	return shim.Success(pResponse.Payload)
//}

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
		fmt.Printf("Error starting Producer chaincode: %s", err)
	}
}

