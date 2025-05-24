#!/bin/bash

if [ -z "$1" ]; then
    echo "Please provide contract address as argument"
    echo "Usage: ./update-contract.sh <contract_address>"
    exit 1
fi

# Update the contract address in index.html
sed -i '' "s|<p id=\"contractAddressText\">.*</p>|<p id=\"contractAddressText\">$1</p>|g" index.html

echo "Contract address updated to: $1" 