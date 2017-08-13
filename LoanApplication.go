package main

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim" 
)

var logger = shim.NewLogger("mylogger")

type SampleChaincode struct {
}

//custom data models
type PersonalInfo struct {
	Id string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	DOB       string `json:"dob"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
}

func GetInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering GetInfo")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing loan application ID")
	}

	var id = args[0]
	bytes, err := stub.GetState(id)
	if err != nil {
		logger.Error("Could not fetch person with id "+id+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreatePersonalInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering CreatePersonalInfo")

	if len(args) < 6 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected six arguments for info creation")
	}
	
	var id = args[0]
	var firstName = args[1]
	var lastName = args[2]
	var dob = args[3]
	var email = args[4]
	var mobile = args[5]
	
	input := `{
		"id": "` + id + `", 
		"firstname": "` + firstName + `", 
		"lastname": ` + lastName + `
		"dob": "` + dob + `", 
		"email": "` + email + `", 
		"mobile": ` + mobile + `
	}`

	err := stub.PutState(id, []byte(input))
	if err != nil {
		logger.Error("Could not save info to ledger", err)
		return nil, err
	}

	var customEvent = "{eventType: 'infoCreation', description:" + id + "' Successfully created'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully saved info")
	return nil, nil

}

/**
Updates the status of the loan application
**/
func UpdatePersonalInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering UpdatePersonalInfo")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for update")
	}

	var loanApplicationId = args[0]
	var mobile = args[1]

	laBytes, err := stub.GetState(loanApplicationId)
	if err != nil {
		logger.Error("Could not fetch info from ledger", err)
		return nil, err
	}
	var loanApplication PersonalInfo
	err = json.Unmarshal(laBytes, &loanApplication)
	loanApplication.Mobile = mobile

	laBytes, err = json.Marshal(&loanApplication)
	if err != nil {
		logger.Error("Could not marshal info post update", err)
		return nil, err
	}

	err = stub.PutState(loanApplicationId, laBytes)
	if err != nil {
		logger.Error("Could not save info post update", err)
		return nil, err
	}

	var customEvent = "{eventType: 'infoUpdate', description:" + loanApplicationId + "' Successfully updated status'}"
	err = stub.SetEvent("evtSender", []byte(customEvent))
	if err != nil {
		return nil, err
	}
	logger.Info("Successfully updated info")
	return nil, nil

}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetInfo" {
		return GetInfo(stub, args)
	}
	return nil, nil
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreatePersonalInfo" {
		return CreatePersonalInfo(stub, args)
	} else if function == "UpdatePersonalInfo" {
		return UpdatePersonalInfo(stub, args)
	}
	return nil, nil
}

type customEvent struct {
	Type       string `json:"type"`
	Decription string `json:"description"`
}

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))

	err := shim.Start(new(SampleChaincode))
	if err != nil {
		logger.Error("Could not start SampleChaincode")
	} else {
		logger.Info("SampleChaincode successfully started")
	}

}
