package main

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

//Define data structure
type File struct {
	Tag string `json:"tag"` //for Find function
	Data []byte //File as Bytes
}

// QueryResult structure used for handling result of query
type QueryResult struct { 
	Key	string `json:"Key"`
	Tag	string `json:"Tag"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "upload" {
		return s.upload(APIstub, args)
	} else if function == "download" {
		return s.download(APIstub, args)
	} else if function == "list" {
		return s.queryAllFiles(APIstub)
	} else if function == "show" {
		return s.queryFile(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getFile(APIstub shim.ChaincodeStubInterface, key string) (*File, error) {

	fileAsBytes, err := APIstub.GetState(key)
	
	if err != nil {
		return nil, err
	}
	if fileAsBytes == nil {
		return nil, nil
	}

	file := new(File)
	json.Unmarshal(fileAsBytes, file)

	return file, nil
}

func (s *SmartContract) putFile(APIstub shim.ChaincodeStubInterface, key string, file *File) error {

	fileAsBytes, err := json.Marshal(file)
	
	if err != nil {
		return err
	}
	if err := APIstub.PutState(key, fileAsBytes); err != nil {
		return err
	}

	return nil
}

func (s *SmartContract) upload(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 { //args : key, tag, bytes
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	file, err := s.getFile(APIstub, args[0]); 
	if err != nil {
		Resp := "Failed to get state for " + args[0]
		return shim.Error(Resp)
	}
	if file != nil { //check if key already exist
		Resp := "file already exist: " + args[0]
		return shim.Error(Resp)
	}

	file = new(File)
	file.Tag = args[1]
	file.Data = []byte(args[2])

	if err := s.putFile(APIstub, args[0], file); err != nil {
		Resp := "Failed to put state for " + args[0]
		return shim.Error(Resp)
	}
	
	return shim.Success(nil)
}

func (s *SmartContract) download(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	file, err := s.getFile(APIstub, args[0]); 
	if err != nil {
		Resp := "Failed to get state for " + args[0]
		return shim.Error(Resp)
	}
	if file == nil { //check if key doesn't exist
		Resp := "file does not exist: " + args[0]
		return shim.Error(Resp)
	}

	return shim.Success(file.Data)
}

func (s *SmartContract) queryFile(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	file, err := s.getFile(APIstub, args[0]); 
	if err != nil {
		Resp := "Failed to get state for " + args[0]
		return shim.Error(Resp)
	}
	if file == nil { //check if key doesn't exist
		Resp := "file does not exist: " + args[0]
		return shim.Error(Resp)
	}
	queryResult := QueryResult{Key: args[0], Tag: file.Tag}
	resultAsBytes, err := json.Marshal(queryResult)

	return shim.Success(resultAsBytes)
}

func (s *SmartContract) queryAllFiles(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}

		file := new(File)
		_ = json.Unmarshal(queryResponse.Value, file)

		queryResult := QueryResult{Key: queryResponse.Key, Tag: file.Tag}
		results = append(results, queryResult)
	}
	resultsAsBytes, err := json.Marshal(results)

	return shim.Success(resultsAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
