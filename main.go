package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"go-contract-indexer/db"
	"go-contract-indexer/erc20"
	"go-contract-indexer/parser"
)

var log = logrus.New()

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	log.Level = logrus.DebugLevel
}

func main() {
	// Load configuration
	rpcURL := viper.GetString("RPC_URL")
	contractAddress := viper.GetString("CONTRACT_ADDRESS")
	dbConnStr := viper.GetString("DB_CONN_STR")

	if rpcURL == "" || contractAddress == "" || dbConnStr == "" {
		log.Fatal("RPC_URL, CONTRACT_ADDRESS, or DB_CONN_STR is not set in the configuration")
	}

	// Initialize the database connection
	database := db.InitDB(dbConnStr)

	// Initialize the ABI
	parser.Init()

	// Connect to the Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddr := common.HexToAddress(contractAddress)

	// Subscribe to the logs of the contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel, sub)

	// Print contract details
	err = printTokenInfo(client, contractAddr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Handle incoming logs
	handleLogs(ctx, logs, sub, database)
}

// printTokenInfo prints the token information for the given contract address.
func printTokenInfo(client *ethclient.Client, contractAddr common.Address) error {
	token, err := erc20.NewErc20(contractAddr, client)
	if err != nil {
		return fmt.Errorf("failed to instantiate token contract: %v", err)
	}

	callOpts := &bind.CallOpts{}

	name, err := token.Name(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get token name: %v", err)
	}

	symbol, err := token.Symbol(callOpts)
	if err != nil {
		return fmt.Errorf("failed to get token symbol: %v", err)
	}

	log.Infof("Starting indexer for ERC-20 contract:")
	log.Infof("Address: %s", contractAddr.Hex())
	log.Infof("Name: %s", name)
	log.Infof("Symbol: %s", symbol)

	return nil
}

// handleLogs processes the logs received from the Ethereum client.
func handleLogs(ctx context.Context, logs chan types.Log, sub ethereum.Subscription, db db.DBInterface) {
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			event, err := parser.UnpackLog(vLog)
			if err != nil {
				log.Printf("Failed to unpack log: %v", err)
				continue
			}

			switch e := event.(type) {
			case *parser.ERC20Transfer:
				handleTransferEvent(e, vLog, db)
			case *parser.ERC20Approval:
				handleApprovalEvent(e, vLog, db)
			default:
				log.Printf("Unknown event type")
			}
		case <-ctx.Done():
			log.Info("Shutting down log handling")
			sub.Unsubscribe()
			return
		}
	}
}

// handleTransferEvent handles the Transfer event logs.
func handleTransferEvent(e *parser.ERC20Transfer, vLog types.Log, db db.DBInterface) {
	from := e.From.Hex()
	to := e.To.Hex()
	log.Infof("Transfer Event: From %s To %s Value %s", from, to, e.Value.String())
	err := db.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Transfer", &from, &to, nil, nil, e.Value)
	if err != nil {
		log.Errorf("Failed to save transfer event: %v", err)
	}
}

// handleApprovalEvent handles the Approval event logs.
func handleApprovalEvent(e *parser.ERC20Approval, vLog types.Log, db db.DBInterface) {
	owner := e.Owner.Hex()
	spender := e.Spender.Hex()
	log.Infof("Approval Event: Owner %s Spender %s Value %s", owner, spender, e.Value.String())
	err := db.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Approval", nil, nil, &owner, &spender, e.Value)
	if err != nil {
		log.Errorf("Failed to save approval event: %v", err)
	}
}

// handleShutdown handles graceful shutdown on receiving a termination signal.
func handleShutdown(cancel context.CancelFunc, sub ethereum.Subscription) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Info("Received shutdown signal")
	cancel()
	sub.Unsubscribe()
}
