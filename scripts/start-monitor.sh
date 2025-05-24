#!/bin/bash

# Проверяем наличие необходимых файлов
if [ ! -f "pump.json" ]; then
    echo "Error: pump.json not found"
    exit 1
fi

if [ ! -f "update-contract.sh" ]; then
    echo "Error: update-contract.sh not found"
    exit 1
fi

# Запускаем мониторинг
go run contract-monitor.go 