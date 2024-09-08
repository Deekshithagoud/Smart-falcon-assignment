const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const app = express();
app.use(express.json());

const ccpPath = path.resolve(__dirname, '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

const walletPath = path.join(process.cwd(), 'wallet');
let wallet;

async function initialize() {
    wallet = await Wallets.newFileSystemWallet(walletPath);
}

async function getContract() {
    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'appUser',  // Ensure 'appUser' is registered and enrolled
        discovery: { enabled: true, asLocalhost: true }
    });

    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('fabcar');  // Change 'fabcar' to your chaincode name
    return { gateway, contract };
}

// Create Asset Endpoint
app.post('/asset', async (req, res) => {
    try {
        const { dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks } = req.body;

        const { gateway, contract } = await getContract();
        await contract.submitTransaction('CreateAsset', dealerID, msisdn, mpin, status, transType, remarks, balance.toString(), transAmount.toString());

        res.json({ message: 'Asset created successfully!' });
        await gateway.disconnect();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Update Asset Endpoint
app.put('/asset', async (req, res) => {
    try {
        const { dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks } = req.body;

        const { gateway, contract } = await getContract();
        await contract.submitTransaction('UpdateAsset', dealerID, msisdn, mpin, status, transType, remarks, balance.toString(), transAmount.toString());

        res.json({ message: 'Asset updated successfully!' });
        await gateway.disconnect();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Read Asset Endpoint
app.get('/asset/:id', async (req, res) => {
    try {
        const { id } = req.params;

        const { gateway, contract } = await getContract();
        const result = await contract.evaluateTransaction('ReadAsset', id);

        res.json(JSON.parse(result.toString()));
        await gateway.disconnect();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

// Get Asset History Endpoint
app.get('/asset/:id/history', async (req, res) => {
    try {
        const { id } = req.params;

        const { gateway, contract } = await getContract();
        const result = await contract.evaluateTransaction('GetAssetHistory', id);

        res.json(JSON.parse(result.toString()));
        await gateway.disconnect();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

const PORT = process.env.PORT || 3000;

app.listen(PORT, async () => {
    await initialize();
    console.log(`REST API server running on port ${PORT}`);
});
