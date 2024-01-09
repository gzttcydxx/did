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
	DID, err := models.Parse(did)
	readIdentity, err := s.ReadIdentity(ctx, did)

	if readIdentity != nil && readIdentity.Did == did {
		return fmt.Errorf("the identity %s already exists", did)
	}

	identityJSON, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(did, identityJSON)
}

func (s *SmartContract) ReadIdentity(ctx contractapi.TransactionContextInterface, did string) (*Identity, error) {
	identityJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, err
	}
	if identityJSON == nil {
		return nil, nil
	}

	var identity Identity
	err = json.Unmarshal(identityJSON, &identity)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (s *SmartContract) UpdateIdentity(ctx contractapi.TransactionContextInterface, did string, didDoc DidDoc) error {
	readIdentity, err := s.ReadIdentity(ctx, did)

	if readIdentity == nil {
		return fmt.Errorf("the identity %s does not exist", did)
	}

	identity := Identity{
		Did:    did,
		DidDoc: didDoc,
	}

	identityJSON, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(did, identityJSON)
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
