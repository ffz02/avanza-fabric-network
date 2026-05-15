package main

type MSPList struct {
	OrgType string `json:"orgType"`
	MSP     string `json:"MSP"`
	ID      string `json:"ID"`
}

type errCode struct {
	ErrorCode string   `json:"errorCode"`
	Options   []string `json:"options"`
}

type global_typedatavalues struct {
	Key        string                 `json:"key"`
	Value      string                 `json:"value"`
	DataObject map[string]interface{} `json:"dataObject"`
	IsActive   bool                   `json:"isActive"`
}

type TypeDataValidateStructure struct {
	Path         string                 `json:"path"`
	TypeDataName string                 `json:"typeDataName"`
	Key          string                 `json:"key"`
	FromType     string                 `json:"fromType"`
	Exist        bool                   `json:"exist"`
	Value        string                 `json:"value"`
	DataObject   map[string]interface{} `json:"dataObjects"`
	IsActive     bool                   `json:"isActive"`
}

/* ===================================================================================
	This is the PostDataToBlockchainRegAuth
  ===================================================================================*/

// ###################### Start ########################
// type SignaturePolicy struct {
// 	Key                    string                   `json:"key"`
// 	DocumentName           string                   `json:"documentName"`
// 	PolicyName             string                   `json:"policyName"`
// 	HasOrder               bool                     `json:"hasOrder"`
// 	StampOrganizations     []StampSignOrganizations `json:"stampOrganizations"`
// 	SignatureOrganizations []StampSignOrganizations `json:"signatureOrganizations"`
// 	// Organizations []Organization `json:"organizations"`
// }

type SignaturePolicy struct {
	Key                            string               `json:"key"`
	DocumentName                   string               `json:"documentName"`
	PolicyName                     string               `json:"policyName"`
	Sequence                       bool                 `json:"sequence"`
	FinalSignatoryOrganizationName string               `json:"finalSignatoryOrgName"`
	FinalSignatoryOrganizationCode string               `json:"finalSignatoryOrgCode"`
	FinalSignatoryOrganizationType string               `json:"finalSignatoryOrgtype"`
	FinalSignatoryGroup            []string             `json:"finalSignatoryGroup"`
	Organizations                  []PolicyOrganization `json:"organizations"`
}

type OverRide struct {
	OrgType string `json:"orgType"`
	OrgCode string `json:"orgCode"`
}

type PolicyOrganization struct {
	OrganizationName string   `json:"orgName"`
	OrganizationCode string   `json:"orgCode"`
	OrganizationType string   `json:"orgtype"`
	Group            []string `json:"group"`
	SequenceNo       int      `json:"sequenceNo"`
	SLA              float64  `json:"sla"`
	Pages            []Pages  `json:"pages"`
	IsReviewer       bool     `json:"isReviewer"`
	IsFlexibleCords  bool     `json:"isFlexibleCords"`
	DispatchStatus   bool     `json:"dispatchStatus"`
}

