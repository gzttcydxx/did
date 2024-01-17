package models

type DID struct {
	Scheme           string
	Method           string
	MethodSpecificID string
}

func (d *DID) String() string {
	return d.Scheme + ":" + d.Method + ":" + d.MethodSpecificID
}

type DIDURL struct {
	DID
	Path     string
	Queries  map[string][]string
	Fragment string
}
