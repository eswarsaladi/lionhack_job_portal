package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

type IDM struct {
}

func main() {
	err := shim.Start(new(IDM))
	if err != nil {
	}
}

type GeneralInfo struct {
	ObjectType        string `json:"docType"`
	UserName          string `json:"username"`
	Name              string `json:"name"`
	FatherHusbandName string `json:"fhname"`
	Age               string `json:"age"`
	Gender            string `json:"gender"`
	DOB               string `json:"dob"`
	PAN               string `json:"pan"`
	Aadhar            string `json:"aadhar"`
	CAddress          string `json:"currentaddress"`
	PAddress          string `json:"permanentaddress"`
	CibilScore        string `json:"cibilscore"`
	CibilDate         string `json:"cibildate"`
	UniqueNumber      string `json:"uniquenumber"`
}

type DocInfo struct {
	ObjectType  string `json:"docType"`
	UniqueKey   string `json:"uniquekey"`
	Certificate string `json:"certificate"`
	Image       string `json:"image"`
}

func (s *IDM) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *IDM) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "InitGeneralInfo" {
		return s.InitGeneralInfo(stub, args)
	} else if function == "ReadGeneralInfo" {
		return s.ReadGeneralInfo(stub, args)
	} else if function == "UpdateGeneralInfo" {
		return s.UpdateGeneralInfo(stub, args)
	} else if function == "InitDocument" {
		return s.InitDocument(stub, args)
	} else if function == "ReadDocument" {
		return s.ReadDocument(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return shim.Error("Received unknown function ")
}

func inputSanitiser(stub shim.ChaincodeStubInterface, args []string) bool {

	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return false
		}
	}
	return true
}

//InitGeneralInfo

func (s *IDM) InitGeneralInfo(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("Init General Details")

	if len(args) != 13 {
		return shim.Error("Incorrect number of Arguments, Expecting 11 arguments")
	}

	isTrue := inputSanitiser(stub, args)
	if isTrue != true {
		return shim.Error("argument must be a non-empty string")
	}

	username := args[0]
	name := args[1]
	fhname := args[2]
	age := args[3]
	gender := args[4]
	dob := args[5]
	pan := args[6]
	aadhar := args[7]
	currentaddress := args[8]
	permanentaddress := args[9]
	cibilscore := args[10]
	cibildate := args[11]
	uniquenumber := args[12]

	infoAsBytes, err := stub.GetState(username)
	if err != nil {
		return shim.Error("Failed to get User Name: " + err.Error())
	} else if infoAsBytes != nil {
		fmt.Println("This User Name already exists: " + username)
		return shim.Error("This User Name already exists: " + username)
	}

	objectType := "info"
	info := &GeneralInfo{objectType, username, name, fhname, age, gender, dob, pan, aadhar, currentaddress, permanentaddress, cibilscore, cibildate, uniquenumber}
	InfoJSONasBytes, err := json.Marshal(info)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(username, InfoJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init General Info")
	return shim.Success(nil)
}

//ReadGeneralInfo

func (s *IDM) ReadGeneralInfo(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	var username, jsonResp string

	fmt.Println("running ReadGeneralInfo()")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	isTrue := inputSanitiser(stub, args)
	if isTrue != true {
		return shim.Error("argument must be a non-empty string")
	}

	username = args[0]
	infoAsbytes, err := stub.GetState(username)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + username + "\"}"
		return shim.Error(jsonResp)
	} else if infoAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + username + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(infoAsbytes)
}

//UpdateGeneralInfo

func (s *IDM) UpdateGeneralInfo(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running UpdateGeneralInfo()")

	if len(args) != 13 {
		return shim.Error("Incorrect number of arguments. Expecting 11")
	}

	isTrue := inputSanitiser(stub, args)
	if isTrue != true {
		return shim.Error("argument must be a non-empty string")
	}

	InfoDetailAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get User Name")
	}
	InfoDetail := GeneralInfo{}
	json.Unmarshal(InfoDetailAsBytes, &InfoDetail)

	InfoDetail.Name = args[1]

	InfoDetail.FatherHusbandName = args[2]
	InfoDetail.Age = args[3]
	InfoDetail.Gender = args[4]
	InfoDetail.DOB = args[5]
	InfoDetail.PAN = args[6]
	InfoDetail.Aadhar = args[7]
	InfoDetail.CAddress = args[8]
	InfoDetail.PAddress = args[9]
	InfoDetail.CibilScore = args[10]
	InfoDetail.CibilDate = args[11]

	jsonAsBytes, _ := json.Marshal(InfoDetail)
	err = stub.PutState(args[0], jsonAsBytes)
	if err != nil {
		return shim.Error(" Failed to update General Info")
	}

	err = stub.PutState(args[0], jsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *IDM) InitDocument(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("Init Documents")

	if len(args) != 3 {
		return shim.Error("Incorrect number of Arguments, Expecting 1 arguments")
	}

	isTrue := inputSanitiser(stub, args)
	if isTrue != true {
		return shim.Error("argument must be a non-empty string")
	}

	uniquekey := args[0]

	certificate := args[1]
	image := args[2]

	docAsBytes, err := stub.GetState(uniquekey)
	if err != nil {
		return shim.Error("Failed to get Unique Key: " + err.Error())
	} else if docAsBytes != nil {
		fmt.Println("This Unique Key already exists: " + uniquekey)
		return shim.Error("This Unique Key already exists: " + uniquekey)
	}

	objectType := "doc"
	doc := &DocInfo{objectType, uniquekey, certificate, image}
	docJSONasBytes, err := json.Marshal(doc)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(uniquekey, docJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init Document")
	return shim.Success(nil)
}

func (s *IDM) ReadDocument(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	var uniquekey, jsonResp string

	fmt.Println("running ReadDocument()")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	isTrue := inputSanitiser(stub, args)
	if isTrue != true {
		return shim.Error("argument must be a non-empty string")
	}

	uniquekey = args[0]
	docAsbytes, err := stub.GetState(uniquekey)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + uniquekey + "\"}"
		return shim.Error(jsonResp)
	} else if docAsbytes == nil {
		jsonResp = "{\"Error\":\"Unique Key does not exist: " + uniquekey + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(docAsbytes)
}
