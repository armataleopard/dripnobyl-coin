const express = require('express');
const { updateContractInfo } = require('./contract-updater');

const app = express();
app.use(express.json());

const CONTRACT_ADDRESS = process.env.CONTRACT_ADDRESS || '';
const WEBHOOK_SECRET = process.env.WEBHOOK_SECRET || '';

app.post('/update-contract', async (req, res) => {
    try {
        // Проверка секретного ключа
        const secret = req.headers['x-webhook-secret'];
        if (secret !== WEBHOOK_SECRET) {
            return res.status(401).json({ error: 'Unauthorized' });
        }

        // Обновление информации о контракте
        const contractData = await updateContractInfo(CONTRACT_ADDRESS);
        
        // Отправляем обновленные данные
        res.json({
            success: true,
            data: contractData
        });
    } catch (error) {
        console.error('Webhook error:', error);
        res.status(500).json({
            success: false,
            error: error.message
        });
    }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Webhook server running on port ${PORT}`);
}); 