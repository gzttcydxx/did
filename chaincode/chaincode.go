package chaincode

import (
	"encoding/json"
	"fmt"

	"did/models"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/hyperledger/fabric/scripts/fabric-samples/did/models"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) CreateIdentity(ctx contractapi.TransactionContextInterface, did string) error {
	readIdentity, err := s.ReadIdentity(ctx, did)

	if readIdentity != nil {
		return fmt.Errorf("the identity %s already exists", did)
	}

	identity := models.DIDDoc{
		Context: []string{"https://www.w3.org/ns/did/v1"},
		ID:      did,
	}

	identityJSON, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(did, identityJSON)
}

func (s *SmartContract) ReadIdentity(ctx contractapi.TransactionContextInterface, did string) (*models.DIDDoc, error) {
	DIDdocJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, err
	}
	if DIDdocJSON == nil {
		return nil, nil
	}

	var DIDdoc *models.DIDDoc
	DIDdoc.UnmarshalJSON(DIDdocJSON)

	return DIDdoc, nil
}

func (s *SmartContract) DeleteIdentity(ctx contractapi.TransactionContextInterface, did string) error {
	readIdentity, err := s.ReadIdentity(ctx, did)
	if err != nil {
		return err
	}

	if readIdentity == nil {
		return fmt.Errorf("the identity %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