type Pages struct {
	PageNo int     `json:"pageNo"`
	Type   string  `json:"type"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	W      float64 `json:"w"`
	H      float64 `json:"h"`
}

type StampSignOrganizations struct {
	OrderNo          int     `json:"orderNo"`
	OrganizationName string  `json:"organizationName"`
	OrganizationCode string  `json:"organizationCode"`
	Pages            []Pages `json:"pages"`
}

type User struct {
	Key          string               `json:"key"`
	DocumentName string               `json:"documentName"`
	UserId       string               `json:"userId"`
	Email        string               `json:"email"`
	OrgCode      string               `json:"orgCode"`
	Current      CurrentHistorySign   `json:"current"`
	History      []CurrentHistorySign `json:"history"`
}

// type UserReq struct {
// 	Key       string `json:"key"`
// 	UserId    string `json:"userId"`
// 	OrgCode   string `json:"orgCode"`
// 	Signature string `json:"signature"`
// 	Publickey string `json:"publickey"`
// }

type CurrentHistorySign struct {
	Type              string `json:"type"`
	Signature         string `json:"signature"`
	CertificateType   string `json:"certificateType"`
	P12Certificate    string `json:"p12Certificate"`
	PublicCertificate string `json:"publicCertificate"`
	Password          string `json:"password"`
	Datetime          string `json:"datetime"`
}

type Organization struct {
	Key          string                `json:"key"`
	DocumentName string                `json:"documentName"`
	OrgCode      string                `json:"orgCode"`
	OrgName      string                `json:"orgName"`
	GroupIDs     []string              `json:"groupIDs"`
	Current      CurrentHistoryStamp   `json:"current"`
	History      []CurrentHistoryStamp `json:"history"`
}

type OrganizationReq struct {
	Key       string `json:"key"`
	OrgCode   string `json:"orgCode"`
	OrgName   string `json:"orgName"`
	Stamp     string `json:"stamp"`
	Publickey string `json:"publickey"`
}

type CurrentHistoryStamp struct {
	Type              string `json:"type"`
	Stamp             string `json:"stamp"`
	P12Certificate    string `json:"p12Certificate"`
	PublicCertificate string `json:"publicCertificate"`
	CertificateType   string `json:"certificateType"`
	Password          string `json:"password"`
	Datetime          string `json:"datetime"`
}

type DocumentType struct {
	Key               string `json:"key"`
	DocumentName      string `json:"documentName"`
	DocumentType      string `json:"documentType"`
	Label             string `json:"label"`
	Value             string `json:"value"`
	SampleDoc         string `json:"sampleDoc"`
	SignaturePolicyId string `json:"signaturePolicyId"`
	LetterTemplateId  string `json:"letterTemplateId"`
	DocumentSource    string `json:"documentSource"`
	DispatchType      string `json:"dispatchType"`
	DispatchValue     string `json:"dispatchValue"`
	EmailTemplate     string `json:"emailTemplate"`
}

type PackageType struct {
	Key          string              `json:"key"`
	DocumentName string              `json:"documentName"`
	PackageType  string              `json:"packageType"`
	Documents    []PackageDocument   `json:"documents"`
	NotifyPolicy PackageNotifyPolicy `json:"notifyPolicy"`
}

type PackageDocument struct {
	DocumentType    string `json:"documentType"`
	SignaturePolicy string `json:"signaturePolicy"`
}

type PackageNotifyPolicy struct {
	RejectAllInvolved  bool          `json:"rejectAllInvolved"`
	ApproveAllInvolved bool          `json:"approveAllInvolved"`
	RejectParties      []NotifyParty `json:"rejectParties"`
	ApproveParties     []NotifyParty `json:"approveParties"`
}

type NotifyParty struct {
	OrgType   string `json:"orgType"`
	OrgCode   string `json:"orgCode"`
	GroupCode string `json:"groupCode"`
	GroupName string `json:"groupName"`
}

type Package struct {
	Key          string `json:"key"`
	DocumentName string `json:"documentName"`
	//Mode        	string      		`json:"mode"`   only required in payload
	PackageType     string              `json:"packageType"`
	SignaturePolicy string              `json:"signaturePolicy"`
	PhysicalStatus  string              `json:"physicalStatus"`
	PackageNumber   string              `json:"packageNumber"`
	PackageName     string              `json:"packageName"`
	Progress        float64             `json:"progress"`
	Status          string              `json:"status"`
	Refno           string              `json:"Refno"`
	RequestedBy     RequestedBy         `json:"requestedBy"`
	RequestedOn     string              `json:"requestedOn"`
	Documents       []UploadedDocuments `json:"documents"`
	OwnerOrg        string              `json:"ownerOrg"`
	OtherDocuments  []OtherDocuments    `json:"otherDocuments"`
	OverrideOrgs    []OverrideOrgs      `json:"overrideOrgs"`
	IsRegenrated    bool                `json:"isRegenrated"`
	IsDispatch      bool                `json:"isDispatch"`
}

type OverrideOrgs struct {
	OrgType string `json:"orgType"`
	OrgCode string `json:"orgCode"`
}

type OtherDocuments struct {
	DocumentName string `json:"documentName"`
	PackageId    string `json:"packageId"`
	DocumentId   string `json:"documentId"`
	DocumentType string `json:"documentType"`
}
type DocumentMapping struct {
	Key         string `json:"key"`
	DocumentKey string `json:"documentKey"`
}
type UploadedDocuments struct {
	Name                           string   `json:"name"`
	Hash                           string   `json:"hash"`
	Extension                      string   `json:"extension"`
	GUID                           string   `json:"guid"`
	DocumentType                   string   `json:"documentType"`
	SignaturePolicy                string   `json:"signaturePolicy"`
	DocumentPolicy                 string   `json:"documentPolicy"`
	OptionalOrgs                   []string `json:"optionalOrgs"`
	SignatureOrgs                  []string `json:"signatureOrgs"`
	SignatureGroups                []string `json:"signatureGroups"`
	Progress                       float64  `json:"progress"`
	IsReadableOnly                 bool     `json:"isReadableOnly"`
	PackageId                      string   `json:"packageId"`
	SetNo                          string   `json:"setNo"`
	FinalSignatoryOrganizationName string   `json:"finalSignatoryOrgName"`
	FinalSignatoryOrganizationCode string   `json:"finalSignatoryOrgCode"`
	FinalSignatoryOrganizationType string   `json:"finalSignatoryOrgtype"`
	FinalSignatoryGroup            []string `json:"finalSignatoryGroup"`
	QRHash                         string   `json:"QRHash"`
	DispatchEmail                  string   `json:"dispatchEmail"`
	EmailTemplate                  string   `json:"emailTemplate"`
	HolderName                     string   `json:"holderName"`
}

type RequestedBy struct {
	OrgCode  string `json:"orgCode"`
	OrgName  string `json:"orgName"`
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}

type Document struct {
	Key                            string            `json:"key"`
	DocumentName                   string            `json:"documentName"`
	Name                           string            `json:"name"`
	Extension                      string            `json:"extension"`
	DocumentType                   string            `json:"documentType"`
	SignaturePolicy                string            `json:"signaturePolicy"`
	DocumentPolicy                 string            `json:"documentPolicy"`
	PackageId                      string            `json:"packageid"`
	Version                        int               `json:"version"`
	Status                         string            `json:"status"`
	Progress                       float64           `json:"progress"`
	Sequence                       bool              `json:"sequence"`
	CurrentVersion                 History           `json:"currentVersion"`
	ReturnedVersions               []History         `json:"returnedVersions"`
	NextSignatory                  string            `json:"nextSignatory"`
	NominationCode                 string            `json:"nominationCode"`
	GUID                           string            `json:"guid"`
	SetNo                          string            `json:"setNo"`
	FinalSignatoryOrganizationName string            `json:"finalSignatoryOrgName"`
	FinalSignatoryOrganizationCode string            `json:"finalSignatoryOrgCode"`
	FinalSignatoryOrganizationType string            `json:"finalSignatoryOrgtype"`
	FinalSignatoryGroup            []string          `json:"finalSignatoryGroup"`
	History                        []History         `json:"history"`
	Signatures                     []Signatures      `json:"signatures"`
	DocumentHistory                []DocumentHistory `json:"documentHistory"`
	DispatchEmail                  string            `json:"dispatchEmail"`
	EmailTemplate                  string            `json:"emailTemplate"`
	HolderName                     string            `json:"holderName"`
}

type PackageDTO struct {
	Key          string `json:"key"`
	DocumentName string `json:"documentName"`
	//Mode          string              `json:"mode"`   only required in payload
	PackageType     string              `json:"packageType"`
	PackageNumber   string              `json:"packageNumber"`
	PackageName     string              `json:"packageName"`
	SignaturePolicy string              `json:"signaturePolicy"`
	PhysicalStatus  string              `json:"physicalStatus"`
	Progress        float64             `json:"progress"`
	Status          string              `json:"status"`
	RequestedBy     RequestedBy         `json:"requestedBy"`
	RequestedOn     string              `json:"requestedOn"`
	Documents       []UploadedDocuments `json:"documents"`
	DocsPackage     []Document          `json:"docsPackage"`
	PackageTypeArr  PackageType         `json:"packageTypeArr"`
	OwnerOrg        string              `json:"ownerOrg"`
	OtherDocuments  []OtherDocuments    `json:"otherDocuments"`
	OverrideOrgs    []OverrideOrgs      `json:"overrideOrgs"`
	Refno           string              `json:"Refno"`
	IsRegenrated    bool                `json:"isRegenrated"`
}

type History struct {
	Type           string `json:"type"`
	OrgCode        string `json:"orgCode"`
	OrgName        string `json:"orgName"`
	UserId         string `json:"userId"`
	UserName       string `json:"userName"`
	Datetime       string `json:"datetime"`
	Hash           string `json:"hash"`
	SequenceNo     int    `json:"sequenceNo"`
	Status         string `json:"status"`
	Comments       string `json:"comments"`
	DocumentType   string `json:"documentType"`
	TranxHash      string `json:"tranxHash"`
	BlockNo        string `json:"blockNo"`
	RejectReason   string `json:"rejectReason"`
	ReturnReason   string `json:"returnReason"`
	ReturnComments string `json:"returnComments"`
	UserSite       string `json:"userSite"`
	UserDepartment string `json:"userDepartment"`
	UserGroup      string `json:"userGroup"`
}

type DocumentHash struct {
	Key          string `json:"key"`
	DocumentName string `json:"documentName"`
	PackageId    string `json:"packageId"`
}

type Signatures struct {
	OrgCode         string   `json:"orgCode"`
	OrgName         string   `json:"orgName"`
	Type            string   `json:"type"`
	Action          string   `json:"action"`
	Datetime        string   `json:"datetime"`
	SequenceNo      int      `json:"sequenceNo"`
	SLA             float64  `json:"sla"`
	Group           []string `json:"group"`
	Pages           []Pages  `json:"pages"`
	IsReviewer      bool     `json:"isReviewer"`
	IsFlexibleCords bool     `json:"isFlexibleCords"`
	IsDispatch      bool     `json:"isDispatch"`
}

type DocumentHistory struct {
	DocumentId string `json:"documentId"`
	Version    int    `json:"version"`
}

type DocumentVerifyDto struct {
	Status          string    `json:"status"`
	Version         int       `json:"version"`
	DocumentDetails History   `json:"documentDetails"`
	History         []History `json:"history"`
}

type OrgUser struct {
	DocumentName string           `json:"documentName"`
	Key          string           `json:"key"` //orgtype_group
	GroupUsers   map[string]Users `json:"groupUsers"`
}

type Users struct {
	Email string `json:"email"`
}

type Tasks struct {
	DocumentName   string   `json:"documentName"`
	Key            string   `json:"key"`
	TaskName       string   `json:"taskName"`
	DocumentType   []string `json:"documentType"`
	OrgCode        string   `json:"orgCode"`
	GroupName      string   `json:"groupName"`
	Port           string   `json:"port"`
	Status         string   `json:"status"`
	PackageId      string   `json:"packageId"`
	NotificationId string   `json:"notificationId"`
	SNNumber       string   `json:"snNumber"`
	SLATime        int      `json:"slaTime"`
	AdditonalData  string   `json:"additonalData"`
	CompletedAt    string   `json:"completedAt"`
	Activity       string   `json:"activity"`
	Stage          string   `json:"stage"`
}

type TaskReqPayload struct {
	DocumentName   string   `json:"documentName"`
	Key            string   `json:"key"`
	TaskName       string   `json:"taskName"`
	DocumentType   []string `json:"documentType"`
	OrgCode        string   `json:"orgCode"`
	GroupName      string   `json:"groupName"`
	Port           string   `json:"port"`
	Status         string   `json:"status"`
	PackageId      string   `json:"packageId"`
	NotificationId string   `json:"notificationId"`
	SNNumber       string   `json:"snNumber"`
	TaskStatus     string   `json:"taskStatus"`
	SLATime        int      `json:"slaTime"`
	AdditonalData  string   `json:"additonalData"`
	Activity       string   `json:"activity"`
	Stage          string   `json:"stage"`
}

type TaskTypes struct {
	SIGNATURE_REQUEST    string
	PSN_UPDATE           string
	KOC_TIMINGS          string
	ACTUAL_LOAD_QUANTITY string
	GENERATE_PACKAGE     string
	FORWARD_TO_CUSTOMER  string
	PRINT_DOCUMENTS      string
	UPLOAD_FINAL_SET     string
}

type NextSignatoryInfo struct {
	OrgCode string   `json:"orgCode"`
	Group   []string `json:"groupId"`
}

type VerifyResponse struct {
	Status  string    `json:"status"`
	History []History `json:"history"`
}
