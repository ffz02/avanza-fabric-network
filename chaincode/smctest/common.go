package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
)

func GetTxID(stub *hypConnect) (string, string) {
	str := stub.Connection.GetTxID()
	_, args := stub.Connection.GetFunctionAndParameters()
	fmt.Println("\n\nARG >> ", args)
	fmt.Println("\n\n\n\n ARGS[req] >> ", args[len(args)-1])
	fmt.Println("\n\n")
	str2 := args[len(args)-1]
	return str, str2
}

func getCidData(stub shim.ChaincodeStubInterface) (string, error) {
	return cid.GetMSPID(stub)
}

func insertData(stub *hypConnect, key string, privateCollection string, data []byte) error {

	err := stub.Connection.PutPrivateData(privateCollection, key, data)
	if err != nil {
		return err
	}

	event := eventDataFormat{}
	event.Key = key
	event.Collection = privateCollection
	stub.EventList = stub.AddEvent(event)

	fmt.Println("Successfully Put State for Key: " + key + " and Private Collection " + privateCollection)
	return nil
}

func deleteData(stub *hypConnect, key string, privateCollection string) error {

	err := stub.Connection.DelPrivateData(privateCollection, key)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Delete for Key: " + key + " and Private Collection " + privateCollection)

	event := eventDataFormat{}
	event.Key = key
	event.Collection = privateCollection
	stub.EventList = stub.AddEvent(event)

	return nil
}
func fetchData(stub hypConnect, key string, privateCollection string) ([]byte, error) {
	bytes, err := stub.Connection.GetPrivateData(privateCollection, key)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func getArguments(stub shim.ChaincodeStubInterface) ([]string, error) {
	transMap, err := stub.GetTransient()
	if err != nil {
		return nil, err
	}
	if _, ok := transMap["PrivateArgs"]; !ok {
		return nil, errors.New("PrivateArgs must be a key in the transient map")
	}
	fmt.Println("Arguments: %v", transMap)
	generalInput := string(transMap["PrivateArgs"])
	retVal := strings.Split(generalInput, "|")
	return retVal, nil
}
func getOrgTypeByMSP(stub shim.ChaincodeStubInterface, MSP string) (string, error) {

	MSPMappingAsBytes, err := stub.GetState("MSPMapping")
	if err != nil {
		return "", err
	}

	if err != nil {
		fmt.Println("MSPMapping - Failed to get state MSP mapping information." + err.Error())
		return "", err
	} else if MSPMappingAsBytes != nil {
		fmt.Println("MSPMapping - This data Fetched from Transactions.")
		var MSPListUnmarshaled []MSPList

		err := json.Unmarshal(MSPMappingAsBytes, &MSPListUnmarshaled)
		if err != nil {
			fmt.Println("MSPMapping-Failed to UnMarshal state.")
			return "", err
		}
		fmt.Println("Unmarshaled: %v", MSPListUnmarshaled)
		for i := 0; i < len(MSPListUnmarshaled); i++ {
			if MSPListUnmarshaled[i].MSP == MSP {
				fmt.Println("OrgType for MSP " + MSP + " is " + MSPListUnmarshaled[i].OrgType)
				return MSPListUnmarshaled[i].OrgType, nil
			}
		}
	}
	return "", nil
}
func RaiseEventData(stub hypConnect, eventName string, args ...interface{}) (string, error) {

	var eventList generalEventStruct
	eventList.EventName = eventName
	eventList.EventList = stub.EventList
	eventList.AdditionalData = args
	eventJSONasBytes, err2 := json.Marshal(eventList)
	if err2 != nil {
		return "", err2
	}
	fmt.Println("Event raised: " + eventName)
	//fmt.Println("\neventJSONasBytes : ", eventList.EventName+"\n")
	mEventName := eventList.EventName
	err3 := stub.Connection.SetEvent("chainCodeEvent", []byte(eventJSONasBytes))
	if err3 != nil {
		return "", err3
	}
	var err4 error
	err4 = nil
	return mEventName, err4

}

func loadXpath(jsonPath string, data map[string]interface{}) (string, interface{}) {

	//	fmt.Println("Path to get " + jsonPath)

	if data[jsonPath] == nil {

		return "No element", nil
	}

	if reflect.TypeOf(data[jsonPath]) == reflect.TypeOf((*string)(nil)).Elem() {
		return "string", data[jsonPath]
	} else if reflect.TypeOf(data[jsonPath]) == reflect.TypeOf((*[]interface{})(nil)).Elem() {

		return "array", data[jsonPath]

	}
	return "object", data[jsonPath]

}

func genericErrorHandler(err error, arrErr []string, errorCode string, mapMainError bool, force bool) []byte {
	fmt.Println("mainError ", err)
	fmt.Println("extraArrayString ", arrErr)

	if !force && err == nil {
		return nil
	}
	var options []string
	options = append(options, arrErr...)
	if mapMainError {
		options = append(options, err.Error())
	}
	errBytes, _ := prepareErrorCode(errorCode, options)
	return errBytes
}

func prepareErrorCode(errorCode string, options []string) ([]byte, error) {
	errCode := errCode{}
	errCode.ErrorCode = errorCode
	errCode.Options = options

	errCodeasBytes, errCodeMarshalError := json.Marshal(errCode)
	if errCodeMarshalError != nil {
		return nil, errCodeMarshalError
	}
	//return stub.Success(errCodeasBytes)
	return errCodeasBytes, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func getTaskTypes() TaskTypes {
	returnVal := &TaskTypes{
		SIGNATURE_REQUEST:     		"SIGNATURE_REQUEST",
		PSN_UPDATE:                 "PSN_UPDATE",
		KOC_TIMINGS:                "KOC_TIMINGS",
		ACTUAL_LOAD_QUANTITY:       "ACTUAL_LOAD_QUANTITY",
		GENERATE_PACKAGE:           "GENERATE_PACKAGE",
		FORWARD_TO_CUSTOMER:        "FORWARD_TO_CUSTOMER",
		PRINT_DOCUMENTS: 			"PRINT_DOCUMENTS",
		UPLOAD_FINAL_SET: 			"UPLOAD_FINAL_SET",
	}
	return *returnVal
}
