class ContractUpdater {
    constructor() {
        this.addressElement = document.getElementById('contractAddressText');
        this.copyButton = document.getElementById('copyAddress');
        this.contractAddress = null;
        this.checkInterval = null;
    }

    // Function to update the contract address on the page
    updateContractAddress(address) {
        this.contractAddress = address;
        this.addressElement.textContent = address;
        this.copyButton.style.display = 'block';
        
        // Save to localStorage
        localStorage.setItem('elmofiContractAddress', address);

        // Add animation class
        this.addressElement.classList.add('address-updated');
        
        // Remove animation class after animation completes
        setTimeout(() => {
            this.addressElement.classList.remove('address-updated');
        }, 1000);
    }

    // Function to copy address to clipboard
    setupCopyButton() {
        this.copyButton.addEventListener('click', () => {
            if (this.contractAddress) {
                navigator.clipboard.writeText(this.contractAddress)
                    .then(() => {
                        const originalText = this.copyButton.textContent;
                        this.copyButton.textContent = 'Copied!';
                        setTimeout(() => {
                            this.copyButton.textContent = originalText;
                        }, 2000);
                    })
                    .catch(err => console.error('Failed to copy:', err));
            }
        });
    }

    // Function to check if contract is already stored
    checkStoredContract() {
        const storedAddress = localStorage.getItem('elmofiContractAddress');
        if (storedAddress) {
            this.updateContractAddress(storedAddress);
        }
    }

    // Function to start checking for contract address
    async checkForContract() {
        try {
            const response = await fetch('YOUR_API_ENDPOINT'); // Replace with your actual API endpoint
            const data = await response.json();
            
            if (data.contractAddress) {
                this.updateContractAddress(data.contractAddress);
                // Stop checking once we have the address
                if (this.checkInterval) {
                    clearInterval(this.checkInterval);
                    this.checkInterval = null;
                }
            }
        } catch (error) {
            console.error('Error checking contract:', error);
        }
    }

    // Initialize the contract updater
    init() {
        this.setupCopyButton();
        this.checkStoredContract();
        
        // Start periodic checking if we don't have the address
        if (!this.contractAddress) {
            this.checkInterval = setInterval(() => this.checkForContract(), 30000); // Check every 30 seconds
            // Also check immediately
            this.checkForContract();
        }
    }

    // Manual update function (can be called from outside)
    manualUpdate(address) {
        if (address && this.isValidSolanaAddress(address)) {
            this.updateContractAddress(address);
            if (this.checkInterval) {
                clearInterval(this.checkInterval);
                this.checkInterval = null;
            }
            return true;
        }
        return false;
    }

    // Validate Solana address format
    isValidSolanaAddress(address) {
        return /^[1-9A-HJ-NP-Za-km-z]{32,44}$/.test(address);
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.contractUpdater = new ContractUpdater();
    window.contractUpdater.init();
});

// Listen for contract address updates
window.addEventListener('message', function(event) {
    // Verify message origin if needed
    // if (event.origin !== "expected_origin") return;
    
    const data = event.data;
    if (data && data.type === 'CONTRACT_ADDRESS_UPDATE') {
        const newAddress = data.address;
        
        if (isValidSolanaAddress(newAddress)) {
            updateContractAddress(newAddress);
            console.log('Contract address updated successfully');
        } else {
            console.error('Invalid Solana address format');
        }
    }
}); 