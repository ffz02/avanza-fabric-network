package main

type TypeDataContractStatus struct {
	ContractRecorded        string `json:"ContractRecorded"`
	InstrumentsAssociated   string `json:"InstrumentsAssociated"`
	CustomerConsentProvided string `json:"CustomerConsentProvided"`
	Terminating             string `json:"Terminating"`
	Terminated              string `json:"Terminated"`
}
