Sure! Here's a comprehensive `README.md` for your project:

---

# Hyperledger Fabric Asset Management System

## Project Description

This project implements a blockchain-based system to manage and track assets for a financial institution. The system ensures the security, transparency, and immutability of asset records while providing an efficient way to track and manage asset-related transactions and histories. The assets represent accounts with specific attributes such as DEALERID, MSISDN, MPIN, BALANCE, STATUS, TRANSAMOUNT, TRANSTYPE, and REMARKS.

## Features

- **Create Asset**: Allows the creation of assets/accounts with specified attributes.
- **Update Asset**: Enables updating the values of existing assets.
- **Read Asset**: Queries the ledger to read the details of an asset.
- **Get Asset History**: Retrieves the transaction history of an asset.

## Prerequisites

- Docker and Docker Compose
- Node.js (version 16 or higher)
- Hyperledger Fabric binaries and Docker images

## Setup

### 1. Hyperledger Fabric Test Network

#### Step 1: Clone the Hyperledger Fabric Samples

```bash
git clone https://github.com/hyperledger/fabric-samples.git
cd fabric-samples/test-network
```

#### Step 2: Generate Crypto Material

```bash
./network.sh down
./network.sh up createChannel -ca
```

#### Step 3: Generate Genesis Block

```bash
export FABRIC_CFG_PATH=$PWD/../config/
configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```

### 2. Deploy the Smart Contract

#### Step 1: Package the Chaincode

```bash
peer lifecycle chaincode package asset-transfer.tar.gz --path ./chaincode-go --lang golang --label asset-transfer_1
```

#### Step 2: Install the Chaincode

```bash
peer lifecycle chaincode install asset-transfer.tar.gz
```

#### Step 3: Approve the Chaincode Definition

```bash
peer lifecycle chaincode approveformyorg --channelID mychannel --name asset-transfer --version 1.0 --sequence 1 --package-id <PACKAGE_ID> --init-required
```

#### Step 4: Commit the Chaincode Definition

```bash
peer lifecycle chaincode commit -o orderer.example.com:7050 --channelID mychannel --name asset-transfer --version 1.0 --sequence 1 --init-required
```

#### Step 5: Initialize the Chaincode

```bash
peer chaincode invoke -o orderer.example.com:7050 --channelID mychannel -n asset-transfer --isInit -c '{"Args":[]}'
```

### 3. Run the REST API

#### Step 1: Install Node.js Dependencies

Navigate to your project directory where `server.js` is located and run:

```bash
npm install
```

#### Step 2: Create Docker Image

Create a `Dockerfile` with the following content:

```dockerfile
# Step 1: Use the official Node.js image as the base image
FROM node:16

# Step 2: Set the working directory inside the container
WORKDIR /usr/src/app

# Step 3: Copy the package.json and package-lock.json to install dependencies
COPY package*.json ./

# Step 4: Install Node.js dependencies
RUN npm install

# Step 5: Copy the rest of the application code to the container
COPY . .

# Step 6: Expose the port the app will run on
EXPOSE 3000

# Step 7: Define the command to run the API server
CMD ["node", "server.js"]
```

Build the Docker image:

```bash
docker build -t rest-api-image .
```

Run the Docker container:

```bash
docker run -p 3000:3000 rest-api-image
```

### 4. Interact with the REST API

- **Create Asset**

  ```bash
  curl -X POST http://localhost:3000/asset -H "Content-Type: application/json" -d '{"dealerID":"D001", "msisdn":"1234567890", "mpin":"1234", "balance":1000, "status":"active", "transAmount":0, "transType":"", "remarks":"Initial balance"}'
  ```

- **Update Asset**

  ```bash
  curl -X PUT http://localhost:3000/asset -H "Content-Type: application/json" -d '{"dealerID":"D001", "msisdn":"1234567890", "mpin":"1234", "balance":2000, "status":"active", "transAmount":1000, "transType":"deposit", "remarks":"Updated balance"}'
  ```

- **Read Asset**

  ```bash
  curl http://localhost:3000/asset/D001
  ```

- **Get Asset History**

  ```bash
  curl http://localhost:3000/asset/D001/history
  ```

### 5. Running Test Cases

You can write test cases using a testing framework like Mocha and Chai. Here’s a basic example:

#### Step 1: Install Mocha and Chai

```bash
npm install --save-dev mocha chai
```

#### Step 2: Create a Test File

Create a file named `test.js` in your project directory:

```javascript
const chai = require('chai');
const chaiHttp = require('chai-http');
const server = require('./server');  // Assuming server.js exports the app

chai.should();
chai.use(chaiHttp);

describe('Asset API', () => {
  it('should create a new asset', (done) => {
    const asset = {
      dealerID: 'D002',
      msisdn: '0987654321',
      mpin: '4321',
      balance: 500,
      status: 'active',
      transAmount: 0,
      transType: '',
      remarks: 'New account'
    };
    chai.request(server)
      .post('/asset')
      .send(asset)
      .end((err, res) => {
        res.should.have.status(200);
        res.body.should.be.a('object');
        res.body.should.have.property('message').eql('Asset created successfully!');
        done();
      });
  });

  // Add more test cases as needed
});
```

#### Step 3: Run the Tests

```bash
npx mocha test.js
```

---