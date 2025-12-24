package vc

import (
	"encoding/json"
	"time"
)

type CreditReadinessSubject struct {
	EnterpriseExists bool   `json:"enterprise_exists"`
	MonthsActive     int    `json:"months_active"`
	ValidationSource string `json:"validation_source"`
	Sector           string `json:"sector"`
	District         string `json:"district"`
}

type CreditReadinessCredential struct {
	Context        []string               `json:"@context"`
	Type           []string               `json:"type"`
	Issuer         string                 `json:"issuer"`
	IssuanceDate   time.Time              `json:"issuanceDate"`
	ExpirationDate time.Time              `json:"expirationDate"`
	CredentialSub  CreditReadinessSubject `json:"credentialSubject"`
}

func (c CreditReadinessCredential) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
