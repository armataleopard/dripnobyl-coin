#!/bin/bash

# Проверяем наличие переменных окружения
if [ -z "$CONTRACT_ADDRESS" ]; then
    echo "Error: CONTRACT_ADDRESS not set"
    exit 1
fi

if [ -z "$WEBHOOK_SECRET" ]; then
    echo "Error: WEBHOOK_SECRET not set"
    exit 1
fi

# Запускаем веб-хук сервер
echo "Starting webhook server..."
node contract-webhook.js &

# Ждем 5 секунд, чтобы сервер успел запуститься
sleep 5

# Запускаем монитор контракта
echo "Starting contract monitor..."
python3 watch_contract.py

# Функция для очистки процессов при выходе
cleanup() {
    echo "Stopping services..."
    pkill -f "node contract-webhook.js"
    pkill -f "python3 watch_contract.py"
}

# Регистрируем функцию очистки
trap cleanup EXIT

# Ждем сигналы завершения
wait 