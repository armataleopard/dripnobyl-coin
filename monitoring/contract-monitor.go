package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Config структура для конфигурации
type Config struct {
	WalletAddress string  `json:"wallet_address"`
	RpcEndpoint   string  `json:"rpc_endpoint"`
	CheckInterval int     `json:"check_interval"`
	AutoStop      bool    `json:"auto_stop"`
	DebugMode     bool    `json:"debug_mode"`
}

var (
	processedTxs = make(map[string]bool)
	config       Config
)

func main() {
	// Настраиваем логирование
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	
	// Загружаем конфигурацию
	loadConfig()

	if config.DebugMode {
		log.Printf("Starting monitor for wallet: %s", config.WalletAddress)
	}

	// Канал для новых контрактов
	contractChan := make(chan string)

	// Запускаем мониторинг в отдельной горутине
	go monitorTransactions(contractChan)

	// Обрабатываем новые контракты
	for contract := range contractChan {
		if contract != "" {
			updateWebsite(contract)
			if config.AutoStop {
				log.Printf("Contract updated, monitoring stopped")
				return
			}
		}
	}
}

func loadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Проверяем обязательные поля
	if config.WalletAddress == "" {
		log.Fatal("Wallet address is required in config")
	}
	if config.RpcEndpoint == "" {
		log.Fatal("RPC endpoint is required in config")
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = 1 // Значение по умолчанию - 1 секунда
	}
}

func monitorTransactions(contractChan chan string) {
	for {
		// Здесь будет ваш код из monitor.go для:
		// 1. Получения последних транзакций вашего кошелька
		// 2. Проверки транзакций на создание токена
		// 3. Если найден новый токен, отправляем его адрес в канал:
		//    contractChan <- tokenAddress

		if config.DebugMode {
			log.Printf("Checking for new transactions...")
		}

		time.Sleep(time.Duration(config.CheckInterval) * time.Second)
	}
}

func updateWebsite(contractAddress string) {
	if config.DebugMode {
		log.Printf("Updating website with contract: %s", contractAddress)
	}

	scriptPath := "./update-contract.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.Printf("Warning: update-contract.sh not found")
		return
	}

	cmd := exec.Command(scriptPath, contractAddress)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Error updating website: %v", err)
		return
	}

	log.Printf("Successfully updated website with contract: %s", contractAddress)
}

func isNewTransaction(txId string) bool {
	if _, exists := processedTxs[txId]; exists {
		return false
	}
	processedTxs[txId] = true
	return true
} 