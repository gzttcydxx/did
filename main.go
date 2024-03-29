/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/gzttcydxx/did/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating did chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting did chaincode: %v", err)
	}
}
