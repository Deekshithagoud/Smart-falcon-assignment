#!/bin/bash

# Exit on first error, print all commands.
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

# Stop the test network
cd ../test-network
./network.sh down

# Remove chaincode docker images
docker rmi $(docker images dev-* -q)

# Clean the keystore
rm -rf ./hfc-key-store

# Remove the generated channel artifacts
rm -rf ./channel-artifacts
rm -rf ./crypto-config
