package main

import (
	"context"
	"errors"
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
	"go-contract-indexer/loghandler"
	"go-contract-indexer/parser"
)

var logger = logrus.New()

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("Error reading config file, %s", err)
	}

	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.Level = logrus.DebugLevel
}

func validateConfig() error {
	if viper.GetString("RPC_URL") == "" {
		return errors.New("RPC_URL is required")
	}
	if viper.GetString("CONTRACT_ADDRESS") == "" {
		return errors.New("CONTRACT_ADDRESS is required")
	}
	if viper.GetString("DB_CONN_STR") == "" {
		return errors.New("DB_CONN_STR is required")
	}
	return nil
}

func main() {
	if err := validateConfig(); err != nil {
		logger.Fatalf("Configuration validation failed: %v", err)
	}
	// Load configuration
	rpcURL := viper.GetString("RPC_URL")
	contractAddress := viper.GetString("CONTRACT_ADDRESS")
	dbConnStr := viper.GetString("DB_CONN_STR")

	if rpcURL == "" || contractAddress == "" || dbConnStr == "" {
		logger.Fatal("RPC_URL, CONTRACT_ADDRESS, or DB_CONN_STR is not set in the configuration")
	}

	// Initialize the database connection with retry mechanism
	database := db.InitDB(dbConnStr)
	defer database.Close()

	// Initialize the ABI
	parser.Init()

	// Connect to the Ethereum client
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		logger.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddr := common.HexToAddress(contractAddress)

	// Subscribe to the logs of the contract
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logger.Fatalf("Failed to subscribe to logs: %v", err)
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel, sub)

	// Print contract details
	err = printTokenInfo(client, contractAddr)
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}

	// Create LogHandler instance
	logHandler := loghandler.NewLogHandler(database, logger)

	// Handle incoming logs
	logHandler.HandleLogs(ctx, logs, sub)
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

	logger.WithFields(logrus.Fields{
		"address": contractAddr.Hex(),
		"name":    name,
		"symbol":  symbol,
	}).Info("Starting indexer for ERC-20 contract")

	return nil
}

// handleShutdown handles graceful shutdown on receiving a termination signal.
func handleShutdown(cancel context.CancelFunc, sub ethereum.Subscription) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	logger.Info("Received shutdown signal")
	cancel()
	sub.Unsubscribe()
}
