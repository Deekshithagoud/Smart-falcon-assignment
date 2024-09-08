#!/bin/bash

# Exit on first error, print all commands.
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

# Clean the keystore
rm -rf ./hfc-key-store

# Launch network; create channel and join peer to channel
cd ../test-network
./network.sh down
./network.sh up createChannel -c mychannel -ca

# Set the environment variables for Org1
export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA

# Package the chaincode
peer lifecycle chaincode package asset-transfer.tar.gz --path ../asset-transfer-basic/chaincode-go --lang golang --label asset-transfer_1

# Install the chaincode
peer lifecycle chaincode install asset-transfer.tar.gz

# Get the package ID for the installed chaincode
PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep asset-transfer_1 | awk '{print $3}' | sed 's/.$//')

# Approve the chaincode for Org1
peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name asset-transfer --version 1.0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA

# Check commit readiness
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name asset-transfer --version 1.0 --sequence 1 --output json

# Commit the chaincode
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name asset-transfer --version 1.0 --sequence 1 --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles $PEER0_ORG1_CA

# Initialize the chaincode
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n asset-transfer --isInit -c '{"Args":[]}'
