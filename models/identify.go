package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type DID struct {
	Scheme           string
	Method           string
	MethodSpecificID string
}

func (d *DID) String() string {
	return d.Scheme + ":" + d.Method + ":" + d.MethodSpecificID
}

func Parse(did string) (*DID, error) {
	const idchar = `a-zA-Z0-9-_\.`
	regex := fmt.Sprintf(`^did:[a-z0-9]+:(:+|[:%s]+)*[%%:%s]+[^:]$`, idchar, idchar)

	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex=%s (this should not have happened!). %w", regex, err)
	}

	if !r.MatchString(did) {
		return nil, fmt.Errorf(
			"invalid did: %s. Make sure it conforms to the DID syntax: https://w3c.github.io/did-core/#did-syntax", did)
	}

	parts := strings.SplitN(did, ":", 3)

	return &DID{
		Scheme:           "did",
		Method:           parts[1],
		MethodSpecificID: parts[2],
	}, nil
}

type DIDURL struct {
	DID
	Path     string
	Queries  map[string][]string
	Fragment string
}

func ParseDIDURL(didURL string) (*DIDURL, error) {
	split := strings.IndexAny(didURL, "?/#")

	didPart := didURL
	pathQueryFragment := ""

	if split != -1 {
		didPart = didURL[:split]
		pathQueryFragment = didURL[split:]
	}

	retDID, err := Parse(didPart)
	if err != nil {
		return nil, err
	}

	if pathQueryFragment == "" {
		return &DIDURL{
			DID:     *retDID,
			Queries: map[string][]string{},
		}, nil
	}

	hasPath := pathQueryFragment[0] == '/'

	if !hasPath {
		pathQueryFragment = "/" + pathQueryFragment
	}

	urlParts, err := url.Parse(pathQueryFragment)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path, query, and fragment components of DID URL: %w", err)
	}

	ret := &DIDURL{
		DID:      *retDID,
		Queries:  urlParts.Query(),
		Fragment: urlParts.Fragment,
	}

	if hasPath {
		ret.Path = urlParts.Path
	}

	return ret, nil
}

type Context interface{}

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

type Verification struct {
	VerificationMethod VerificationMethod       `json:"verificationMethod,omitempty"`
	Relationship       VerificationRelationship `json:"relationship,omitempty"`
	Embedded           bool                     `json:"embedded,omitempty"`
}

type Proof struct {
	Type         string
	Created      *time.Time
	Creator      string
	ProofValue   []byte
	Domain       string
	Nonce        []byte
	ProofPurpose string
	relativeURL  bool
}

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

func ParseDocument(data []byte) (*DIDDoc, error) {
	var raw rawDIDDoc

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal did doc: %w", err)
	}

	doc := &DIDDoc{
		Context:              raw.Context,
		ID:                   raw.ID,
		AlsoKnownAs:          any2Array(raw.AlsoKnownAs),
		Authentication:       any2Verification(raw.Authentication),
		AssertionMethod:      any2Verification(raw.AssertionMethod),
		CapabilityDelegation: any2Verification(raw.CapabilityDelegation),
		CapabilityInvocation: any2Verification(raw.CapabilityInvocation),
		KeyAgreement:         any2Verification(raw.KeyAgreement),
		Created:              raw.Created,
		Updated:              raw.Updated,
		Proof:                any2Proof(raw.Proof),
	}

	for _, vm := range raw.VerificationMethod {
		verificationMethod, err := parseVerificationMethod(vm)
		if err != nil {
			return nil, err
		}

		doc.VerificationMethod = append(doc.VerificationMethod, *verificationMethod)
	}

	return doc, nil
}

func parseVerificationMethod(vm map[string]interface{}) (*VerificationMethod, error) {
	id := any2String(vm["id"])
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	relativeURL := false

	if strings.HasPrefix(id, "#") {
		relativeURL = true
	}

	return &VerificationMethod{
		ID:                 id,
		Type:               any2String(vm["type"]),
		Controller:         any2String(vm["controller"]),
		Value:              any2Bytes(vm["value"]),
		relativeURL:        relativeURL,
		PublicKeyMultibase: any2String(vm["publicKeyMultibase"]),
	}, nil
}

func any2Verification(i interface{}) []Verification {
	if i == nil {
		return nil
	}

	is, ok := i.([]interface{})
	if !ok {
		return nil
	}

	var result []Verification

	for _, e := range is {
		if e != nil {
			result = append(result, any2VerificationElement(e))
		}
	}

	return result
}

func any2VerificationElement(i interface{}) Verification {
	switch e := i.(type) {
	case map[string]interface{}:
		vm, err := parseVerificationMethod(e)
		if err != nil {
			return Verification{}
		}

		return Verification{
			VerificationMethod: *vm,
		}
	case string:
		return Verification{
			VerificationMethod: VerificationMethod{
				ID: e,
			},
		}
	default:
		return Verification{}
	}
}

func any2Proof(i interface{}) []Proof {
	if i == nil {
		return nil
	}

	is, ok := i.([]interface{})
	if !ok {
		return nil
	}

	var result []Proof

	for _, e := range is {
		if e != nil {
			result = append(result, any2ProofElement(e))
		}
	}

	return result
}

func any2ProofElement(i interface{}) Proof {
	switch e := i.(type) {
	case map[string]interface{}:
		return Proof{
			Type:         any2String(e["type"]),
			Created:      any2Time(e["created"]),
			Creator:      any2String(e["creator"]),
			ProofValue:   any2Bytes(e["proofValue"]),
			Domain:       any2String(e["domain"]),
			Nonce:        any2Bytes(e["nonce"]),
			ProofPurpose: any2String(e["proofPurpose"]),
			relativeURL:  false,
		}
	default:
		return Proof{}
	}
}

func any2Time(i interface{}) *time.Time {
	if i == nil {
		return nil
	}

	switch e := i.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, e)
		if err != nil {
			return nil
		}

		return &t
	default:
		return nil
	}
}

func any2Bytes(i interface{}) []byte {
	if i == nil {
		return nil
	}

	switch e := i.(type) {
	case string:
		return []byte(e)
	default:
		return nil
	}
}

func any2String(i interface{}) string {
	if i == nil {
		return ""
	}

	if e, ok := i.(string); ok {
		return e
	}

	return ""
}

func any2Array(i interface{}) []string {
	if i == nil {
		return nil
	}

	is, ok := i.([]interface{})
	if !ok {
		return nil
	}

	var result []string

	for _, e := range is {
		if e != nil {
			result = append(result, any2String(e))
		}
	}

	return result
}
