package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	"go-contract-indexer/db"
	"go-contract-indexer/erc20"
	"go-contract-indexer/parser"
)

func main() {
	// Load configuration
	rpcURL, contractAddress, dbConnStr := loadConfig()

	// Initialize the database connection
	db.InitDB(dbConnStr)

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

	// Print contract details
	err = printTokenInfo(client, contractAddr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Handle incoming logs
	handleLogs(logs, sub)
}

// loadConfig loads the configuration from the .env file.
func loadConfig() (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	rpcURL := os.Getenv("RPC_URL")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")
	dbConnStr := os.Getenv("DB_CONN_STR")

	if rpcURL == "" || contractAddress == "" || dbConnStr == "" {
		log.Fatal("RPC_URL, CONTRACT_ADDRESS, or DB_CONN_STR is not set in the .env file")
	}

	return rpcURL, contractAddress, dbConnStr
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

	fmt.Printf("Starting indexer for ERC-20 contract: \n")
	fmt.Printf("Address: %s\n", contractAddr.Hex())
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Symbol: %s\n", symbol)

	return nil
}

// handleLogs processes the logs received from the Ethereum client.
func handleLogs(logs chan types.Log, sub ethereum.Subscription) {
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
				handleTransferEvent(e, vLog)
			case *parser.ERC20Approval:
				handleApprovalEvent(e, vLog)
			default:
				log.Printf("Unknown event type")
			}
		}
	}
}

// handleTransferEvent handles the Transfer event logs.
func handleTransferEvent(e *parser.ERC20Transfer, vLog types.Log) {
	from := e.From.Hex()
	to := e.To.Hex()
	fmt.Printf("Transfer Event: From %s To %s Value %s\n", from, to, e.Value.String())
	err := db.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Transfer", &from, &to, nil, nil, e.Value)
	if err != nil {
		log.Printf("Failed to save transfer event: %v", err)
	}
}

// handleApprovalEvent handles the Approval event logs.
func handleApprovalEvent(e *parser.ERC20Approval, vLog types.Log) {
	owner := e.Owner.Hex()
	spender := e.Spender.Hex()
	fmt.Printf("Approval Event: Owner %s Spender %s Value %s\n", owner, spender, e.Value.String())
	err := db.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Approval", nil, nil, &owner, &spender, e.Value)
	if err != nil {
		log.Printf("Failed to save approval event: %v", err)
	}
}
