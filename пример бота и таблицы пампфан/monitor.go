package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os/exec"
	"time"

	pb "matteo_solana/proto"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
)

const (
	// BuyDiscriminator = 0x66063d1201daeaea // старый вариант
	BuyDiscriminator  uint64 = 0xeaebda01123d0666 // правильный порядок байтов
	LAMPORTS_PER_SOL         = 1_000_000_000
	MinPurchaseAmount        = 1.9 * float64(LAMPORTS_PER_SOL)
)

var (
	pumpProgram = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	endpoint    = "fra.pixellabz.io:9000"
	kacp        = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
)

type PurchaseMonitor struct {
	client     pb.GeyserClient
	rpcClient  *rpc.Client
	wsClient   *ws.Client
	walletAddr solana.PublicKey
}

func NewPurchaseMonitor(walletAddr solana.PublicKey) (*PurchaseMonitor, error) {
	log.Println("Начинаем создание монитора...")

	// Создаем gRPC подключение
	log.Printf("Подключаемся к gRPC серверу %s...", endpoint)
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к gRPC серверу: %w", err)
	}
	log.Println("Успешно подключились к gRPC серверу")

	// Создаем клиенты
	log.Println("Создаем RPC клиент...")
	client := pb.NewGeyserClient(conn)
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	log.Println("Подключаемся к websocket...")
	wsClient, err := ws.Connect(context.Background(), "wss://api.mainnet-beta.solana.com")
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ошибка подключения к websocket: %w", err)
	}
	log.Println("Успешно подключились к websocket")

	return &PurchaseMonitor{
		client:     client,
		rpcClient:  rpcClient,
		wsClient:   wsClient,
		walletAddr: walletAddr,
	}, nil
}

func (pm *PurchaseMonitor) Start() error {
	log.Println("Запускаем мониторинг...")

	// Подписываемся на транзакции через Geyser
	var subscription pb.SubscribeRequest
	subscription.Transactions = make(map[string]*pb.SubscribeRequestFilterTransactions)
	f := false
	subscription.Transactions["pump_filter"] = &pb.SubscribeRequestFilterTransactions{
		AccountInclude: []string{pm.walletAddr.String()},
		Failed:         &f,
	}
	subscription.Commitment = pb.CommitmentLevel_PROCESSED.Enum()

	log.Println("Создаем gRPC stream...")
	stream, err := pm.client.Subscribe(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка подключения к gRPC серверу (возможно, ваш IP не в whitelist): %w", err)
	}
	log.Println("Успешно создали gRPC stream")

	log.Println("Отправляем подписку...")
	if err := stream.Send(&subscription); err != nil {
		return fmt.Errorf("failed to send subscription: %w", err)
	}
	log.Println("Успешно отправили подписку")
	log.Println("Ожидаем транзакции...")

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Ошибка получения подписки: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		if tx := resp.GetTransaction(); tx != nil {
			log.Println("Получена новая транзакция, проверяем...")
			message := tx.GetTransaction().GetTransaction().GetMessage()
			if message == nil {
				continue
			}

			// Проверяем, что транзакция содержит наш кошелек
			containsOurWallet := false
			for _, key := range message.AccountKeys {
				if solana.PublicKeyFromBytes(key).Equals(pm.walletAddr) {
					containsOurWallet = true
					break
				}
			}

			if !containsOurWallet {
				log.Printf("Пропускаем транзакцию - не содержит наш кошелек %s", pm.walletAddr.String())
				continue
			}

			log.Printf("Транзакция содержит наш кошелек %s, анализируем инструкции...", pm.walletAddr.String())

			for _, ix := range message.Instructions {
				program_idx := ix.ProgramIdIndex
				if int(program_idx) >= len(message.AccountKeys) {
					continue
				}
				program_id := message.AccountKeys[program_idx]
				if !solana.PublicKeyFromBytes(program_id).Equals(pumpProgram) {
					continue
				}

				log.Printf("Найдена инструкция для программы pump.fun")
				log.Printf("Данные инструкции (hex): %x", ix.Data)

				// Проверяем дискриминатор инструкции buy
				if len(ix.Data) < 8 {
					log.Println("Данные инструкции слишком короткие для дискриминатора")
					continue
				}
				discriminator := binary.LittleEndian.Uint64(ix.Data[:8])
				log.Printf("Найденный дискриминатор: %x", discriminator)
				log.Printf("Ожидаемый дискриминатор: %x", BuyDiscriminator)

				if discriminator != BuyDiscriminator {
					log.Printf("Неверный дискриминатор: %x (ожидался: %x)", discriminator, BuyDiscriminator)
					continue
				}

				log.Println("Найдена инструкция buy")

				// Проверяем, что транзакция от нашего кошелька
				user_idx := ix.Accounts[6]
				if int(user_idx) >= len(message.AccountKeys) {
					log.Println("Индекс пользователя вне диапазона")
					continue
				}
				user_pk := solana.PublicKeyFromBytes(message.AccountKeys[user_idx])
				if !user_pk.Equals(pm.walletAddr) {
					log.Printf("Транзакция от другого кошелька: %s (ожидался: %s)", user_pk.String(), pm.walletAddr.String())
					continue
				}

				log.Println("Транзакция от нашего кошелька")

				// Получаем сумму покупки
				if len(ix.Data) < 16 {
					log.Println("Данные инструкции слишком короткие для суммы")
					continue
				}
				purchaseAmount := binary.LittleEndian.Uint64(ix.Data[8:16])
				amountInSol := float64(purchaseAmount) / float64(LAMPORTS_PER_SOL)
				log.Printf("Сумма покупки: %.9f SOL (в lamports: %d)", amountInSol, purchaseAmount)
				log.Printf("Минимальная требуемая сумма: %.9f SOL (в lamports: %d)", MinPurchaseAmount/float64(LAMPORTS_PER_SOL), MinPurchaseAmount)

				// Проверяем, что сумма покупки >= 1.9 SOL
				if float64(purchaseAmount) >= MinPurchaseAmount {
					log.Printf("🎯 ОБНАРУЖЕНА ПОКУПКА! Сумма: %.9f SOL (>= %.9f SOL)\n", amountInSol, MinPurchaseAmount/float64(LAMPORTS_PER_SOL))
					log.Printf("🔔 Воспроизводим звуковой сигнал...")

					// Воспроизводим звуковой сигнал
					cmd := exec.Command("cmd", "/C", "start", "sistema-poiska-pidorasov.wav")
					if err := cmd.Run(); err != nil {
						log.Printf("❌ Ошибка воспроизведения звука: %v\n", err)
					} else {
						log.Printf("✅ Звуковой сигнал успешно воспроизведен")
					}
				} else {
					log.Printf("ℹ️ Обнаружена покупка на сумму %.9f SOL (меньше минимальной %.9f SOL)\n",
						amountInSol, MinPurchaseAmount/float64(LAMPORTS_PER_SOL))
				}
			}
		}
	}
}

func main() {
	// Создаем монитор для указанного кошелька
	walletAddr := solana.MustPublicKeyFromBase58("FgWxH43h72i43vQwaSo8Zd43nG9Eh5ErrG2ZkNShzk44")
	monitor, err := NewPurchaseMonitor(walletAddr)
	if err != nil {
		log.Fatalf("Ошибка создания монитора: %v", err)
	}

	// Запускаем мониторинг
	if err := monitor.Start(); err != nil {
		log.Fatalf("Ошибка мониторинга: %v", err)
	}
}
