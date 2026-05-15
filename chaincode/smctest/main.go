package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	// "github.com/hyperledger/fabric/core/chaincode/lib/cid"
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	// pb "github.com/hyperledger/fabric/protos/peer"
)

type CoreChainCode struct {
}

//var logger = shim.NewLogger("Core")

func main() {

	//TestAll()
	fmt.Println("Core ChainCode Started")
	err := shim.Start(new(CoreChainCode))
	if err != nil {
		fmt.Printf("Error starting UR chaincode: %s", err)
	}
}

func (t *CoreChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Core ChainCode Initiated")

	_, args := stub.GetFunctionAndParameters()
	fmt.Printf("Init: %v", args)
	if len(args[0]) <= 0 {
		return shim.Error("MSP Mapping information is required for initiating the chain code")
	}

	var MSPListUnmarshaled []MSPList
	err := json.Unmarshal([]byte(args[0]), &MSPListUnmarshaled)

	if err != nil {
		return shim.Error("An error occurred while Unmarshiling MSPMapping: " + err.Error())
	}
	MSPMappingJSONasBytes, err := json.Marshal(MSPListUnmarshaled)
	if err != nil {
		return shim.Error("An error occurred while Marshiling MSPMapping :" + err.Error())
	}

	_Key := "MSPMapping"
	err = stub.PutState(_Key, []byte(MSPMappingJSONasBytes))
	if err != nil {
		return shim.Error("An error occurred while inserting MSPMapping:" + err.Error())
	}
	return shim.Success(nil)
}

func (t *CoreChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	certOrgType, err := getCidData(stub)
	if err != nil {
		return shim.Error("Enrolment mspid Type invalid!!! " + err.Error())
	}
	fmt.Println("MSP:" + certOrgType)

	orgType, err := getOrgTypeByMSP(stub, string(certOrgType))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("OrgType Is => : " + orgType)

	function, _ := stub.GetFunctionAndParameters()
	fmt.Println("Invoke is running for function: " + function)

	args, errTrans := getArguments(stub)
	if errTrans != nil {
		return shim.Error(errTrans.Error())
	}
	fmt.Println("Arguments Loaded Successfully!!")

	if orgType == "Core" {
		fmt.Println("Inside the OrgType condition")

		connection := hypConnect{}
		connection.Connection = stub
		switch functionName := function; functionName {
		// case "initiatePackage":
		// 	return t.initiatePackage(connection, args, "initiatePackage")
		// case "AddDocuments":
		// 	return t.AddDocuments(connection, args, "AddDocuments")
		case "signDocument":
			return t.signDocument(connection, args, "signDocument")
		case "getSignPolicy":
			return t.getSignPolicy(connection, args, "getSignPolicy")
		case "upsertSignPolicy":
			return t.upsertSignPolicy(connection, args, "upsertSignPolicy")
		case "upsertUser":
			return t.upsertUser(connection, args, "upsertUser")
		case "getUser":
			return t.getUser(connection, args, "getUser")
		case "upsertOrganization":
			return t.upsertOrganization(connection, args, "upsertOrganization")
		case "getOrganization":
			return t.getOrganization(connection, args, "getOrganization")
		case "upsertDocumentType":
			return t.upsertDocumentType(connection, args, "upsertDocumentType")
		case "upsertDocumentTypeEx":
			return t.upsertDocumentTypeEx(connection, args, "upsertDocumentType")
		case "upsertPackageType":
			return t.upsertPackageType(connection, args, "upsertPackageType")
		case "getPackageType":
			return t.getPackageType(connection, args, "getPackageType")
		case "upsertDocument":
			return t.upsertDocument(connection, args, "upsertDocument")
		case "upsertPackage":
			return t.upsertPackage(connection, args, "upsertPackage")
		case "regeneratePkg":
			return t.regeneratePkg(connection, args, "regeneratePkg")
		case "updatePackagePhysicalStatus":
			return t.updatePackagePhysicalStatus(connection, args, "upsertPackage")
		case "rejectPackage":
			return t.rejectPackage(connection, args, "rejectPackage")
		case "getDocument":
			return t.getDocument(connection, args, "getDocument")
		case "updateDocument":
			return t.updateDocument(connection, args, "updateDocument")
		case "getPackage":
			return t.getPackage(connection, args, "getPackage")
		case "getPackage2":
			return t.getPackage2(connection, args, "getPackage2")
		case "verifyDocument":
			return t.verifyDocument(connection, args, "verifyDocument")
		case "verifyDocument2":
			return t.verifyDocument2(connection, args, "verifyDocument2")
		case "addUpdateUser":
			return t.addUpdateUser(connection, args, "addUpdateUser")
		case "getOrgTypeGroups":
			return t.getOrgTypeGroups(connection, args, "getOrgTypeGroups")
		case "deleteDataFromCollection":
			return t.deleteDataFromCollection(connection, args, "deleteDataFromCollection")
		case "upsertTask":
			return t.upsertTask(connection, args, "upsertTask")
		case "upsertTaskWithoutPackage":
			return t.upsertTaskWithoutPackage(connection, args, "upsertTaskWithoutPackage")
		case "completePendingTasks":
			return t.completePendingTasks(connection, args, "completePendingTasks")
		case "continuePackage":
			return t.continuePackage(connection, args, "continuePackage")

		default:
			//logger.Warning("Invoke did not find function: " + function)
			return shim.Error("Received unknown function invocation: " + function)
		}

	} else {
		return shim.Error("Invalid MSP: " + orgType)
	}
}

func (t *CoreChainCode) continuePackage(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Println("continuePackage >>>>>>>>>>>>>>> ", args)

	var packageNo = sanitize(args[0], "string").(string)
	var options []string
	packageBytes, errorFetch := fetchData(stub, packageNo, "package")
	if errorFetch != nil {
		fmt.Println("Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
	} else if packageBytes == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", packageNo)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
	}

	var packageData Package
	errorUnmarshal := json.Unmarshal(packageBytes, &packageData)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
	} else if packageData.Status == "Completed" {
		fmt.Println("Package already Completed")
		options = append(options, "Package", packageNo)
		errBytes, _ := prepareErrorCode("2032", options)
		return shim.Error(string(errBytes))
	} else if !packageData.IsDispatch {
		fmt.Println("Package not in Dispatch state")
		options = append(options, "Package", packageNo)
		errBytes, _ := prepareErrorCode("2028", options)
		return shim.Error(string(errBytes))
	}

	packageData.IsDispatch = false
	packageMarshalled, errorMarshal := json.Marshal(packageData)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
	}

	errorInsert := insertData(&stub, packageData.Key, "package", packageMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
	}

	fmt.Println("continuePackage Successfully!")
	return shim.Success(nil)
}

func (t *CoreChainCode) completePendingTasks(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Println("completePendingTasks >>>>>>>>>>>>>>> ", args)
	var options []string
	var taskIds []string

	RequestTime := sanitize(args[1], "string").(string)

	errorUnmarshal := json.Unmarshal([]byte(args[0]), &taskIds)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
	}

	for _, taskId := range taskIds {
		TaskOBJ := Tasks{}

		TaskData, errorFetch := fetchData(stub, taskId, "tasks")
		if errorFetch != nil {
			fmt.Println("\n Task fetch error>>", errorFetch)
			options = append(options, "Task")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))

		} else if TaskData == nil {
			fmt.Println("Task not found")
			options = append(options, "Task", taskId)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
		}
		errorUnmarshal := json.Unmarshal(TaskData, &TaskOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
		}

		TaskOBJ.Status = "Completed"
		TaskOBJ.CompletedAt = RequestTime

		taskMarshalled, errorMarshal := json.Marshal(TaskOBJ)
		if errorMarshal != nil {
			fmt.Println("Error occurred while marshalling the data")
			options = append(options)
			errBytes, _ := prepareErrorCode("2003", options)
			return shim.Error(string(errBytes))
		}
		fmt.Println("Marshalled task Array: ", taskMarshalled)

		errorInsert := insertData(&stub, TaskOBJ.Key, "tasks", taskMarshalled)
		//raise shim error message if insertion fails.
		if errorInsert != nil {
			fmt.Println("Error occurred while inserting the document in private collection")
			options = append(options)
			errBytes, _ := prepareErrorCode("2004", options)
			return shim.Error(string(errBytes))
		}

	}

	return shim.Success(nil)
}

func (t *CoreChainCode) upsertTask(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Println("UPSERT TASK >>>>>>>>>>>>>>> ", args)
	var options []string
	TaskList := []TaskReqPayload{}

	errorUnmarshal := json.Unmarshal([]byte(args[0]), &TaskList)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	RequestTime := sanitize(args[1], "string").(string)
	for i := 0; i < len(TaskList); i++ {
		fmt.Println("object value >>>", TaskList[i])

		// Check for package exist...
		PackageData, errorFetch := fetchData(stub, TaskList[i].PackageId, "package")
		if errorFetch != nil {
			fmt.Println("\n Package fetch error>>", errorFetch)
			options = append(options, "Package")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Package fetch error")

		} else if PackageData == nil {
			fmt.Println("Package not found")
			options = append(options, "Package", TaskList[i].PackageId)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Package not found")
		}
		UUID := TaskList[i].PackageId + "_" + TaskList[i].TaskName + "_" + TaskList[i].OrgCode + "_" + TaskList[i].GroupName

		taskError := addTask(stub, UUID, TaskList[i].TaskName, TaskList[i].DocumentType, TaskList[i].OrgCode, TaskList[i].GroupName, TaskList[i].Port, TaskList[i].Status, TaskList[i].PackageId, TaskList[i].NotificationId, TaskList[i].SNNumber, TaskList[i].TaskStatus, TaskList[i].SLATime, TaskList[i].AdditonalData, RequestTime, TaskList[i].Stage, TaskList[i].Activity)
		if taskError != "" {
			return shim.Error(taskError)
		}
	}
	return shim.Success(nil)
}

func (t *CoreChainCode) upsertTaskWithoutPackage(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Println("UPSERT TASK >>>>>>>>>>>>>>> ", args)

	TaskList := []TaskReqPayload{}
	var options []string
	errorUnmarshal := json.Unmarshal([]byte(args[0]), &TaskList)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	RequestTime := sanitize(args[1], "string").(string)
	for i := 0; i < len(TaskList); i++ {
		fmt.Println("object value >>>", TaskList[i])

		// Check for package exist...

		UUID := TaskList[i].PackageId + "_" + TaskList[i].TaskName + "_" + TaskList[i].OrgCode + "_" + TaskList[i].GroupName

		taskError := addTask(stub, UUID, TaskList[i].TaskName, TaskList[i].DocumentType, TaskList[i].OrgCode, TaskList[i].GroupName, TaskList[i].Port, TaskList[i].Status, TaskList[i].PackageId, TaskList[i].NotificationId, TaskList[i].SNNumber, TaskList[i].TaskStatus, TaskList[i].SLATime, TaskList[i].AdditonalData, RequestTime, TaskList[i].Stage, TaskList[i].Activity)
		if taskError != "" {
			return shim.Error(taskError)
		}
	}
	return shim.Success(nil)
}

