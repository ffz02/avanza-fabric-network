package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
)

var (
	warningChars = []string{"'", "--", "&"}
	escapedChars = []string{"\\'", "", ""}
)

func hash(val string) string {
	hasher := sha512.New()
	hasher.Write([]byte(val))
	sha1_hash := hex.EncodeToString(hasher.Sum(nil))
	return sha1_hash
}

func sanitize(input interface{}, t string) interface{} {
	m := input.(string)
	switch t {
	case "bool":
		feetFloat, _ := strconv.ParseBool(strings.TrimSpace(m))
		return feetFloat
	case "float64":
		feetFloat, _ := strconv.ParseFloat(strings.TrimSpace(m), 64)
		return feetFloat
	case "float32":
		feetFloat, _ := strconv.ParseFloat(strings.TrimSpace(m), 32)
		return feetFloat
	case "string":
		outString := m
		for i := 0; i < len(warningChars); i++ {
			outString = strings.Replace(outString, warningChars[i], escapedChars[i], -1)
		}
		return outString
	case "int64":
		intVal, _ := strconv.ParseInt(strings.TrimSpace(m), 10, 64)
		return intVal
	case "int":
		intVal, _ := strconv.Atoi(strings.TrimSpace(m))
		return intVal
	default:
		panic(fmt.Sprintf("unexpected type: %T", m))
	}
}

func deleteDataFromCollection(stub *hypConnect, args []string, functionName string) error {
	key := sanitize(args[0], "string").(string)
	collectionName := sanitize(args[1], "string").(string)

	err := deleteData(stub, key, collectionName)
	if err != nil {
		return err
	}
	return nil
}

type hypConnect struct {
	Connection shim.ChaincodeStubInterface
	EventList  []eventDataFormat
}

func (ref *hypConnect) AddEvent(event eventDataFormat) []eventDataFormat {
	ref.EventList = append(ref.EventList, event)
	return ref.EventList
}

type generalEventStruct struct {
	EventName      string            `json:"eventName"`
	EventList      []eventDataFormat `json:"events"`
	AdditionalData interface{}       `json:"additionalData"`
}

type eventDataFormat struct {
	Key        string `json:"Key"`
	Collection string `json:"Collection"`
}
