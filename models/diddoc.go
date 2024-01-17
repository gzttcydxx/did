package models

import (
	"fmt"
	"time"
)

type Context interface{}

type DIDDoc struct {
	Context              Context
	ID                   string
	AlsoKnownAs          []string
	VerificationMethod   []VerificationMethod
	Authentication       []Verification
	AssertionMethod      []Verification
	CapabilityDelegation []Verification
	CapabilityInvocation []Verification
	KeyAgreement         []Verification
	Created              *time.Time
	Updated              *time.Time
	Proof                []Proof
}

type rawDIDDoc struct {
	Context              Context                  `json:"@context,omitempty"`
	ID                   string                   `json:"id,omitempty"`
	AlsoKnownAs          []interface{}            `json:"alsoKnownAs,omitempty"`
	VerificationMethod   []map[string]interface{} `json:"verificationMethod,omitempty"`
	Authentication       []interface{}            `json:"authentication,omitempty"`
	AssertionMethod      []interface{}            `json:"assertionMethod,omitempty"`
	CapabilityDelegation []interface{}            `json:"capabilityDelegation,omitempty"`
	CapabilityInvocation []interface{}            `json:"capabilityInvocation,omitempty"`
	KeyAgreement         []interface{}            `json:"keyAgreement,omitempty"`
	Created              *time.Time               `json:"created,omitempty"`
	Updated              *time.Time               `json:"updated,omitempty"`
	Proof                []interface{}            `json:"proof,omitempty"`
}

// UnmarshalJSON unmarshals a DID Document.
func (doc *DIDDoc) UnmarshalJSON(data []byte) error {
	_doc, err := ParseDocument(data)
	if err != nil {
		return fmt.Errorf("failed to parse did doc: %w", err)
	}

	*doc = *_doc

	return nil
}
