const fs = require('fs');
const path = require('path');

// Get contract address from command line argument
const contractAddress = process.argv[2];

if (!contractAddress) {
    console.error('Please provide contract address as argument');
    process.exit(1);
}

// Read the current localStorage data if it exists
const localStoragePath = path.join(__dirname, 'localStorage.json');
let storage = {};

try {
    if (fs.existsSync(localStoragePath)) {
        storage = JSON.parse(fs.readFileSync(localStoragePath, 'utf8'));
    }
} catch (error) {
    console.log('No existing localStorage found, creating new one');
}

// Update the contract address
storage.elmofiContractAddress = contractAddress;

// Save the updated localStorage
fs.writeFileSync(localStoragePath, JSON.stringify(storage, null, 2));

console.log(`Contract address updated to: ${contractAddress}`); 