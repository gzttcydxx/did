package models

import "strings"

type DID struct {
	Scheme     string `json:"scheme" required:"false" default:"did"`
	Method     string `json:"method" required:"false" default:"ethr"`
	ChainID    string `json:"chain_id" required:"false" default:"1"`
	SpecificID string `json:"specific_id" required:"false" default:"0x1234567890"`
	Fragment   string `json:"fragment" required:"false" default:"#1234567890"`
}

func (d *DID) ToString() string {
	result := d.Scheme + ":" + d.Method
	if d.ChainID != "" {
		result += ":" + d.ChainID
	}
	result += ":" + d.SpecificID
	if d.Fragment != "" {
		result += "#" + d.Fragment
	}
	return result
}

func (d *DID) FromString(s string) {
	parts := strings.Split(s, "#")
	if len(parts) > 1 {
		d.Fragment = parts[1]
	}
	parts = strings.Split(parts[0], ":")
	d.Scheme = parts[0]
	d.Method = parts[1]
	if len(parts) == 3 {
		d.SpecificID = parts[2]
	} else if len(parts) >= 4 {
		d.ChainID = parts[2]
		d.SpecificID = strings.Join(parts[3:], ":")
	}
}

func NewDID(s string) (*DID, error) {
	d := &DID{}
	d.FromString(s)
	return d, nil
}