func (t *CoreChainCode) getSignPolicy(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign policy: %v", args)

	Key := sanitize(args[0], "string").(string)
	var options []string
	SignPolicyOBJ := SignaturePolicy{}
	SignPolicyCheck, errorFetch := fetchData(stub, Key, "SignaturePolicy")
	fmt.Println("\n sign policy Key getSignPolicy>>", Key)
	if errorFetch != nil {
		fmt.Println("\n SignPolicy fetch error>>", errorFetch)
		options = append(options, "SignPolicy")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"SignPolicy fetch error"}, "90200001", false, true)
		// return shim.Error(string(errBytes))

	} else if SignPolicyCheck == nil {
		fmt.Println("SignPolicy not found")
		options = append(options, "SignPolicy", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"SignPolicy not found"}, "90200002", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("SignPolicy exist")

		errorUnmarshal := json.Unmarshal(SignPolicyCheck, &SignPolicyOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nSign Policy >>>>>", SignPolicyOBJ)

	SignPolicyMarshalled, errorMarshal := json.Marshal(SignPolicyOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("SignPolicyMarshalled marshalled successful !!!")

	return shim.Success(SignPolicyMarshalled)
}

func (t *CoreChainCode) upsertSignPolicy(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign policy: %v", args)

	Key := sanitize(args[0], "string").(string)
	fmt.Print("Key >>>", Key)

	SignPolicyOBJ := SignaturePolicy{}
	var options []string
	SignPolicyData, errorFetch := fetchData(stub, Key, "SignaturePolicy")
	fmt.Println("\n sign policy Key upsertSignPolicy>>", Key)
	if errorFetch != nil {
		fmt.Println("\n SignPolicy fetch error>>", errorFetch)
		options = append(options, "SignPolicy")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"SignPolicy fetch error"}, "90200001", false, true)
		// return shim.Error(string(errBytes))

	} else if SignPolicyData != nil {
		fmt.Println("SignPolicy exist")
		errorUnmarshal := json.Unmarshal(SignPolicyData, &SignPolicyOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	//else {
	// 	fmt.Println("SignPolicy not found")
	// }
	// fmt.Println("\n\n\nSign Policy >>>>>", SignPolicyOBJ)

	SignPolicyOBJ.Key = Key
	SignPolicyOBJ.PolicyName = Key
	SignPolicyOBJ.DocumentName = "signaturepolicy"
	SignPolicyOBJ.Sequence = sanitize(args[2], "bool").(bool)

	err := json.Unmarshal([]byte(args[1]), &SignPolicyOBJ.Organizations)
	if err != nil {
		fmt.Printf("\n Error occurred while unmarshalling Policy Organizations")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}

	//upsert
	SignPolicyMarshalled, errorMarshal := json.Marshal(SignPolicyOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "SignaturePolicy", SignPolicyMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in SignaturePolicy Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) upsertUser(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("user: %v", args)

	UserID := sanitize(args[0], "string").(string)
	Email := sanitize(args[1], "string").(string)
	key := hash(UserID)
	fmt.Print("Key >>>", key)

	Signature := sanitize(args[2], "string").(string)
	PublicCertificate := sanitize(args[3], "string").(string)
	P12Certificate := sanitize(args[4], "string").(string)
	DateTime := sanitize(args[5], "string").(string)
	Password := sanitize(args[6], "string").(string)
	OrgCode := sanitize(args[7], "string").(string)
	CertificateType := sanitize(args[8], "string").(string)

	UserOBJ := User{}
	var options []string
	UserData, errorFetch := fetchData(stub, key, "user")
	if errorFetch != nil {
		fmt.Println("\n User fetch error>>", errorFetch)
		options = append(options, "User")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"User fetch error"}, "90300001", false, true)
		// return shim.Error(string(errBytes))

	} else if UserData != nil {
		fmt.Println("User already exists")

		errorUnmarshal := json.Unmarshal(UserData, &UserOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}

		// UserOBJ.History = append(UserOBJ.History, UserOBJ.Current)

	}

	fmt.Println("\n\n\n UserOBJ Current 0 >>>>", UserOBJ.Current)

	UserOBJ.Key = key
	UserOBJ.DocumentName = "user"
	UserOBJ.UserId = UserID
	UserOBJ.Email = Email
	UserOBJ.OrgCode = OrgCode

	updateType := ""

	if P12Certificate != "" && Signature != "" {
		updateType = "Both"
	} else if Signature != "" {
		updateType = "Signature"
	} else if P12Certificate != "" {
		updateType = "Certificate"
	}

	fmt.Println("\n\n\n update Type >>>>", updateType)

	UserOBJ.Current.Datetime = DateTime
	UserOBJ.Current.Type = updateType
	fmt.Println("\n\n\n UserOBJ Current 1 >>>>", UserOBJ.Current)

	if updateType == "Signature" {
		if UserOBJ.Current.Signature == Signature {
			options = append(options)
			errBytes, _ := prepareErrorCode("2022", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Signature is already updated"}, "90400001", false, true)
			// return shim.Error(string(errBytes))
		} else {
			UserOBJ.Current.Signature = Signature
		}
	}

	if updateType == "Certificate" {
		if UserOBJ.Current.P12Certificate == P12Certificate {
			options = append(options)
			errBytes, _ := prepareErrorCode("2027", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Certificate is already updated"}, "90400002", false, true)
			// return shim.Error(string(errBytes))
		} else {
			UserOBJ.Current.P12Certificate = P12Certificate
			UserOBJ.Current.PublicCertificate = PublicCertificate
			UserOBJ.Current.Password = Password
			UserOBJ.Current.CertificateType = CertificateType
		}
	}

	fmt.Println("\n\n\n p12Certificate Current", UserOBJ.Current.P12Certificate)
	fmt.Println("\n\n\n p12Certificate Old", P12Certificate)
	fmt.Println("\n\n\n signature Current", UserOBJ.Current.Signature)
	fmt.Println("\n\n\n signature Current", Signature)

	if updateType == "Both" {
		if UserOBJ.Current.P12Certificate == P12Certificate && UserOBJ.Current.Signature == Signature {
			options = append(options)
			errBytes, _ := prepareErrorCode("2037", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Certificate and Signature already updated"}, "90400003", false, true)
			// return shim.Error(string(errBytes))
		} else {
			UserOBJ.Current.Signature = Signature
			UserOBJ.Current.P12Certificate = P12Certificate
			UserOBJ.Current.PublicCertificate = PublicCertificate
			UserOBJ.Current.Password = Password
			UserOBJ.Current.CertificateType = CertificateType
		}
	}

	// if UserOBJ.History == nil {
	UserOBJ.History = append(UserOBJ.History, UserOBJ.Current)
	// }

	//insert

	UserMarshalled, errorMarshal := json.Marshal(UserOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, key, "user", UserMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in User Collection!")
	return shim.Success(nil)
}

func (t *CoreChainCode) deleteDataFromCollection(stub hypConnect, args []string, functionName string) pb.Response {

	err := deleteDataFromCollection(&stub, args, functionName)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *CoreChainCode) getUser(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf(" get user: %v", args)

	UserID := sanitize(args[0], "string").(string)
	Key := UserID
	var options []string
	UserOBJ := User{}
	UserCheck, errorFetch := fetchData(stub, Key, "user")
	if errorFetch != nil {
		fmt.Println("\n User fetch error>>", errorFetch)
		options = append(options, "User")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"User fetch error"}, "90300001", false, true)
		// return shim.Error(string(errBytes))

	} else if UserCheck == nil {
		fmt.Println("User not found")
		options = append(options, "User", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"User not found"}, "90300002", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("User exist")

		errorUnmarshal := json.Unmarshal(UserCheck, &UserOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nUser Object >>>>>", UserOBJ)

	UserMarshalled, errorMarshal := json.Marshal(UserOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("UserMarshalled marshalled successful !!!")

	return shim.Success(UserMarshalled)
}

func (t *CoreChainCode) upsertOrganization(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("org: %v", args)

	OrgCode := sanitize(args[0], "string").(string)
	key := OrgCode
	fmt.Print("Key >>>", key)

	Stamp := sanitize(args[1], "string").(string)
	PublicCertificate := sanitize(args[2], "string").(string)
	P12Certificate := sanitize(args[3], "string").(string)
	DateTime := sanitize(args[4], "string").(string)
	Password := sanitize(args[5], "string").(string)
	OrgName := sanitize(args[6], "string").(string)
	CertificateType := sanitize(args[7], "string").(string)

	OrganizationOBJ := Organization{}
	var options []string
	OrganizationData, errorFetch := fetchData(stub, key, "organization")
	if errorFetch != nil {
		fmt.Println("\n Organization fetch error>>", errorFetch)
		options = append(options, "Organization")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Organization fetch error"}, "90500001", false, true)
		// return shim.Error(string(errBytes))

	} else if OrganizationData != nil {
		fmt.Println("Organization already exists")

		errorUnmarshal := json.Unmarshal(OrganizationData, &OrganizationOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}

		// OrganizationOBJ.History = append(OrganizationOBJ.History, OrganizationOBJ.Current)

	}

	OrganizationOBJ.Key = key
	OrganizationOBJ.DocumentName = "organization"
	OrganizationOBJ.OrgName = OrgName
	OrganizationOBJ.OrgCode = OrgCode

	errGroups := json.Unmarshal([]byte(args[8]), &OrganizationOBJ.GroupIDs)
	fmt.Println("Organization GroupIds >>>", OrganizationOBJ.GroupIDs)

	if errGroups != nil {
		fmt.Printf("\n Error occurred while unmarshalling Groups")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}

	updateType := ""
	if P12Certificate != "" && Stamp != "" {
		updateType = "Both"
	} else if Stamp != "" {
		updateType = "Stamp"
	} else if P12Certificate != "" {
		updateType = "Certificate"
	}

	fmt.Println("UpdateType >>>", updateType)

	OrganizationOBJ.Current.Datetime = DateTime
	OrganizationOBJ.Current.Type = updateType

	fmt.Println("OrganizationOBJ >>>", OrganizationOBJ)

	if updateType == "Stamp" {
		if OrganizationOBJ.Current.Stamp == Stamp {
			options = append(options)
			errBytes, _ := prepareErrorCode("2038", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Stamp is already updated"}, "90400004", false, true)
			// return shim.Error(string(errBytes))
		} else {
			OrganizationOBJ.Current.Stamp = Stamp
		}
	}

	if updateType == "Certificate" {
		if OrganizationOBJ.Current.P12Certificate == P12Certificate {
			options = append(options)
			errBytes, _ := prepareErrorCode("2027", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Certificate is already updated"}, "90400002", false, true)
			// return shim.Error(string(errBytes))
		} else {
			OrganizationOBJ.Current.P12Certificate = P12Certificate
			OrganizationOBJ.Current.PublicCertificate = PublicCertificate
			OrganizationOBJ.Current.Password = Password
			OrganizationOBJ.Current.CertificateType = CertificateType
		}
	}

	if updateType == "Both" {
		if OrganizationOBJ.Current.P12Certificate == P12Certificate && OrganizationOBJ.Current.Stamp == Stamp {
			options = append(options)
			errBytes, _ := prepareErrorCode("2039", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Certificate and Stamp already updated"}, "90400005", false, true)
			// return shim.Error(string(errBytes))
		} else {
			OrganizationOBJ.Current.Stamp = Stamp
			OrganizationOBJ.Current.P12Certificate = P12Certificate
			OrganizationOBJ.Current.PublicCertificate = PublicCertificate
			OrganizationOBJ.Current.Password = Password
			OrganizationOBJ.Current.CertificateType = CertificateType

		}
	}

	// if OrganizationOBJ.History == nil {
	OrganizationOBJ.History = append(OrganizationOBJ.History, OrganizationOBJ.Current)
	// }

	//insert

	OrganizationMarshalled, errorMarshal := json.Marshal(OrganizationOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, key, "organization", OrganizationMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in Organization Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) getOrganization(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf(" get org: %v", args)

	OrgCode := sanitize(args[0], "string").(string)
	var options []string
	OrganizationOBJ := Organization{}
	OrganizationCheck, errorFetch := fetchData(stub, OrgCode, "organization")
	if errorFetch != nil {
		fmt.Println("\n Organization fetch error>>", errorFetch)
		options = append(options, "Organization")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Organization fetch error"}, "90500001", false, true)
		// return shim.Error(string(errBytes))

	} else if OrganizationCheck == nil {
		fmt.Println("Organization not found")
		options = append(options, "Organization", OrgCode)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Organization not found"}, "90500002", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("Organization exist")

		errorUnmarshal := json.Unmarshal(OrganizationCheck, &OrganizationOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nOrganization Object >>>>>", OrganizationOBJ)

	OrganizationMarshalled, errorMarshal := json.Marshal(OrganizationOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("OrganizationMarshalled marshalled successful !!!")

	return shim.Success(OrganizationMarshalled)
}

func (t *CoreChainCode) upsertDocumentType(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("document type: %v", args)

	Key := sanitize(args[0], "string").(string)
	fmt.Print("Key >>>", Key)

	DocumentTypeOBJ := DocumentType{}
	var options []string
	DocumentTypeData, errorFetch := fetchData(stub, Key, "documentType")
	if errorFetch != nil {
		fmt.Println("\n DocumentType fetch error>>", errorFetch)
		options = append(options, "DocumentType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"DocumentType fetch error"}, "90600001", false, true)
		// return shim.Error(string(errBytes))

	} else if DocumentTypeData != nil {
		fmt.Println("DocumentType already exists")
		errorUnmarshal := json.Unmarshal(DocumentTypeData, &DocumentTypeOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}

	fmt.Println("\n\n\nDocument Type >>>>>", DocumentTypeOBJ)

	DocumentTypeOBJ.Key = Key
	DocumentTypeOBJ.DocumentName = "documenttype"
	DocumentTypeOBJ.DocumentType = Key
	DocumentTypeOBJ.Label = sanitize(args[1], "string").(string)
	DocumentTypeOBJ.Value = sanitize(args[2], "string").(string)

	//upsert
	DocumentTypeMarshalled, errorMarshal := json.Marshal(DocumentTypeOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "documentType", DocumentTypeMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in DocumentType Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) upsertDocumentTypeEx(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("document type: %v", args)
	var options []string
	SignPolicyOBJ := SignaturePolicy{}
	errorUnmarshal := json.Unmarshal([]byte(args[1]), &SignPolicyOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	Key := SignPolicyOBJ.PolicyName
	fmt.Print("Key >>>", Key)

	SignPolicyData, errorFetch := fetchData(stub, Key, "SignaturePolicy")
	fmt.Println("\n sign policy Key upsertSignPolicy>>", Key)
	if errorFetch != nil {
		fmt.Println("\n SignPolicy fetch error>>", errorFetch)
		options = append(options, "SignPolicy")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"SignPolicy fetch error"}, "90200001", false, true)
		// return shim.Error(string(errBytes))

	} else if SignPolicyData != nil {
		fmt.Println("SignPolicy exist")

	}
	//else {
	// 	fmt.Println("SignPolicy not found")
	// }
	// fmt.Println("\n\n\nSign Policy >>>>>", SignPolicyOBJ)

	SignPolicyOBJ.Key = Key
	SignPolicyOBJ.PolicyName = Key
	SignPolicyOBJ.DocumentName = "signaturepolicy"

	// err := json.Unmarshal([]byte(SignPolicyOBJ.), &SignPolicyOBJ.Organizations)
	// if err != nil {
	// 	fmt.Printf("\n Error occurred while unmarshalling Policy Organizations")
	// 	errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
	// 	return shim.Error(string(errBytes))

	// }

	//upsert
	SignPolicyMarshalled, errorMarshal := json.Marshal(SignPolicyOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "SignaturePolicy", SignPolicyMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in SignaturePolicy Collection!")

	DocumentTypeOBJ := DocumentType{}
	errorUnmarshal = json.Unmarshal([]byte(args[0]), &DocumentTypeOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	Key = DocumentTypeOBJ.DocumentType
	fmt.Print("Key >>>", Key)

	DocumentTypeData, errorFetch := fetchData(stub, Key, "documentType")
	if errorFetch != nil {
		fmt.Println("\n DocumentType fetch error>>", errorFetch)
		options = append(options, "DocumentType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"DocumentType fetch error"}, "90600001", false, true)
		// return shim.Error(string(errBytes))

	} else if DocumentTypeData != nil {
		fmt.Println("DocumentType already exists")
		// errorUnmarshal := json.Unmarshal(DocumentTypeData, &DocumentTypeOBJ)
		// if errorUnmarshal != nil {
		// 	errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// 	return shim.Error(string(errBytes))
		// }
	}

	fmt.Println("\n\n\nDocument Type >>>>>", DocumentTypeOBJ)

	DocumentTypeOBJ.Key = Key
	DocumentTypeOBJ.DocumentName = "documenttype"
	DocumentTypeOBJ.SignaturePolicyId = SignPolicyOBJ.Key

	//upsert
	DocumentTypeMarshalled, errorMarshal := json.Marshal(DocumentTypeOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert = insertData(&stub, Key, "documentType", DocumentTypeMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in DocumentType Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) upsertPackageType(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("package type: %v", args)

	Key := sanitize(args[0], "string").(string)
	fmt.Print("Key >>>", Key)
	var options []string

	PackageTypeOBJ := PackageType{}

	PackageTypeData, errorFetch := fetchData(stub, Key, "packageType")
	if errorFetch != nil {
		fmt.Println("\n PackageType fetch error>>", errorFetch)
		options = append(options, "PackageType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType fetch error"}, "90700001", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageTypeData != nil {
		fmt.Println("DocumentType already exists")
		errorUnmarshal := json.Unmarshal(PackageTypeData, &PackageTypeOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nPackage Type >>>>>", PackageTypeOBJ)

	PackageTypeOBJ.Key = Key
	PackageTypeOBJ.PackageType = Key
	PackageTypeOBJ.DocumentName = "packagetype"

	errSign := json.Unmarshal([]byte(args[1]), &PackageTypeOBJ.Documents)
	if errSign != nil {
		fmt.Printf("\n Error occurred while unmarshalling package documents")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}

	err := json.Unmarshal([]byte(args[2]), &PackageTypeOBJ.NotifyPolicy)
	if err != nil {
		fmt.Printf("\n Error occurred while unmarshalling package notify parties")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	for i := 0; i < len(PackageTypeOBJ.Documents); i++ {
		fmt.Println("array value >>>", PackageTypeOBJ.Documents[i])
		DocumentCheck, errorFetch := fetchData(stub, PackageTypeOBJ.Documents[i].DocumentType, "documentType")
		if errorFetch != nil {
			fmt.Println("\n DocumentType fetch error>>", errorFetch)
			options = append(options, "DocumentType")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"DocumentType fetch error"}, "90600001", false, true)
			// return shim.Error(string(errBytes))

		} else if DocumentCheck == nil {
			options = append(options, "DocumentType", PackageTypeOBJ.Documents[i].DocumentType)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Document Type not found"}, "90600002", false, true)
			// return shim.Error(string(errBytes))
		}

		SignaturePolCheck, errorFetch := fetchData(stub, PackageTypeOBJ.Documents[i].SignaturePolicy, "SignaturePolicy")
		if errorFetch != nil {
			fmt.Println("\n SignaturePolicy fetch error>>", errorFetch)
			options = append(options, "SignaturePolicy")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"SignaturePolicy fetch error"}, "90200001", false, true)
			// return shim.Error(string(errBytes))

		} else if SignaturePolCheck == nil {
			options = append(options, "Signature Policy", PackageTypeOBJ.Documents[i].SignaturePolicy)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Signature Policy not found"}, "90200002", false, true)
			// return shim.Error(string(errBytes))
		}
	}

	//upsert
	PackageTypeMarshalled, errorMarshal := json.Marshal(PackageTypeOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "packageType", PackageTypeMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in PackageType Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) getPackageType(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign policy: %v", args)

	Key := sanitize(args[0], "string").(string)
	var options []string
	PackageTypeOBJ := PackageType{}
	PackageTypeCheck, errorFetch := fetchData(stub, Key, "packageType")
	if errorFetch != nil {
		fmt.Println("\n PackageType fetch error>>", errorFetch)
		options = append(options, "PackageType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType fetch error"}, "90700001", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageTypeCheck == nil {
		fmt.Println("PackageType not found")
		options = append(options, "PackageType", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType not found"}, "90700002", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("PackageType exist")

		errorUnmarshal := json.Unmarshal(PackageTypeCheck, &PackageTypeOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nPkg Type >>>>>", PackageTypeOBJ)

	PackageTypeMarshalled, errorMarshal := json.Marshal(PackageTypeOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("PackageTypeMarshalled marshalled successful !!!")

	return shim.Success(PackageTypeMarshalled)
}

func (t *CoreChainCode) updatePackagePhysicalStatus(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("package type: %v", args)

	Key := sanitize(args[0], "string").(string)
	status := sanitize(args[1], "string").(string)

	fmt.Print("Key >>>", Key)

	PackageOBJ := Package{}
	var options []string
	//Package Check Logic
	PackageData, errorFetch := fetchData(stub, Key, "package")
	if errorFetch != nil {
		fmt.Println("\n PackageType fetch error>>", errorFetch)
		options = append(options, "PackageType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType fetch error"}, "90700001", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageData == nil {
		options = append(options, "PackageType", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType not found"}, "90700002", false, true)
		// return shim.Error(string(errBytes))
	}
	errorUnmarshal := json.Unmarshal(PackageData, &PackageOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("PackageType exist >>>>>>>>>>>>>>.", PackageOBJ)

	if PackageOBJ.Status == "Completed" {
		// if isDownload == true {
		// 	PackageOBJ.PhysicalStatus = "Downloaded"
		// } else {
		// 	PackageOBJ.PhysicalStatus = "Printed"
		// }
		if strings.ToLower(status) == strings.ToLower("Printed") {

			PackageOBJ.PhysicalStatus = "Printed"
		} else if strings.ToLower(status) == strings.ToLower("Uploaded") {

			if PackageOBJ.PhysicalStatus == "Printed" {
				PackageOBJ.PhysicalStatus = "Uploaded"
			} else {
				options = append(options)
				errBytes, _ := prepareErrorCode("2007", options)
				return shim.Error(string(errBytes))
				//return shim.Error("PhysicalStatus is not Printed")
			}

		} else if strings.ToLower(status) == strings.ToLower("Forwarded") {

			if PackageOBJ.PhysicalStatus == "Uploaded" {
				PackageOBJ.PhysicalStatus = "Forwarded"

				// send notification event
			} else {
				options = append(options)
				errBytes, _ := prepareErrorCode("2008", options)
				return shim.Error(string(errBytes))
				//return shim.Error("PhysicalStatus is not Uploaded")
			}

		} else {
			options = append(options)
			errBytes, _ := prepareErrorCode("2023", options)
			return shim.Error(string(errBytes))
			//	return shim.Error("Invalid Status")

		}
	} else {
		options = append(options)
		errBytes, _ := prepareErrorCode("2028", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Package is In Progress")
	}
	fmt.Println("PackageType updated payload >>>>>>>>>>>>>>.", PackageOBJ)

	//upsert
	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "package", PackageMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully PhysicalStatus in Package Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) regeneratePkg(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Println("packagessssssss regeneratePkg:", args, len(args))

	tnxId := stub.Connection.GetTxID()
	Key := sanitize(args[0], "string").(string)
	var options []string

	PackageOBJ := Package{}
	previousDocs := []UploadedDocuments{}
	fmt.Println("\n\n\nPackage Args >>>>>", PackageOBJ, previousDocs)

	//Uploaded Document Type Check
	var UploadedDocumentsOBJ []UploadedDocuments

	err := json.Unmarshal([]byte(args[1]), &UploadedDocumentsOBJ)
	if err != nil {
		fmt.Printf("\n Error occurred while unmarshalling StampOrg")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("\n\n\nPackage Documents from RequestPayload >>>>>", UploadedDocumentsOBJ)
	fmt.Println("------------------------------------------------------------------------------")

	PackageData, errorFetch := fetchData(stub, Key, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
		// return shim.Error(string(errBytes))
	} else if PackageData == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
		// return shim.Error(string(errBytes))
	}

	errorUnmarshal := json.Unmarshal(PackageData, &PackageOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("\n\n\nPackage >>>>>", PackageOBJ)

	previousDocs = PackageOBJ.Documents
	fmt.Println("Updated Documents!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", PackageOBJ.Documents)

	//upsert
	fmt.Println("PackageOBJ >>>", PackageOBJ)

	for j := 0; j < len(UploadedDocumentsOBJ); j++ {

		UOBJ := UploadedDocumentsOBJ[j]

		// Previous Docs for maintaining DocHistory
		prevUploadDocOBJ := UploadedDocuments{}
		prevDocOBJ := Document{}
		for p := 0; p < len(previousDocs); p++ {
			if previousDocs[p].DocumentType == UOBJ.DocumentType {
				prevUploadDocOBJ = previousDocs[p]
			}
		}

		//get Previous Documents
		if len(prevUploadDocOBJ.Hash) > 0 {
			PrevDocData, errorFetch := fetchData(stub, PackageOBJ.Key+"-"+prevUploadDocOBJ.Hash, "document")
			fmt.Println("\n PrevDocData upsertPackage3>>")
			if errorFetch != nil {
				fmt.Println("\n PrevDoc fetch error upsertPackage>>", errorFetch)
				options = append(options, "PrevDoc")
				errBytes, _ := prepareErrorCode("2001", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"PrevDoc fetch error"}, "90200001", false, true)
				// return shim.Error(string(errBytes))
			} else if PrevDocData == nil {
				fmt.Println("PrevDoc not found")
				options = append(options, "PrevDoc", PackageOBJ.Key+"-"+prevUploadDocOBJ.Hash)
				errBytes, _ := prepareErrorCode("2002", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"PrevDoc not found"}, "90200002", false, true)
				// return shim.Error(string(errBytes))
			}

			errorUnmarshal := json.Unmarshal(PrevDocData, &prevDocOBJ)
			if errorUnmarshal != nil {
				options = append(options)
				errBytes, _ := prepareErrorCode("2000", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
				// return shim.Error(string(errBytes))
			}
			fmt.Println("\n\n\nPREVIOUS DOCUMENT OBJ >>>>>", prevDocOBJ)

			prevDocOBJ.History = append(prevDocOBJ.History, prevDocOBJ.CurrentVersion)
			prevDocOBJ.Version = prevDocOBJ.Version + 1
			prevDocOBJ.GUID = UOBJ.GUID
			prevDocOBJ.Status = "In Progress"
			prevDocOBJ.Progress = 0
			resetSign := []Signatures{}
			for rs := 0; rs < len(prevDocOBJ.Signatures); rs++ {
				prevDocOBJ.Signatures[rs].Action = "Pending"
				resetSign = append(resetSign, prevDocOBJ.Signatures[rs])
			}
			prevDocOBJ.Signatures = resetSign
			fmt.Println("RresetSign >>> @@@@@@@@@@@@@@@@@@@@@@@@@@@@", prevDocOBJ.Signatures, resetSign)

			prevDocOBJ.CurrentVersion = History{
				Type:       "Uploaded",
				OrgCode:    PackageOBJ.RequestedBy.OrgCode,
				OrgName:    PackageOBJ.RequestedBy.OrgName,
				UserId:     PackageOBJ.RequestedBy.UserId,
				UserName:   PackageOBJ.RequestedBy.UserName,
				Datetime:   PackageOBJ.RequestedOn,
				Hash:       UOBJ.Hash,
				SequenceNo: 0,
				Status:     "Uploaded",
				TranxHash:  tnxId,
			}

			if len(prevDocOBJ.Signatures) > 0 {
				smallest := prevDocOBJ.Signatures[0]            // set the smallest number to the first element of the list
				for _, num := range prevDocOBJ.Signatures[1:] { // iterate over the rest of the list
					if num.SequenceNo < smallest.SequenceNo { // if num is smaller than the current smallest number
						smallest = num // set smallest to num
					}
				}
				fmt.Println("\n\n\nPackage >>>>>", smallest)

				prevDocOBJ.NextSignatory = smallest.OrgCode
			}

			fmt.Println("Before DocumentMarshalled >>>", prevDocOBJ)
			//upsert document
			DocumentMarshalled, errorMarshal := json.Marshal(prevDocOBJ)
			if errorMarshal != nil {
				fmt.Println("Error occurred while marshalling the data")
				options = append(options)
				errBytes, _ := prepareErrorCode("2003", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
				// return shim.Error(string(errBytes))
			}
			errorInsert := insertData(&stub, prevDocOBJ.Key, "document", DocumentMarshalled)
			if errorInsert != nil {
				fmt.Println("Error occurred while inserting the Document in private collection")
				options = append(options)
				errBytes, _ := prepareErrorCode("2004", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
				// return shim.Error(string(errBytes))
			}
			fmt.Println("Upsert Successfully in Document Collection!")

		}
	}

	PackageOBJ.IsRegenrated = true

	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, PackageOBJ.Key, "package", PackageMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("Upsert Successfully in package Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) upsertPackage(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("packagessssssss:", args, len(args))

	tnxId := stub.Connection.GetTxID()

	Key := sanitize(args[0], "string").(string)

	packageMode := sanitize(args[7], "string").(string)

	packageType := sanitize(args[1], "string").(string)

	//collection Name
	var DocumentMappingCollection = "documentMapping"
	var DocumentMappingStruct DocumentMapping

	PackageOBJ := Package{}
	previousDocs := []UploadedDocuments{}
	fmt.Println("\n\n\nPackage BEF >>>>>", PackageOBJ, previousDocs)
	var options []string
	//Package Type Check Logic
	PackageTypeOBJ := PackageType{}
	PackageTypeData, errorFetch := fetchData(stub, packageType, "packageType")
	if errorFetch != nil {
		fmt.Println("\n PackageType fetch error>>", errorFetch)
		options = append(options, "PackageType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType fetch error"}, "90700001", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageTypeData == nil {
		options = append(options, "PackageType", packageType)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType not found"}, "90700002", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("PackageType exist")

	errorUnmarshal := json.Unmarshal(PackageTypeData, &PackageTypeOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	//Uploaded Document Type Check
	var UploadedDocumentsOBJ []UploadedDocuments
	var OverRideOrg []OverRide

	errorUnmarshal = json.Unmarshal([]byte(args[10]), &OverRideOrg)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	errorUnmarshal = json.Unmarshal([]byte(args[10]), &PackageOBJ.OverrideOrgs)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error PackageOBJ.OverrideOrgs"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	errorUnmarshal = json.Unmarshal([]byte(args[14]), &PackageOBJ.OtherDocuments)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("\n\n\nOverRide ORGS 1 >>>>>", OverRideOrg)

	err := json.Unmarshal([]byte(args[9]), &UploadedDocumentsOBJ)
	if err != nil {
		fmt.Printf("\n Error occurred while unmarshalling StampOrg")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("\n\n\nPackage BEF 1 >>>>>", UploadedDocumentsOBJ)

	for i := 0; i < len(PackageTypeOBJ.Documents); i++ {
		found := false
		for j := 0; j < len(UploadedDocumentsOBJ); j++ {
			if PackageTypeOBJ.Documents[i].DocumentType == UploadedDocumentsOBJ[j].DocumentType {
				SignPolicyOBJ := SignaturePolicy{}
				SignPolicyData, errorFetch := fetchData(stub, PackageTypeOBJ.Documents[i].SignaturePolicy, "SignaturePolicy")
				fmt.Println("\n PackageTypeOBJ.Documents[i].SignaturePolicy upsertPackage>>", PackageTypeOBJ.Documents[i].SignaturePolicy)
				if errorFetch != nil {
					fmt.Println("\n SignPolicy fetch error>>", errorFetch)
					options = append(options, "SignPolicy")
					errBytes, _ := prepareErrorCode("2001", options)
					return shim.Error(string(errBytes))
					// errBytes := genericErrorHandler(nil, []string{"SignPolicy fetch error"}, "90200001", false, true)
					// return shim.Error(string(errBytes))

				} else if SignPolicyData == nil {
					fmt.Println("SignPolicy not found")
					options = append(options, "SignPolicy", PackageTypeOBJ.Documents[i].SignaturePolicy)
					errBytes, _ := prepareErrorCode("2002", options)
					return shim.Error(string(errBytes))
					// errBytes := genericErrorHandler(nil, []string{"SignPolicy not found"}, "90200002", false, true)
					// return shim.Error(string(errBytes))
				}

				errorUnmarshal := json.Unmarshal(SignPolicyData, &SignPolicyOBJ)
				if errorUnmarshal != nil {
					options = append(options)
					errBytes, _ := prepareErrorCode("2000", options)
					return shim.Error(string(errBytes))
					// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
					// return shim.Error(string(errBytes))
				}

				fmt.Println("\n\n\nSign Policy >>>>>", SignPolicyOBJ)
				if PackageTypeOBJ.Documents[i].SignaturePolicy == "" {
					fmt.Println("\n\n\nPackageTypeOBJ.Documents[i] ln 809", PackageTypeOBJ.Documents[i])
					fmt.Println("\n\n\nUploadedDocumentsOBJ[j] ln 810", UploadedDocumentsOBJ[j])
				}

				UploadedDocumentsOBJ[j].SignaturePolicy = PackageTypeOBJ.Documents[i].SignaturePolicy

				for o := 0; o < len(SignPolicyOBJ.Organizations); o++ {
					optionalOrgs := false
					if SignPolicyOBJ.Organizations[o].OrganizationCode != "" {
						for _, value := range UploadedDocumentsOBJ[j].OptionalOrgs {
							if value == SignPolicyOBJ.Organizations[o].OrganizationCode {
								// logic
								optionalOrgs = true
							}
						}
						//do something here
						if !optionalOrgs {
							UploadedDocumentsOBJ[j].SignatureOrgs = append(UploadedDocumentsOBJ[j].SignatureOrgs, SignPolicyOBJ.Organizations[o].OrganizationCode)
						}
					} else {
						for _, value := range OverRideOrg {
							if value.OrgType == SignPolicyOBJ.Organizations[o].OrganizationType {
								// logic
								for _, value1 := range UploadedDocumentsOBJ[j].OptionalOrgs {
									if value1 == value.OrgCode {
										// logic
										optionalOrgs = true
									}
								}
								if !optionalOrgs {
									UploadedDocumentsOBJ[j].SignatureOrgs = append(UploadedDocumentsOBJ[j].SignatureOrgs, value.OrgCode)
								}
							}
						}
						var emptyArray []string
						if len(OverRideOrg) == 0 {
							UploadedDocumentsOBJ[j].SignatureOrgs = emptyArray
						}
					}

					for g := 0; g < len(SignPolicyOBJ.Organizations[o].Group); g++ {
						groupAlreadyExist := false
						for u := 0; u < len(UploadedDocumentsOBJ[j].SignatureGroups); u++ {
							if UploadedDocumentsOBJ[j].SignatureGroups[u] == SignPolicyOBJ.Organizations[o].Group[g] {
								groupAlreadyExist = true
							}
						}
						if !groupAlreadyExist {
							UploadedDocumentsOBJ[j].SignatureGroups = append(UploadedDocumentsOBJ[j].SignatureGroups, SignPolicyOBJ.Organizations[o].Group...)
						}
					}

					if len(UploadedDocumentsOBJ[j].SignatureGroups) == 0 {
						UploadedDocumentsOBJ[j].SignatureGroups = append(UploadedDocumentsOBJ[j].SignatureGroups, SignPolicyOBJ.Organizations[o].Group...)
					}
				}

				UploadedDocumentsOBJ[j].FinalSignatoryOrganizationCode = SignPolicyOBJ.FinalSignatoryOrganizationCode
				UploadedDocumentsOBJ[j].FinalSignatoryOrganizationName = SignPolicyOBJ.FinalSignatoryOrganizationName
				UploadedDocumentsOBJ[j].FinalSignatoryOrganizationType = SignPolicyOBJ.FinalSignatoryOrganizationType
				UploadedDocumentsOBJ[j].FinalSignatoryGroup = SignPolicyOBJ.FinalSignatoryGroup

				found = true
			}
		}

		fmt.Println("------------------------------------------------------------------------------", found)
		fmt.Println("------------------------------------------------------------------------------", found)
		// if !found {

		// 	errMsg := PackageTypeOBJ.Documents[i].DocumentType + " not found in this Package Type Definition"
		// 	errBytes := genericErrorHandler(nil, []string{errMsg}, "90700003", false, true)
		// 	return shim.Error(string(errBytes))
		// }
	}

	fmt.Println("------------------------------------------------------------------------------")
	fmt.Println("------------------------------------------------------------------------------", args[10], args[11])

	if packageMode == "NEW" {
		PackageOBJ.Key = Key
		PackageOBJ.DocumentName = "package"
		PackageOBJ.PackageType = packageType
		PackageOBJ.OwnerOrg = sanitize(args[13], "string").(string)

		PackageOBJ.Refno = sanitize(args[12], "string").(string)

		PackageOBJ.SignaturePolicy = sanitize(args[11], "string").(string)
		PackageOBJ.PackageNumber = Key
		PackageOBJ.PackageName = sanitize(args[2], "string").(string)
		PackageOBJ.RequestedBy = RequestedBy{
			OrgCode:  sanitize(args[3], "string").(string),
			OrgName:  sanitize(args[4], "string").(string),
			UserId:   sanitize(args[5], "string").(string),
			UserName: sanitize(args[6], "string").(string),
		}
		PackageOBJ.RequestedOn = sanitize(args[8], "string").(string)
		PackageOBJ.Progress = 0
		PackageOBJ.Status = "In Progress"
		PackageOBJ.Documents = UploadedDocumentsOBJ
	}
	if packageMode == "UPDATE" {
		PackageData, errorFetch := fetchData(stub, Key, "package")
		if errorFetch != nil {
			fmt.Println("\n Package fetch error>>", errorFetch)
			options = append(options, "Package")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
			// return shim.Error(string(errBytes))
		} else if PackageData == nil {
			fmt.Println("Package not found")
			options = append(options, "Package", Key)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
			// return shim.Error(string(errBytes))
		}

		errorUnmarshal := json.Unmarshal(PackageData, &PackageOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
		fmt.Println("\n\n\nPackage >>>>>", PackageOBJ)

		// PackageOBJ.Documents = append(UploadedDocumentsOBJ, PackageOBJ.Documents...)
		previousDocs = PackageOBJ.Documents

		for pd := 0; pd < len(previousDocs); pd++ {

			for p := 0; p < len(UploadedDocumentsOBJ); p++ {
				if previousDocs[pd].DocumentType == UploadedDocumentsOBJ[p].DocumentType {
					UploadedDocumentsOBJ[p].SignatureOrgs = previousDocs[pd].SignatureOrgs
					UploadedDocumentsOBJ[p].SignatureGroups = previousDocs[pd].SignatureGroups

				}
			}

		}

		PackageOBJ.Documents = UploadedDocumentsOBJ // for rejection flow
		fmt.Println("Updated Documents!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", PackageOBJ.Documents)
		//maintain document level history
	}

	if packageMode != "UPDATE" && packageMode != "NEW" {
		options = append(options)
		errBytes, _ := prepareErrorCode("2029", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Requested Mode is not valid"}, "90700006", false, true)
		// return shim.Error(string(errBytes))
	}

	//upsert
	fmt.Println("PackageOBJ >>>", PackageOBJ)

	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "package", PackageMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("------------------------------------------------------------------------------")
	fmt.Println("Upsert Successfully in Package Collection!")

	//
	for j := 0; j < len(UploadedDocumentsOBJ); j++ {
		UOBJ := UploadedDocumentsOBJ[j]

		DocumentOBJ := Document{}
		DocumentOBJ.Key = PackageOBJ.Key + "-" + UOBJ.Hash
		DocumentOBJ.Version = 0
		DocumentOBJ.GUID = UOBJ.GUID
		DocumentOBJ.Status = "In Progress"
		DocumentOBJ.Progress = 0
		DocumentOBJ.PackageId = Key
		DocumentOBJ.DocumentName = "document"
		DocumentOBJ.DocumentType = UOBJ.DocumentType
		DocumentOBJ.Name = UOBJ.Name
		DocumentOBJ.SetNo = UOBJ.SetNo
		DocumentOBJ.Extension = UOBJ.Extension
		DocumentOBJ.DispatchEmail = UOBJ.DispatchEmail
		DocumentOBJ.EmailTemplate = UOBJ.EmailTemplate
		DocumentOBJ.HolderName = UOBJ.HolderName
		DocumentOBJ.SignaturePolicy = UOBJ.SignaturePolicy
		DocumentOBJ.DocumentPolicy = UOBJ.DocumentPolicy
		DocumentOBJ.FinalSignatoryOrganizationCode = UOBJ.FinalSignatoryOrganizationCode
		DocumentOBJ.FinalSignatoryOrganizationName = UOBJ.FinalSignatoryOrganizationName
		DocumentOBJ.FinalSignatoryOrganizationType = UOBJ.FinalSignatoryOrganizationType
		DocumentOBJ.FinalSignatoryGroup = UOBJ.FinalSignatoryGroup

		// DocumentOBJ.Sequence = false
		DocumentOBJ.CurrentVersion = History{
			Type:       "Uploaded",
			OrgCode:    PackageOBJ.RequestedBy.OrgCode,
			OrgName:    PackageOBJ.RequestedBy.OrgName,
			UserId:     PackageOBJ.RequestedBy.UserId,
			UserName:   PackageOBJ.RequestedBy.UserName,
			Datetime:   PackageOBJ.RequestedOn,
			Hash:       UOBJ.Hash,
			SequenceNo: 0,
			Status:     "Uploaded",
			TranxHash:  tnxId,
			UserGroup:  sanitize(args[15], "string").(string),
		}

		if UploadedDocumentsOBJ[j].QRHash != "" {
			DocumentMappingStruct.Key = "QR_" + UploadedDocumentsOBJ[j].QRHash
			DocumentMappingStruct.DocumentKey = DocumentOBJ.Key
			//upsert documentMapping
			DocumentMappingMarshalled, errorMarshal := json.Marshal(DocumentMappingStruct)
			if errorMarshal != nil {
				fmt.Println("Error occurred while marshalling the data")
				options = append(options)
				errBytes, _ := prepareErrorCode("2003", options)
				return shim.Error(string(errBytes))
			}
			errorInsert := insertData(&stub, DocumentMappingStruct.Key, DocumentMappingCollection, DocumentMappingMarshalled)
			if errorInsert != nil {
				fmt.Println("Error occurred while inserting the DocumentMapping in private collection")
				options = append(options)
				errBytes, _ := prepareErrorCode("2004", options)
				return shim.Error(string(errBytes))
			}
		}

		// Previous Docs for maintaining DocHistory

		if packageMode == "UPDATE" {
			prevUploadDocOBJ := UploadedDocuments{}
			prevDocOBJ := Document{}
			for p := 0; p < len(previousDocs); p++ {
				if previousDocs[p].DocumentType == UOBJ.DocumentType {
					prevUploadDocOBJ = previousDocs[p]
				}
			}

			//get Previous Documents
			if len(prevUploadDocOBJ.Hash) > 0 {
				PrevDocData, errorFetch := fetchData(stub, PackageOBJ.Key+"-"+prevUploadDocOBJ.Hash, "document")
				fmt.Println("\n PrevDocData upsertPackage3>>")
				if errorFetch != nil {
					fmt.Println("\n PrevDoc fetch error upsertPackage>>", errorFetch)
					options = append(options, "PrevDoc")
					errBytes, _ := prepareErrorCode("2001", options)
					return shim.Error(string(errBytes))
					// errBytes := genericErrorHandler(nil, []string{"PrevDoc fetch error"}, "90200001", false, true)
					// return shim.Error(string(errBytes))
				} else if PrevDocData == nil {
					fmt.Println("PrevDoc not found")
					options = append(options, "PrevDoc", PackageOBJ.Key+"-"+prevUploadDocOBJ.Hash)
					errBytes, _ := prepareErrorCode("2002", options)
					return shim.Error(string(errBytes))
					// errBytes := genericErrorHandler(nil, []string{"PrevDoc not found"}, "90200002", false, true)
					// return shim.Error(string(errBytes))
				}

				errorUnmarshal := json.Unmarshal(PrevDocData, &prevDocOBJ)
				if errorUnmarshal != nil {
					options = append(options)
					errBytes, _ := prepareErrorCode("2000", options)
					return shim.Error(string(errBytes))
					// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
					// return shim.Error(string(errBytes))
				}
				fmt.Println("\n\n\nPREVIOUS DOCUMENT OBJ >>>>>", prevDocOBJ)

				DocumentOBJ.DocumentHistory = append(DocumentOBJ.DocumentHistory, DocumentHistory{
					DocumentId: prevDocOBJ.Key,
					Version:    prevDocOBJ.Version,
				})
				DocumentOBJ.Version = prevDocOBJ.Version + 1

			}
		}

		//insert in document hash collection
		DocumentHashOBJ := DocumentHash{}
		DocumentHashOBJ.Key = UOBJ.Hash
		DocumentHashOBJ.DocumentName = "documentHash"
		DocumentHashOBJ.PackageId = PackageOBJ.Key + "-" + UOBJ.Hash

		SignPolicyOBJ := SignaturePolicy{}
		SignPolicyData, errorFetch := fetchData(stub, UOBJ.SignaturePolicy, "SignaturePolicy")
		fmt.Println("\n sign policy Key upsertPackage2>>", UOBJ.SignaturePolicy)
		if errorFetch != nil {
			fmt.Println("\n SignPolicy fetch error upsertPackage>>", errorFetch)
			options = append(options, "SignPolicy")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"SignPolicy fetch error"}, "90200001", false, true)
			// return shim.Error(string(errBytes))
		} else if SignPolicyData == nil {
			fmt.Println("SignPolicy not found")
			options = append(options, "SignPolicy", UOBJ.SignaturePolicy)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"SignPolicy not found"}, "90200002", false, true)
			// return shim.Error(string(errBytes))
		}

		errorUnmarshal := json.Unmarshal(SignPolicyData, &SignPolicyOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
		fmt.Println("\n\n\nPackage >>>>>", PackageOBJ)

		fmt.Println("\n\n\nSign Policy >>>>>", SignPolicyOBJ)

		//update sequence
		DocumentOBJ.Sequence = SignPolicyOBJ.Sequence

		optionalORGS := 0

		for i := 0; i < len(SignPolicyOBJ.Organizations); i++ {
			orgObj := SignPolicyOBJ.Organizations[i]
			signature := false
			stamp := false
			optionalOrgs := false
			ActionType := ""

			if orgObj.IsReviewer == true {
				ActionType = "Reviewer"
			} else {
				for _, page := range orgObj.Pages {
					if page.Type == "Signature" {
						signature = true
					}
					if page.Type == "Stamp" {
						stamp = true
					}
				}

				if stamp && signature {
					ActionType = "Both"
				} else if signature {
					ActionType = "Signature"
				} else if stamp {
					ActionType = "Stamp"
				}
			}

			var orgCode = ""

			if orgObj.OrganizationCode != "" {
				//do something here
				orgCode = orgObj.OrganizationCode
			} else {
				for _, value := range OverRideOrg {
					if value.OrgType == orgObj.OrganizationType {
						// logic
						orgCode = value.OrgCode
					}
				}

				if len(OverRideOrg) == 0 {
					orgCode = ""
				}

			}

			for _, value := range UOBJ.OptionalOrgs {
				if value == orgCode {
					// logic
					optionalOrgs = true
					optionalORGS = optionalORGS + 1
				}
			}

			var emptyArrayOrgs []string
			if len(UOBJ.OptionalOrgs) == 0 {
				UOBJ.OptionalOrgs = emptyArrayOrgs
			}

			fmt.Println("\n\n\nORG CODE >>>>>", orgCode)
			fmt.Println("OPTIONAL ORG CODE >>>>>", optionalOrgs)
			fmt.Println("OPTIONAL ORG CODE >>>>>", optionalORGS)
			fmt.Println("OPTIONAL ORG CODE >>>>>", orgObj.SequenceNo-optionalORGS)
			if !optionalOrgs {
				SIGOBJ := Signatures{
					OrgCode:         orgCode,
					OrgName:         orgObj.OrganizationName,
					Type:            ActionType,
					Group:           orgObj.Group,
					Action:          "Pending",
					SLA:             orgObj.SLA,
					SequenceNo:      orgObj.SequenceNo - optionalORGS,
					Pages:           orgObj.Pages,
					IsReviewer:      orgObj.IsReviewer,
					IsFlexibleCords: orgObj.IsFlexibleCords,
					IsDispatch:      orgObj.DispatchStatus,
				}
				DocumentOBJ.Signatures = append(DocumentOBJ.Signatures, SIGOBJ)
			}

		}

		if len(DocumentOBJ.Signatures) > 0 {
			smallest := DocumentOBJ.Signatures[0]            // set the smallest number to the first element of the list
			for _, num := range DocumentOBJ.Signatures[1:] { // iterate over the rest of the list
				if num.SequenceNo < smallest.SequenceNo { // if num is smaller than the current smallest number
					smallest = num // set smallest to num
				}
			}
			fmt.Println("\n\n\nPackage >>>>>", smallest)

			DocumentOBJ.NextSignatory = smallest.OrgCode

			// @REVERTED
			// var taskTypes = getTaskTypes()
			// var TaskName = taskTypes.SIGNATURE_REQUEST
			// var UUID string = DocumentOBJ.PackageId + "_" + strings.Replace(DocumentOBJ.DocumentType, " ", "_", -1) + "_" + smallest.OrgCode
			// error := addTask(stub, UUID, TaskName, DocumentOBJ.DocumentType, smallest.OrgCode, smallest.Group, "", "pending", DocumentOBJ.PackageId, "", "NEW")
			// if error != "" {
			// 	return shim.Error(error)
			// }
		}

		//upsert document
		DocumentMarshalled, errorMarshal := json.Marshal(DocumentOBJ)
		if errorMarshal != nil {
			fmt.Println("Error occurred while marshalling the data")
			options = append(options)
			errBytes, _ := prepareErrorCode("2003", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
		errorInsert := insertData(&stub, DocumentOBJ.Key, "document", DocumentMarshalled)
		if errorInsert != nil {
			fmt.Println("Error occurred while inserting the Document in private collection")
			options = append(options)
			errBytes, _ := prepareErrorCode("2004", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
			// return shim.Error(string(errBytes))
		}
		fmt.Println("Upsert Successfully in Document Collection!")

		//upsert document hash
		DocumentHashMarshalled, errorMarshal := json.Marshal(DocumentHashOBJ)
		if errorMarshal != nil {
			fmt.Println("Error occurred while marshalling the data")
			options = append(options)
			errBytes, _ := prepareErrorCode("2003", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
		errorInsert = insertData(&stub, DocumentHashOBJ.Key, "documentHash", DocumentHashMarshalled)
		if errorInsert != nil {
			fmt.Println("Error occurred while inserting the Document Hash in private collection")
			options = append(options)
			errBytes, _ := prepareErrorCode("2004", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
			// return shim.Error(string(errBytes))
		}
		fmt.Println("Upsert Successfully in Document Hash Collection!")
	}

	return shim.Success(nil)
}

func (t *CoreChainCode) getPackage(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("package: %v", args)
	Key := sanitize(args[0], "string").(string)
	var options []string
	PackageOBJ := PackageDTO{}
	PackageCheck, errorFetch := fetchData(stub, Key, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
		// return shim.Error(string(errBytes))
	} else if PackageCheck == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("Package exist")
		errorUnmarshal := json.Unmarshal(PackageCheck, &PackageOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nSign Policy >>>>>", PackageOBJ)
	for i := 0; i < len(PackageOBJ.Documents); i++ {
		docs := PackageOBJ.Documents[i]
		actdocs := t.getDocumentLocally(stub, Key+"-"+docs.Hash, "getDocumentInternally")
		PackageOBJ.DocsPackage = append(PackageOBJ.DocsPackage, actdocs)
		PackageOBJ.Documents[i].PackageId = Key
	}

	for i := 0; i < len(PackageOBJ.OtherDocuments); i++ {
		docs := PackageOBJ.OtherDocuments[i]
		actpack := t.getPackageLocally(stub, docs.PackageId, "getPackageInternally")
		PackageOBJ.DocsPackage = append(PackageOBJ.DocsPackage, actpack.DocsPackage...)

		actpack.Documents[0].Name = docs.DocumentType
		actpack.Documents[0].IsReadableOnly = true
		actpack.Documents[0].PackageId = docs.PackageId
		PackageOBJ.Documents = append(PackageOBJ.Documents, actpack.Documents...)

	}

	fmt.Println("\n\n\nPkg Type >>>>>", PackageOBJ)

	Key1 := PackageOBJ.PackageType
	fmt.Println("\n\n\nPackage Type Key >>>>>", Key1)
	PackageTypeCheck, errorFetch := fetchData(stub, Key1, "packageType")
	if errorFetch != nil {
		fmt.Println("\n PackageType fetch error>>", errorFetch)
		options = append(options, "PackageType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType fetch error"}, "90700001", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageTypeCheck == nil {
		fmt.Println("PackageType not found")
		options = append(options, "PackageType", Key1)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType not found"}, "90700002", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("PackageType exist")

		errorUnmarshal := json.Unmarshal(PackageTypeCheck, &PackageOBJ.PackageTypeArr)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	for i := 0; i < len(PackageOBJ.PackageTypeArr.NotifyPolicy.ApproveParties); i++ {
		fmt.Println("ApproveParties exist", PackageOBJ.PackageTypeArr.NotifyPolicy.ApproveParties[i].OrgType)
		for _, override := range PackageOBJ.OverrideOrgs {
			fmt.Println("OverrideOrgs exist", override.OrgType)
			if PackageOBJ.PackageTypeArr.NotifyPolicy.ApproveParties[i].OrgType == override.OrgType {
				fmt.Println("OverrideOrgs And ApproveParties exist",
					override.OrgType+" = "+PackageOBJ.PackageTypeArr.NotifyPolicy.ApproveParties[i].OrgType)
				PackageOBJ.PackageTypeArr.NotifyPolicy.ApproveParties[i].OrgCode = override.OrgCode
			}
		}
	}
	for i := 0; i < len(PackageOBJ.PackageTypeArr.NotifyPolicy.RejectParties); i++ {
		for _, override := range PackageOBJ.OverrideOrgs {
			if PackageOBJ.PackageTypeArr.NotifyPolicy.RejectParties[i].OrgType == override.OrgType {
				PackageOBJ.PackageTypeArr.NotifyPolicy.RejectParties[i].OrgCode = override.OrgCode
			}
		}
	}
	fmt.Println("ApproveParties exist", PackageOBJ.PackageTypeArr.NotifyPolicy.ApproveParties)
	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("PackageMarshalled marshalled successful !!!")
	return shim.Success(PackageMarshalled)
}

func (t *CoreChainCode) getPackage2(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("package: %v", args)
	Key := sanitize(args[0], "string").(string)
	var options []string
	PackageOBJ := PackageDTO{}
	PackageCheck, errorFetch := fetchData(stub, Key, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
		// return shim.Error(string(errBytes))
	} else if PackageCheck == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("Package exist")
		errorUnmarshal := json.Unmarshal(PackageCheck, &PackageOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nSign Policy >>>>>", PackageOBJ)
	for i := 0; i < len(PackageOBJ.Documents); i++ {
		docs := PackageOBJ.Documents[i]
		actdocs := t.getDocumentLocally(stub, Key+"-"+docs.Hash, "getDocumentInternally")
		PackageOBJ.DocsPackage = append(PackageOBJ.DocsPackage, actdocs)
	}

	fmt.Println("\n\n\nPkg Type >>>>>", PackageOBJ)

	Key1 := PackageOBJ.PackageType
	fmt.Println("\n\n\nPackage Type Key >>>>>", Key1)
	PackageTypeCheck, errorFetch := fetchData(stub, Key1, "packageType")
	if errorFetch != nil {
		fmt.Println("\n PackageType fetch error>>", errorFetch)
		options = append(options, "PackageType")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType fetch error"}, "90700001", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageTypeCheck == nil {
		fmt.Println("PackageType not found")
		options = append(options, "PackageType", Key1)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"PackageType not found"}, "90700002", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("PackageType exist")

		errorUnmarshal := json.Unmarshal(PackageTypeCheck, &PackageOBJ.PackageTypeArr)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("PackageMarshalled marshalled successful !!!")
	return shim.Success(PackageMarshalled)
}

func (t *CoreChainCode) getDocumentLocally(stub hypConnect, Key string, functionName string) Document {
	fmt.Printf("Key:" + Key)
	//Key := sanitize(args[0], "string").(string)
	DocumentOBJ := Document{}
	DocumentCheck, errorFetch := fetchData(stub, Key, "document")
	if errorFetch != nil {
		fmt.Println("\n Document fetch error>>", errorFetch)
		errBytes := genericErrorHandler(nil, []string{"Document fetch error"}, "90600003", false, true)
		fmt.Println("Document not found", errBytes)
		return DocumentOBJ
	} else if DocumentCheck == nil {
		fmt.Println("Document not found")
		errBytes := genericErrorHandler(nil, []string{"Document not found"}, "90600004", false, true)
		fmt.Println("Document not found", errBytes)
		return DocumentOBJ
	} else {
		fmt.Println("Document exist")
		errorUnmarshal := json.Unmarshal(DocumentCheck, &DocumentOBJ)
		if errorUnmarshal != nil {
			errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			fmt.Println("Document not found", errBytes)
			return DocumentOBJ
		}
	}
	fmt.Println("\n\n\nSign Policy >>>>>", DocumentOBJ)
	return DocumentOBJ
}

func (t *CoreChainCode) getPackageLocally(stub hypConnect, Key string, functionName string) PackageDTO {
	fmt.Printf("package: %v", Key)
	PackageOBJ := PackageDTO{}
	PackageCheck, errorFetch := fetchData(stub, Key, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		return PackageOBJ
	} else if PackageCheck == nil {
		fmt.Println("Package not found")
		return PackageOBJ
	} else {
		fmt.Println("Package exist")
		errorUnmarshal := json.Unmarshal(PackageCheck, &PackageOBJ)
		if errorUnmarshal != nil {
			return PackageOBJ
		}
	}
	fmt.Println("\n\n\nSign Policy >>>>>", PackageOBJ)
	for i := 0; i < len(PackageOBJ.Documents); i++ {
		docs := PackageOBJ.Documents[i]
		actdocs := t.getDocumentLocally(stub, Key+"-"+docs.Hash, "getDocumentInternally")
		PackageOBJ.DocsPackage = append(PackageOBJ.DocsPackage, actdocs)
	}

	return PackageOBJ
}
func (t *CoreChainCode) rejectPackage(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("reject packagesssss: %v", args)

	Key := sanitize(args[0], "string").(string)
	OrgCode := sanitize(args[1], "string").(string)
	OrgName := sanitize(args[2], "string").(string)
	UserId := sanitize(args[3], "string").(string)
	UserName := sanitize(args[4], "string").(string)
	Datetime := sanitize(args[5], "string").(string)
	var options []string
	PackageOBJ := Package{}
	PackageCheck, errorFetch := fetchData(stub, Key, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageCheck == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Package exist")

	errorUnmarshal := json.Unmarshal(PackageCheck, &PackageOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("\n\n\nPackage unmarshal >>>>>", PackageOBJ)

	//action allowed check
	if PackageOBJ.RequestedBy.OrgCode != OrgCode || PackageOBJ.RequestedBy.UserId != UserId {
		fmt.Println("You're not allowed to perform this action")
		options = append(options)
		errBytes, _ := prepareErrorCode("2030", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"You're not allowed to perform this action"}, "90700009", false, true)
		// return shim.Error(string(errBytes))
	}

	if PackageOBJ.Status == "Rejected" {
		fmt.Println("Package is already Rejected")
		options = append(options)
		errBytes, _ := prepareErrorCode("2031", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package is already Rejected"}, "90700007", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageOBJ.Status == "Completed" {
		fmt.Println("Package is already completed.")
		options = append(options)
		errBytes, _ := prepareErrorCode("2032", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package is already completed."}, "90700008", false, true)
		// return shim.Error(string(errBytes))

	}

	fmt.Println("Package is In Progress.", PackageOBJ.Status)

	PackageOBJ.Status = "Rejected"

	PackageOBJ.RequestedBy = RequestedBy{
		OrgCode:  OrgCode,
		OrgName:  OrgName,
		UserId:   UserId,
		UserName: UserName,
	}
	PackageOBJ.RequestedOn = Datetime

	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("Reject PackageMarshalled marshalled successful !!!")

	errorInsert := insertData(&stub, Key, "package", PackageMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Reject Successfully in Package Collection!")

	return shim.Success(nil)
}

func addTask(stub hypConnect, UUID string, TaskName string, DocumentType []string, OrgCode string, GroupName string, Port string, Status string, PackageId string, NotificationId string, SNNumber string, TaskStatus string, SLATime int, AdditonalData string, RequestTime string, stage string, activity string) string {
	fmt.Println("ADD TASK >>>>>>>>>>>>>> ")
	fmt.Println(UUID, TaskName, DocumentType, OrgCode, GroupName, Port, Status, PackageId, SNNumber, TaskStatus, SLATime, AdditonalData)
	taskGroup := GroupName
	var options []string
	if TaskStatus == "UPDATE" {
		//fetch task and throw error if not exist
		TaskOBJ := Tasks{}

		TaskData, errorFetch := fetchData(stub, UUID, "tasks")
		if errorFetch != nil {
			fmt.Println("\n Task fetch error>>", errorFetch)
			options = append(options, "Task")
			errBytes, _ := prepareErrorCode("2001", options)
			return string(errBytes)
			//return "Task fetch error"

		} else if TaskData == nil {
			fmt.Println("Task not found")
			options = append(options, "Task", UUID)
			errBytes, _ := prepareErrorCode("2002", options)
			return string(errBytes)
			//return "Task not found"
		}
		errorUnmarshal := json.Unmarshal(TaskData, &TaskOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return string(errBytes)
			//return "Unmarshalling error"
		}

		taskGroup = TaskOBJ.GroupName
		SLATime = TaskOBJ.SLATime
		if TaskOBJ.Stage != "" {
			stage = TaskOBJ.Stage
		}
		if TaskOBJ.Activity != "" {
			activity = TaskOBJ.Activity
		}
	}

	task := &Tasks{
		DocumentName:   "usertask",
		Key:            UUID,
		TaskName:       TaskName,
		DocumentType:   DocumentType,
		OrgCode:        OrgCode,
		GroupName:      taskGroup,
		Port:           Port,
		Status:         Status,
		PackageId:      PackageId,
		NotificationId: NotificationId,
		SNNumber:       SNNumber,
		SLATime:        SLATime,
		AdditonalData:  AdditonalData,
		Stage:          stage,
		Activity:       activity,
	}
	if Status == "Completed" {
		task.CompletedAt = RequestTime
	}

	taskMarshalled, errorMarshal := json.Marshal(task)
	if errorMarshal != nil {
		fmt.Println("Marshal failed for task!!!!!!", errorMarshal)
	}
	fmt.Println("Marshalled task Array: ", taskMarshalled)

	errorInsert := insertData(&stub, task.Key, "tasks", taskMarshalled)
	//raise shim error message if insertion fails.
	if errorInsert != nil {
		fmt.Println("Insertion failed of task !!!!!!", errorInsert)
	}

	return ""
}

func (t *CoreChainCode) signDocument(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign document: %v", args)
	NextSignatoryOBJ := NextSignatoryInfo{}

	// @REVERTED
	// isNextSignatoryPresent := false

	key := sanitize(args[0], "string").(string)
	PackageNo := sanitize(args[10], "string").(string)
	tnxId := stub.Connection.GetTxID()
	currentUserGroup := sanitize(args[17], "string").(string)

	// Current Version
	CurrentVersion := History{
		Type:           sanitize(args[5], "string").(string),
		OrgCode:        sanitize(args[1], "string").(string),
		OrgName:        sanitize(args[2], "string").(string),
		UserId:         sanitize(args[3], "string").(string),
		UserName:       sanitize(args[4], "string").(string),
		Datetime:       sanitize(args[8], "string").(string),
		Hash:           sanitize(args[7], "string").(string),
		Status:         sanitize(args[6], "string").(string),
		Comments:       sanitize(args[9], "string").(string),
		RejectReason:   sanitize(args[11], "string").(string),
		ReturnComments: sanitize(args[12], "string").(string),
		ReturnReason:   sanitize(args[13], "string").(string),
		UserSite:       sanitize(args[14], "string").(string),
		UserDepartment: sanitize(args[15], "string").(string),
		TranxHash:      tnxId,
		UserGroup:      currentUserGroup,
	}

	var options []string

	//insert in document hash collection
	var DocumentMappingCollection = "documentMapping"
	DocumentHashOBJ := DocumentHash{}
	DocumentHashOBJ.Key = sanitize(args[7], "string").(string)
	DocumentHashOBJ.DocumentName = "documentHash"
	DocumentHashOBJ.PackageId = key

	// upsert Package
	var DocumentMappingStruct DocumentMapping

	if sanitize(args[7], "string").(string) != "" {
		DocumentMappingStruct.Key = "FILE_" + sanitize(args[7], "string").(string)
		DocumentMappingStruct.DocumentKey = key
		fmt.Println("\n\n\nDocumentMappingStruct >>>>>", DocumentMappingStruct)
		//upsert documentMapping
		DocumentMappingMarshalled, errorMarshal := json.Marshal(DocumentMappingStruct)
		if errorMarshal != nil {
			fmt.Println("Error occurred while marshalling the data")
			options = append(options)
			errBytes, _ := prepareErrorCode("2003", options)
			return shim.Error(string(errBytes))
		}
		errorInsert := insertData(&stub, DocumentMappingStruct.Key, DocumentMappingCollection, DocumentMappingMarshalled)
		if errorInsert != nil {
			fmt.Println("Error occurred while inserting the DocumentMapping in private collection")
			options = append(options)
			errBytes, _ := prepareErrorCode("2004", options)
			return shim.Error(string(errBytes))
		}
	}

	DocumentOBJ := Document{}
	DocumentData, errorFetch := fetchData(stub, key, "document")
	if errorFetch != nil {
		fmt.Println("\n Document fetch error>>", errorFetch)
		options = append(options, "Document")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document fetch error"}, "90600003", false, true)
		// return shim.Error(string(errBytes))
	} else if DocumentData == nil {
		fmt.Println("Document not found")
		options = append(options, "Document", key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document not found"}, "90600004", false, true)
		// return shim.Error(string(errBytes))
	}
	errorUnmarshal := json.Unmarshal(DocumentData, &DocumentOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("\n\n\nDocument Object >>>>>", DocumentOBJ)

	if DocumentOBJ.Status == "Rejected" {
		options = append(options)
		errBytes, _ := prepareErrorCode("2040", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document is already rejected, cannot be signed"}, "90600005", false, true)
		// return shim.Error(string(errBytes))
	}

	//insert current object in hostory
	DocumentOBJ.History = append(DocumentOBJ.History, DocumentOBJ.CurrentVersion)

	TotalSigned := 0.0
	SigFound := false
	isDispatch := false

	isLastReviewer := false

	for s := 0; s < len(DocumentOBJ.Signatures); s++ {
		SIGNOBJ := DocumentOBJ.Signatures[s]
		if SIGNOBJ.Action == "Signed" {
			TotalSigned++
		}
		isSignatoryGroup := false
		for _, group := range SIGNOBJ.Group {
			if group == currentUserGroup {
				isSignatoryGroup = true
				break
			}
		}
		if SIGNOBJ.OrgCode == CurrentVersion.OrgCode && SIGNOBJ.Type == CurrentVersion.Type && isSignatoryGroup {
			SigFound = true
			TotalSigned++
			DocumentOBJ.NextSignatory = ""
			fmt.Println("\n\n\nSequence NO >>>>>", SIGNOBJ.SequenceNo)
			fmt.Println("\n\n\nNext Signatory sequence No >>>>>", SIGNOBJ.SequenceNo+1)
			for n := 0; n < len(DocumentOBJ.Signatures); n++ {
				fmt.Println("\n\n\nNext Signatory sequence No IN LOOP >>>>>", DocumentOBJ.Signatures[n].SequenceNo)
				if SIGNOBJ.SequenceNo+1 == DocumentOBJ.Signatures[n].SequenceNo {
					DocumentOBJ.NextSignatory = DocumentOBJ.Signatures[n].OrgCode
					fmt.Println("\n\n\nNext Signatory For this Doc >>>>>", DocumentOBJ.Signatures[n].OrgCode)

					//set NextSignatoryOBJ
					// @REVERTED
					// isNextSignatoryPresent = true
					NextSignatoryOBJ.OrgCode = DocumentOBJ.Signatures[n].OrgCode
					NextSignatoryOBJ.Group = DocumentOBJ.Signatures[n].Group
				}
			}

			DocumentOBJ.Signatures[s].Action = CurrentVersion.Status
			DocumentOBJ.Signatures[s].Datetime = CurrentVersion.Datetime
			//check whether the given signature is in order
			if DocumentOBJ.Sequence && DocumentOBJ.CurrentVersion.SequenceNo+1 != SIGNOBJ.SequenceNo {
				options = append(options)
				errBytes, _ := prepareErrorCode("2033", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Organization is not allowed not perform this action"}, "90500003", false, true)
				// return shim.Error(string(errBytes))
			}

			if CurrentVersion.Type == "Reviewer" && s == (len(DocumentOBJ.Signatures)-1) {
				isLastReviewer = true
			}
			isDispatch = SIGNOBJ.IsDispatch
		}
	}
	fmt.Println("\n\nDocument Sequence >>>", DocumentOBJ.Sequence, DocumentOBJ.CurrentVersion.SequenceNo+1)

	if DocumentOBJ.Sequence {
		CurrentVersion.SequenceNo = DocumentOBJ.CurrentVersion.SequenceNo + 1
	}
	fmt.Println("\n\nDocument Sequence No updated>>>", DocumentOBJ.CurrentVersion.SequenceNo)

	//don't raise progress if status is REJECTED
	if CurrentVersion.Status == "Rejected" {
		DocumentOBJ.Status = "Rejected"

		if isLastReviewer {
			DocumentOBJ.ReturnedVersions = append(DocumentOBJ.ReturnedVersions, DocumentOBJ.CurrentVersion)
		}
	} else if CurrentVersion.Status == "Signed" {
		fmt.Println("\n\n\nProgress >>>", TotalSigned, len(DocumentOBJ.Signatures), (TotalSigned/float64(len(DocumentOBJ.Signatures)))*100)
		DocumentOBJ.Progress = toFixed((TotalSigned/float64(len(DocumentOBJ.Signatures)))*100, 2)
	}

	fmt.Println("\n\n\nProgress Values>>>", DocumentOBJ.Progress)
	if DocumentOBJ.Progress == 100 {
		DocumentOBJ.Status = "Signed"
	}

	// set overall package status and progress
	PackageOBJ := Package{}
	PackageCheck, errorFetch := fetchData(stub, PackageNo, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
		// return shim.Error(string(errBytes))
	} else if PackageCheck == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", PackageNo)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("Package exist")

	errorUnmarshal1 := json.Unmarshal(PackageCheck, &PackageOBJ)
	if errorUnmarshal1 != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	comments := sanitize(args[16], "string").(string)
	if PackageOBJ.Progress == 0 && PackageOBJ.IsRegenrated && CurrentVersion.Status != "Rejected" && CurrentVersion.Type != "Reviewer" {
		if comments == "" {
			options = []string{"Comments are required when signing regenerated package"}
			errBytes, _ := prepareErrorCode("2036", options)
			return shim.Error(string(errBytes))
		} else {
			CurrentVersion.Comments = comments
		}
	}

	SignedDocument := 0
	for i := 0; i < len(PackageOBJ.Documents); i++ {
		fmt.Println("Documents array value >>>", PackageOBJ.Documents[i])
		DocumentKey := PackageNo + "-" + PackageOBJ.Documents[i].Hash
		fmt.Println("compare >>>", key, DocumentKey)

		if key != DocumentKey {
			DocumentOBJ1 := Document{}
			DocumentData, errorFetch := fetchData(stub, DocumentKey, "document")
			if errorFetch != nil {
				fmt.Println("\n Document fetch error>>", errorFetch)
				options = append(options, "Document")
				errBytes, _ := prepareErrorCode("2001", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Document fetch error"}, "90600003", false, true)
				// return shim.Error(string(errBytes))

			} else if DocumentData == nil {
				options = append(options, "Document", DocumentKey)
				errBytes, _ := prepareErrorCode("2002", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Document not found"}, "90600004", false, true)
				// return shim.Error(string(errBytes))
			}

			errorUnmarshal := json.Unmarshal(DocumentData, &DocumentOBJ1)
			if errorUnmarshal != nil {
				options = append(options)
				errBytes, _ := prepareErrorCode("2000", options)
				return shim.Error(string(errBytes))
				// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
				// return shim.Error(string(errBytes))
			}

			if DocumentOBJ1.Status == "Signed" {
				SignedDocument++
			}

			if isDispatch {
				fmt.Println("chck for other docs  -----------")
			signatureLoop:
				for s := 0; s < len(DocumentOBJ1.Signatures); s++ {
					SIGNOBJ1 := DocumentOBJ1.Signatures[s]
					fmt.Println("SIGNOBJ1  -----------", SIGNOBJ1, currentUserGroup)
					for _, group := range SIGNOBJ1.Group {
						if group == currentUserGroup && SIGNOBJ1.Action != "Signed" {
							isDispatch = false
							break signatureLoop
						}
					}
				}
			}
		}

		//update Progress of individual Douments in Pkg
		fmt.Println("\n\n Progress for Pkg individual doc!!!!!!!!!!!!!! -----------", isDispatch)

		if key == DocumentKey {
			PackageOBJ.Documents[i].Progress = DocumentOBJ.Progress
		}
		fmt.Println("\n\n Progress-----------", PackageOBJ.Documents[i].Progress)

	}
	if DocumentOBJ.Status == "Signed" {
		SignedDocument++
	}
	PackageOBJ.Progress = toFixed((float64(SignedDocument)/float64(len(PackageOBJ.Documents)))*100, 2)
	if PackageOBJ.Progress == 100 {
		PackageOBJ.Status = "Completed"
	}

	if PackageOBJ.IsRegenrated {
		PackageOBJ.IsRegenrated = false
	}
	PackageOBJ.IsDispatch = isDispatch

	//upsert
	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, PackageNo, "package", PackageMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in Package Collection!")

	//-------------------------------------------------------------------------------------------------------------------------

	if !SigFound {
		options = append(options)
		errBytes, _ := prepareErrorCode("2034", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Organization is not part of Signatuire Policy"}, "90200003", false, true)
		// return shim.Error(string(errBytes))
	}

	DocumentOBJ.CurrentVersion = CurrentVersion

	//upsert document
	DocumentMarshalled, errorMarshal := json.Marshal(DocumentOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert1 := insertData(&stub, key, "document", DocumentMarshalled)
	if errorInsert1 != nil {
		fmt.Println("Error occurred while inserting the Document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("Upsert Successfully in Document Collection!")

	//upsert document hash
	DocumentHashMarshalled, errorMarshal := json.Marshal(DocumentHashOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert = insertData(&stub, DocumentHashOBJ.Key, "documentHash", DocumentHashMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Document Hash in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}
	fmt.Println("Upsert Successfully in Document Hash Collection!")

	// @REVERTED
	// var taskTypes = getTaskTypes()
	// var TaskName = taskTypes.SIGNATURE_REQUEST

	// // Add Task for Next Signatory
	// if isNextSignatoryPresent == true {
	// 	var UUID string = DocumentOBJ.PackageId + "_" + strings.Replace(DocumentOBJ.DocumentType, " ", "_", -1) + "_" + NextSignatoryOBJ.OrgCode
	// 	error := addTask(stub, UUID, TaskName, DocumentOBJ.DocumentType, NextSignatoryOBJ.OrgCode, NextSignatoryOBJ.Group, "", "pending", DocumentOBJ.PackageId, "", "NEW")
	// 	if error != "" {
	// 		return shim.Error(error)
	// 	}
	// }
	// var UUID string = DocumentOBJ.PackageId + "_" + strings.Replace(DocumentOBJ.DocumentType, " ", "_", -1) + "_" + sanitize(args[1], "string").(string)
	// var emptyGroup [] string
	// error := addTask(stub, UUID, TaskName, DocumentOBJ.DocumentType, sanitize(args[1], "string").(string), emptyGroup, "", "completed", DocumentOBJ.PackageId, "", "UPDATE")
	// if error != "" {
	// 	return shim.Error(error)
	// }

	return shim.Success(nil)
}

func (t *CoreChainCode) upsertDocument(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("document: %v", args)
	Key := sanitize(args[0], "string").(string)
	IPFS := sanitize(args[17], "string").(string)
	Key = Key + hash(IPFS)

	DocumentOBJ := Document{}
	DocumentOBJ.Key = Key
	DocumentOBJ.DocumentName = "document"
	DocumentOBJ.Name = sanitize(args[1], "string").(string)
	DocumentOBJ.Extension = sanitize(args[2], "string").(string)
	DocumentOBJ.DocumentType = sanitize(args[3], "string").(string)
	DocumentOBJ.SignaturePolicy = sanitize(args[4], "string").(string)
	DocumentOBJ.Version = sanitize(args[19], "int").(int)
	DocumentOBJ.Status = sanitize(args[5], "string").(string)
	DocumentOBJ.Progress = sanitize(args[6], "float64").(float64)
	DocumentOBJ.Sequence = sanitize(args[18], "bool").(bool)

	Type := sanitize(args[7], "string").(string)
	OrgCode := sanitize(args[8], "string").(string)
	OrgName := sanitize(args[9], "string").(string)
	UserId := sanitize(args[10], "string").(string)
	UserName := sanitize(args[11], "string").(string)
	Hash := sanitize(args[12], "string").(string)
	Order := sanitize(args[13], "int").(int)
	Status := sanitize(args[14], "string").(string)

	// Current Version
	CurrentVersion := History{
		Type:       Type,
		OrgCode:    OrgCode,
		OrgName:    OrgName,
		UserId:     UserId,
		UserName:   UserName,
		Datetime:   time.Now().Format("20060102150405"),
		Hash:       Hash,
		SequenceNo: Order,
		Status:     Status,
	}
	var options []string
	DocumentOBJ.CurrentVersion = CurrentVersion
	DocumentOBJ.History = append(DocumentOBJ.History, CurrentVersion)

	err := json.Unmarshal([]byte(args[15]), &DocumentOBJ.Signatures)
	if err != nil {
		fmt.Printf("\n Error occurred while unmarshalling StampOrg")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	errDocHis := json.Unmarshal([]byte(args[16]), &DocumentOBJ.DocumentHistory)
	if errDocHis != nil {
		fmt.Printf("\n Error occurred while unmarshalling StampOrg")
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("\n\n\nDocument BEF 1 >>>>>", DocumentOBJ)

	DocumentCheck, errorFetch := fetchData(stub, Key, "document")
	if errorFetch != nil {
		fmt.Println("\n Document fetch error>>", errorFetch)
		options = append(options, "Document")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document fetch error"}, "90600003", false, true)
		// return shim.Error(string(errBytes))

	} else if DocumentCheck != nil {
		fmt.Println("Document already exists")
	} else {
		fmt.Println("Document not found")
	}
	fmt.Println("\n\n\nDocument >>>>>", DocumentOBJ)

	//upsert
	DocumentMarshalled, errorMarshal := json.Marshal(DocumentOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert := insertData(&stub, Key, "document", DocumentMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in Document Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) getDocument(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("document: %v", args)

	Key := sanitize(args[0], "string").(string)
	var options []string
	DocumentOBJ := Document{}
	DocumentCheck, errorFetch := fetchData(stub, Key, "document")
	if errorFetch != nil {
		fmt.Println("\n Document fetch error>>", errorFetch)
		options = append(options, "Document")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document fetch error"}, "90600003", false, true)
		// return shim.Error(string(errBytes))

	} else if DocumentCheck == nil {
		fmt.Println("Document not found")
		options = append(options, "Document", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document not found"}, "90600004", false, true)
		// return shim.Error(string(errBytes))
	} else {
		fmt.Println("Document exist")

		errorUnmarshal := json.Unmarshal(DocumentCheck, &DocumentOBJ)
		if errorUnmarshal != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
			// return shim.Error(string(errBytes))
		}
	}
	fmt.Println("\n\n\nSign Policy >>>>>", DocumentOBJ)

	DocumentMarshalled, errorMarshal := json.Marshal(DocumentOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("DocumentMarshalled marshalled successful !!!")

	return shim.Success(DocumentMarshalled)
}

func (t *CoreChainCode) updateDocument(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("document: %v", args)

	Key := sanitize(args[0], "string").(string)
	OrgCode := sanitize(args[1], "string").(string)
	OrgName := sanitize(args[2], "string").(string)
	UserId := sanitize(args[3], "string").(string)
	UserName := sanitize(args[4], "string").(string)
	Datetime := sanitize(args[5], "string").(string)
	PackageNo := sanitize(args[6], "string").(string)
	Name := sanitize(args[7], "string").(string)
	Hash := sanitize(args[8], "string").(string)
	Extension := sanitize(args[9], "string").(string)
	DocumentType := sanitize(args[10], "string").(string)
	var options []string
	// 1- check if Package Exist
	PackageOBJ := Package{}
	PackageCheck, errorFetch := fetchData(stub, PackageNo, "package")
	if errorFetch != nil {
		fmt.Println("\n Package fetch error>>", errorFetch)
		options = append(options, "Package")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package fetch error"}, "90700004", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageCheck == nil {
		fmt.Println("Package not found")
		options = append(options, "Package", PackageNo)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package not found"}, "90700005", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Package exist")

	errorUnmarshal := json.Unmarshal(PackageCheck, &PackageOBJ)
	if errorUnmarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("\n\n\nPackage unmarshal >>>>>", PackageOBJ)

	//2- action Authorization
	if PackageOBJ.RequestedBy.OrgCode != OrgCode || PackageOBJ.RequestedBy.UserId != UserId {
		fmt.Println("You're not allowed to perform this action")
		options = append(options)
		errBytes, _ := prepareErrorCode("2030", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"You're not allowed to perform this action"}, "90700009", false, true)
		// return shim.Error(string(errBytes))
	}

	if PackageOBJ.Status == "Rejected" {
		fmt.Println("Package is Rejected")
		options = append(options)
		errBytes, _ := prepareErrorCode("2031", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package is already Rejected"}, "90700007", false, true)
		// return shim.Error(string(errBytes))

	} else if PackageOBJ.Status == "Completed" {
		fmt.Println("Package is completed.")
		options = append(options)
		errBytes, _ := prepareErrorCode("2032", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Package is already completed."}, "90700008", false, true)
		// return shim.Error(string(errBytes))

	}

	fmt.Println("Package is In Progress.", PackageOBJ.Status)

	//3- fetch Document and check if it exists or not
	DocumentOBJ := Document{}
	DocumentCheck, errorFetch := fetchData(stub, Key, "document")
	if errorFetch != nil {
		fmt.Println("\n Document fetch error>>", errorFetch)
		options = append(options, "Document")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document fetch error"}, "90600003", false, true)
		// return shim.Error(string(errBytes))

	} else if DocumentCheck == nil {
		fmt.Println("Document not found")
		options = append(options, "Document", Key)
		errBytes, _ := prepareErrorCode("2002", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Document not found"}, "90600004", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Document exist")

	errorUnmarshalDoc := json.Unmarshal(DocumentCheck, &DocumentOBJ)
	if errorUnmarshalDoc != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("\n\n\nDocument Unmarshal >>>>>", DocumentOBJ)

	//3.5 Check DocumentType is Same or not.
	if DocumentType != DocumentOBJ.DocumentType {
		fmt.Println("Un matched Document Type")
		options = append(options)
		errBytes, _ := prepareErrorCode("2035", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmatched DocumentType"}, "90700010", false, true)
		// return shim.Error(string(errBytes))
	}
	//10- save Previous Current Version Hash for Package Progress
	splitKey := PackageNo + "-"
	PrevHash := strings.Split(DocumentOBJ.Key, splitKey)[1]

	//4- Create new copied Document
	NewDocumentOBJ := Document{}
	NewDocumentOBJ.Key = DocumentOBJ.Key + "_" + strconv.Itoa(DocumentOBJ.Version)
	NewDocumentOBJ.Version = DocumentOBJ.Version
	NewDocumentOBJ.Status = DocumentOBJ.Status
	NewDocumentOBJ.Progress = DocumentOBJ.Progress
	NewDocumentOBJ.DocumentName = "document"
	NewDocumentOBJ.Name = DocumentOBJ.Name
	NewDocumentOBJ.Extension = DocumentOBJ.Extension
	NewDocumentOBJ.DocumentType = DocumentOBJ.DocumentType
	NewDocumentOBJ.SignaturePolicy = DocumentOBJ.SignaturePolicy
	NewDocumentOBJ.Sequence = DocumentOBJ.Sequence
	NewDocumentOBJ.CurrentVersion = DocumentOBJ.CurrentVersion
	NewDocumentOBJ.History = DocumentOBJ.History
	NewDocumentOBJ.Signatures = DocumentOBJ.Signatures

	//5- Insert New Document
	NewDocumentMarshalled, errorMarshal := json.Marshal(NewDocumentOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("NewDocumentMarshalled marshalled successful !!!", NewDocumentOBJ.Key)

	errorInsert := insertData(&stub, NewDocumentOBJ.Key, "document", NewDocumentMarshalled)
	if errorInsert != nil {
		fmt.Println("Error occurred while inserting the Document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in Document Copy Collection!")

	//6- Update Existing Document
	NewSignatures := []Signatures{}
	for s := 0; s < len(DocumentOBJ.Signatures); s++ {
		DocumentOBJ.Signatures[s].Action = "Pending"
		DocumentOBJ.Signatures[s].Datetime = ""
		NewSignatures = append(NewSignatures, DocumentOBJ.Signatures[s])
	}

	DocumentOBJ.Name = Name
	DocumentOBJ.Extension = Extension
	DocumentOBJ.DocumentType = DocumentType
	DocumentOBJ.Version = DocumentOBJ.Version + 1
	DocumentOBJ.CurrentVersion = History{
		Type:       "Uploaded",
		OrgCode:    OrgCode,
		OrgName:    OrgName,
		UserId:     UserId,
		UserName:   UserName,
		Datetime:   Datetime,
		Hash:       Hash,
		SequenceNo: 0,
		Status:     "Uploaded",
	}
	DocumentOBJ.History = []History{}
	DocumentOBJ.Signatures = NewSignatures
	DocumentOBJ.DocumentHistory = append(DocumentOBJ.DocumentHistory, DocumentHistory{
		DocumentId: NewDocumentOBJ.Key,
		Version:    NewDocumentOBJ.Version,
	})

	DocumentMarshalled, errorMarshal := json.Marshal(DocumentOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("DocumentMarshalled marshalled successful !!!", DocumentOBJ.Key)
	errorInsertOld := insertData(&stub, DocumentOBJ.Key, "document", DocumentMarshalled)
	if errorInsertOld != nil {
		fmt.Println("Error occurred while inserting the Document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in Document Collection!")
	fmt.Println("PACKAGE DOCUMENTS >>>", PackageOBJ.Documents)
	for i := 0; i < len(PackageOBJ.Documents); i++ {
		fmt.Println("\n\n\n\nComparision!", PackageOBJ.Documents[i].Hash, PrevHash)

		if PackageOBJ.Documents[i].Hash == PrevHash {
			fmt.Println("Matched")
			PackageOBJ.Documents[i].Progress = 0
			PackageOBJ.Documents[i].Name = Name
		}
	}

	//upsert
	PackageMarshalled, errorMarshal := json.Marshal(PackageOBJ)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))
	}
	errorInsert2 := insertData(&stub, PackageNo, "package", PackageMarshalled)
	if errorInsert2 != nil {
		fmt.Println("Error occurred while inserting the Package in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	fmt.Println("Upsert Successfully in Package Collection!")

	return shim.Success(nil)
}

func (t *CoreChainCode) verifyDocument(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign document: %v", args)
	var options []string
	//if no args passed
	if len(args[0]) <= 0 {
		options = append(options)
		errBytes, _ := prepareErrorCode("2009", options)
		return shim.Error(string(errBytes))
		//	return shim.Error("Invalid Argument")
	}

	//Initialization the Structure
	var documentHashStructure DocumentHash
	var documentStructure Document

	//collections
	var documentHashCollection = "documentHash"
	var documentCollection = "document"

	hash := sanitize(args[0], "string").(string)
	status := ""
	verison := 0
	tranxHistory := []History{}
	documentDetails := History{}

	//check document hash exisiting record
	documentHashData, err := fetchData(stub, hash, documentHashCollection)
	if err != nil {
		options = append(options, "Document Hash Structure")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Error while fetching data from Document Hash Structure " + err.Error())
	}
	//if no document found
	if documentHashData == nil {
		status = "Not Verified"
		tranxHistory = nil
	} else {
		//unmarshal documentHashDataData
		err = json.Unmarshal([]byte(documentHashData), &documentHashStructure)
		if err != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Error while unmarshal documentHashData " + err.Error())
		}

		//check document exisiting record
		documentData, err := fetchData(stub, documentHashStructure.PackageId, documentCollection)
		if err != nil {
			options = append(options, "Document Structure")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Error while fetching data from Document Structure " + err.Error())
		}
		if documentData == nil {
			options = append(options, "Document Structure", documentHashStructure.PackageId)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
			//return shim.Error("No data found with key: " + documentHashStructure.PackageId)
		}

		//unmarshal documentData
		err = json.Unmarshal([]byte(documentData), &documentStructure)
		if err != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Error while unmarshal documentData " + err.Error())
		}

		//set document version
		verison = documentStructure.Version

		// if request hash is equal to current version hash then doc is genuine
		if documentStructure.CurrentVersion.Hash == hash {
			status = "Genuine"
			documentDetails = documentStructure.CurrentVersion
			tranxHistory = append(documentStructure.History, documentDetails)
		} else {
			//if hash is not available as current hash and not in history so it will considered as not verified
			if len(documentStructure.History) == 0 {
				status = "Not Verified"
				tranxHistory = nil
			} else {
				status = "Stale"
				for i := 0; i < len(documentStructure.History); i++ {
					//if doc is not genuine then find hash in doc history, if found then doc is stale
					if documentStructure.History[i].Hash == hash {
						documentDetails = documentStructure.History[i]
					} else {
						tranxHistory = append(tranxHistory, documentStructure.History[i])
					}
				}
				//append current version as well
				tranxHistory = append(tranxHistory, documentStructure.CurrentVersion)
			}
		}
	}

	documentVerifyDto := DocumentVerifyDto{
		Status:          status,
		Version:         verison,
		DocumentDetails: documentDetails,
		History:         tranxHistory,
	}

	data, errorMarshal := json.Marshal(documentVerifyDto)
	if errorMarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Error while Marshalling documentVerifyDto -----> " + errorMarshal.Error())
	}

	return shim.Success(data)
}

func (t *CoreChainCode) verifyDocument2(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("Args ===> %v", args)
	var options []string
	//if no args passed
	if len(args[0]) <= 0 {
		options = append(options)
		errBytes, _ := prepareErrorCode("2009", options)
		return shim.Error(string(errBytes))
	}

	mode := sanitize(args[0], "string").(string)
	data := sanitize(args[1], "string").(string)
	hash := mode + "_" + data
	var verifyResponse VerifyResponse

	//Initialization the Structure
	var documentMappingStructure DocumentMapping
	var documentStructure Document

	//collections
	var documentMappingCollection = "documentMapping"
	var documentCollection = "document"

	//check document hash exisiting record
	documentMappingData, err := fetchData(stub, hash, documentMappingCollection)
	if err != nil {
		options = append(options, "Document Mapping Structure")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
	}
	//if no document found
	if documentMappingData == nil {
		fmt.Println("No data found with key: " + hash)
		verifyResponse.Status = "Not Verified"
	} else {
		//unmarshal documentMappingData
		err = json.Unmarshal([]byte(documentMappingData), &documentMappingStructure)
		if err != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
		}

		//check document exisiting record
		documentData, err := fetchData(stub, documentMappingStructure.DocumentKey, documentCollection)
		if err != nil {
			options = append(options, "Document Structure")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
		}
		if documentData == nil {
			options = append(options, "Document Structure", documentMappingStructure.DocumentKey)
			errBytes, _ := prepareErrorCode("2002", options)
			return shim.Error(string(errBytes))
		}

		//unmarshal documentData
		err = json.Unmarshal([]byte(documentData), &documentStructure)
		if err != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
		}

		verifyResponse.Status = documentStructure.Status
		if verifyResponse.Status == "In Progress" {
			for _, signature := range documentStructure.Signatures {
				if signature.IsDispatch {
					verifyResponse.Status = signature.Action
					break
				}
			}
		}
		// Creating History
		var documentHistory []History
		if verifyResponse.Status == "Signed" {
			history := documentStructure.History
			for i := 0; i < len(history); i++ {
				if history[i].Status != "Uploaded" {
					documentHistory = append(documentHistory, history[i])
				}
			}
			documentHistory = append(documentHistory, documentStructure.CurrentVersion)
		}
		verifyResponse.History = documentHistory
	}

	fmt.Println("Response =>>>> ", verifyResponse)

	response, errorMarshal := json.Marshal(verifyResponse)
	if errorMarshal != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Error while Marshalling documentVerifyDto -----> " + errorMarshal.Error())
	}

	return shim.Success(response)
}

func (t *CoreChainCode) addUpdateUser(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign document: %v", args)
	var options []string
	//if no args passed
	if len(args[0]) <= 0 {
		options = append(options)
		errBytes, _ := prepareErrorCode("2009", options)
		return shim.Error(string(errBytes))
		//	return shim.Error("Invalid Argument")
	}

	//Initialization the Structure
	var orgUser OrgUser

	//collections
	var orgUserCollection = "orgUsers"

	orgType_groups_key := sanitize(args[0], "string").(string)
	userid := sanitize(args[1], "string").(string)
	email := sanitize(args[2], "string").(string)
	status := sanitize(args[3], "bool").(bool)

	//check document hash exisiting record
	orgTypebytes, err := fetchData(stub, orgType_groups_key, orgUserCollection)
	if err != nil {
		options = append(options, "Document Hash Structure")
		errBytes, _ := prepareErrorCode("2001", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Error while fetching data from Document Hash Structure " + err.Error())
	}
	//if no document found
	if orgTypebytes != nil {
		err = json.Unmarshal([]byte(orgTypebytes), &orgUser)
		if err != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Error while unmarshal documentHashData " + err.Error())
		}
	} else {

		orgUser = OrgUser{
			DocumentName: "orguser",
			Key:          orgType_groups_key,
			GroupUsers:   make(map[string]Users),
		}

	}

	if orgUser.GroupUsers == nil {

		orgUser.GroupUsers = make(map[string]Users)
	}

	if !status {

		//delete(orgUser.GroupUsers, userid)
	} else {

		orgUser.GroupUsers[userid] = Users{Email: email}
	}

	orgusermarshalled, errorMarshal := json.Marshal(orgUser)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}
	fmt.Println("DocumentMarshalled marshalled successful !!!", orgUser.Key)
	errorInsertOld := insertData(&stub, orgType_groups_key, orgUserCollection, orgusermarshalled)
	if errorInsertOld != nil {
		fmt.Println("Error occurred while inserting the Document in private collection")
		options = append(options)
		errBytes, _ := prepareErrorCode("2004", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Error while Insertion"}, "90100002", false, true)
		// return shim.Error(string(errBytes))
	}

	return shim.Success(nil)
}

func (t *CoreChainCode) getOrgTypeGroups(stub hypConnect, args []string, functionName string) pb.Response {
	fmt.Printf("sign document: %v", args)
	var options []string
	//if no args passed
	if len(args[0]) <= 0 {
		options = append(options)
		errBytes, _ := prepareErrorCode("2009", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Invalid Argument")
	}

	//Initialization the Structure
	//var orgUser OrgUser

	//collections
	var orgUserCollection = "orgUsers"

	//orgType_groups_keys := sanitize(args[0], "string").(string)

	var orggroupkeys []string
	err := json.Unmarshal([]byte(args[0]), &orggroupkeys)
	if err != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		//return shim.Error("Error while unmarshal documentHashData " + err.Error())
	}

	if err != nil {
		options = append(options)
		errBytes, _ := prepareErrorCode("2000", options)
		return shim.Error(string(errBytes))
		//	return shim.Error("Error while unmarshal documentHashData " + err.Error())
	}

	var email []string

	for i := 0; i < len(orggroupkeys); i++ {
		orgTypebytes, err := fetchData(stub, orggroupkeys[i], orgUserCollection)
		var orguser OrgUser
		if err != nil {
			options = append(options, "Document Hash Structure")
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			//return shim.Error("Error while fetching data from Document Hash Structure " + err.Error())
		}
		//if no document found
		if orgTypebytes == nil {
			options = append(options, "Group", orggroupkeys[i])
			errBytes, _ := prepareErrorCode("2001", options)
			return shim.Error(string(errBytes))
			//return shim.Error("No data found for this group " + orggroupkeys[i])

		}

		err = json.Unmarshal([]byte(orgTypebytes), &orguser)

		if err != nil {
			options = append(options)
			errBytes, _ := prepareErrorCode("2000", options)
			return shim.Error(string(errBytes))
			//	return shim.Error("An error occurred while orguser: " + err.Error())
		}

		for key, element := range orguser.GroupUsers {

			fmt.Println("Key:", key, "=>", "Element:", element)

			email = append(email, element.Email)
		}

	}

	emailmarshalled, errorMarshal := json.Marshal(email)
	if errorMarshal != nil {
		fmt.Println("Error occurred while marshalling the data")
		options = append(options)
		errBytes, _ := prepareErrorCode("2003", options)
		return shim.Error(string(errBytes))
		// errBytes := genericErrorHandler(nil, []string{"Unmarshalling error"}, "90100001", false, true)
		// return shim.Error(string(errBytes))

	}

	return shim.Success(emailmarshalled)
}
