const { Connection, PublicKey } = require('@solana/web3.js');
const fs = require('fs');

const SOLANA_RPC_URL = process.env.SOLANA_RPC_URL || 'https://api.mainnet-beta.solana.com';

async function updateContractInfo(contractAddress) {
    try {
        const connection = new Connection(SOLANA_RPC_URL, 'confirmed');
        const pubKey = new PublicKey(contractAddress);
        
        // Получаем информацию о контракте
        const accountInfo = await connection.getAccountInfo(pubKey);
        const balance = await connection.getBalance(pubKey);
        
        // Получаем последние транзакции
        const signatures = await connection.getSignaturesForAddress(pubKey, { limit: 10 });
        
        const contractData = {
            address: contractAddress,
            balance: balance / 1e9, // Конвертируем в SOL
            lastUpdate: new Date().toISOString(),
            recentTransactions: signatures.map(sig => ({
                signature: sig.signature,
                slot: sig.slot,
                timestamp: sig.blockTime ? new Date(sig.blockTime * 1000).toISOString() : null
            }))
        };

        // Сохраняем данные в файл
        fs.writeFileSync('config.json', JSON.stringify(contractData, null, 2));
        console.log('Contract info updated successfully');
        
        return contractData;
    } catch (error) {
        console.error('Error updating contract info:', error);
        throw error;
    }
}

// Экспортируем функцию для использования в других файлах
module.exports = { updateContractInfo }; 