package models

type VerificationMethod struct {
	ID                 string `json:"id,omitempty"`
	Type               string `json:"type,omitempty"`
	Controller         string `json:"controller,omitempty"`
	Value              []byte `json:"value,omitempty"`
	relativeURL        bool   `json:"-"`
	PublicKeyMultibase string `json:"publicKeyMultibase,omitempty"`
	// multibaseEncoding multibase.Encoding
}

type VerificationRelationship int

const (
	VerificationRelationshipGeneral VerificationRelationship = iota

	Authentication

	AssertionMethod

	CapabilityDelegation

	CapabilityInvocation

	KeyAgreement
)

type Verification struct {
	VerificationMethod VerificationMethod       `json:"verificationMethod,omitempty"`
	Relationship       VerificationRelationship `json:"relationship,omitempty"`
	Embedded           bool                     `json:"embedded,omitempty"`
}
