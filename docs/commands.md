# ElmoFi Token Website Management System

## Project Structure
```
elmofi/
├── config/           # Configuration files
├── docs/            # Documentation and requirements
├── monitoring/      # Contract monitoring services
├── scripts/         # Shell scripts for automation
└── src/            # Source files for websites
    ├── assets/     # Images and media files
    ├── js/         # JavaScript files
    └── styles/     # CSS files
```

## Available Commands

### Site Management

1. Create a new site:
```bash
./scripts/site-manager.sh create <site-name>
```
- Creates a new website with a unique design
- Automatically sets up monitoring
- Creates a timestamped directory

2. Delete a site:
```bash
./scripts/site-manager.sh delete <site-name>
```
- Removes the specified site
- Cleans up associated monitoring

3. List all sites:
```bash
./scripts/site-manager.sh list
```
- Shows all active sites with their creation times
- Displays monitoring status

### Contract Monitoring

1. Start contract monitoring:
```bash
./scripts/start-monitor.sh
```
- Initiates the Solana contract monitoring service
- Watches for new token contracts
- Logs activities to monitoring/logs

2. Update contract address:
```bash
./scripts/update-contract.sh <contract-address> <site-directory>
```
- Updates the contract address on a specific site
- Validates the contract address format
- Logs the update in monitoring/logs

### Configuration

1. Update RPC settings:
```bash
nano config/config.json
```
Edit the following parameters:
- `rpcUrl`: Your Solana RPC node URL
- `walletAddress`: Your wallet address
- `apiKey`: Your API key (if required)

### Monitoring Service

1. Start the Go monitoring service:
```bash
cd monitoring && go run contract-monitor.go
```
- Monitors Solana blockchain for new contracts
- Automatically updates website contract addresses
- Provides real-time logging

## Important Notes

1. All sites are automatically deleted after 10 minutes of creation
2. Contract addresses are automatically updated when new tokens are detected
3. Each site gets a unique design while maintaining core functionality
4. Logs are stored in monitoring/logs directory
5. Configuration changes require service restart

## Troubleshooting

If a command fails:
1. Check logs in monitoring/logs
2. Verify config/config.json settings
3. Ensure all required services are running
4. Check file permissions in scripts directory

## File Locations

- Website templates: src/
- Script files: scripts/
- Configuration: config/
- Monitoring service: monitoring/
- Documentation: docs/
- Assets and images: src/assets/

Remember to keep the config.json file updated with correct RPC and wallet information. 