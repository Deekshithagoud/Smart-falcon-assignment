package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Asset represents the structure for accounts with various attributes
type Asset struct {
    DealerID     string  `json:"dealerID"`
    MSISDN       string  `json:"msisdn"`
    MPIN         string  `json:"mpin"`
    Balance      float64 `json:"balance"`
    Status       string  `json:"status"`
    TransAmount  float64 `json:"transAmount"`
    TransType    string  `json:"transType"`
    Remarks      string  `json:"remarks"`
}

// SmartContract defines the structure for the chaincode
type SmartContract struct {
    contractapi.Contract
}

// CreateAsset creates a new asset and stores it in the ledger
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID, msisdn, mpin, status, transType, remarks string, balance, transAmount float64) error {
    asset := Asset{
        DealerID:    dealerID,
        MSISDN:      msisdn,
        MPIN:        mpin,
        Balance:     balance,
        Status:      status,
        TransAmount: transAmount,
        TransType:   transType,
        Remarks:     remarks,
    }

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return fmt.Errorf("failed to marshal asset: %v", err)
    }

    return ctx.GetStub().PutState(dealerID, assetJSON)
}

// UpdateAsset updates an existing asset's attributes in the ledger
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerID, msisdn, mpin, status, transType, remarks string, balance, transAmount float64) error {
    assetJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return fmt.Errorf("failed to read asset from world state: %v", err)
    }
    if assetJSON == nil {
        return fmt.Errorf("asset %s does not exist", dealerID)
    }

    asset := Asset{}
    if err := json.Unmarshal(assetJSON, &asset); err != nil {
        return fmt.Errorf("failed to unmarshal asset JSON: %v", err)
    }

    asset.MSISDN = msisdn
    asset.MPIN = mpin
    asset.Status = status
    asset.TransType = transType
    asset.Remarks = remarks
    asset.Balance = balance
    asset.TransAmount = transAmount

    updatedAssetJSON, err := json.Marshal(asset)
    if err != nil {
        return fmt.Errorf("failed to marshal updated asset: %v", err)
    }

    return ctx.GetStub().PutState(dealerID, updatedAssetJSON)
}

// ReadAsset reads an asset from the ledger using the dealerID
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerID string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(dealerID)
    if err != nil {
        return nil, fmt.Errorf("failed to read asset from world state: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("asset %s does not exist", dealerID)
    }

    var asset Asset
    if err := json.Unmarshal(assetJSON, &asset); err != nil {
        return nil, fmt.Errorf("failed to unmarshal asset JSON: %v", err)
    }

    return &asset, nil
}

// GetAssetHistory retrieves the transaction history of an asset
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, dealerID string) ([]*Asset, error) {
    resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve history for asset %s: %v", dealerID, err)
    }
    defer resultsIterator.Close()

    var assets []*Asset
    for resultsIterator.HasNext() {
        response, err := resultsIterator.Next()
        if err != nil {
            return nil, fmt.Errorf("failed to get next asset history: %v", err)
        }

        var asset Asset
        if err := json.Unmarshal(response.Value, &asset); err != nil {
            return nil, fmt.Errorf("failed to unmarshal asset history: %v", err)
        }
        assets = append(assets, &asset)
    }

    return assets, nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
    }

    if err := chaincode.Start(); err != nil {
        log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
    }
}
