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

	"go-contract-indexer/erc20"
	"go-contract-indexer/parser"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the RPC_URL and CONTRACT_ADDRESS from the .env file
	rpcURL := os.Getenv("RPC_URL")
	contractAddress := os.Getenv("CONTRACT_ADDRESS")

	if rpcURL == "" || contractAddress == "" {
		log.Fatal("RPC_URL or CONTRACT_ADDRESS is not set in the .env file")
	}

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

	// Create a channel to receive the logs from the contract and subscribe to the logs
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	// Print contract address, token name, and symbol
	err = printTokenInfo(client, contractAddr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error: %v", err)
		case vLog := <-logs:
			event, err := parser.UnpackLog(vLog)
			if err != nil {
				log.Fatalf("Failed to unpack log: %v", err)
			}

			switch e := event.(type) {
			case *parser.ERC20Transfer:
				fmt.Printf("Transfer Event: From %s To %s Value %s\n", e.From.Hex(), e.To.Hex(), e.Value.String())
			case *parser.ERC20Approval:
				fmt.Printf("Approval Event: Owner %s Spender %s Value %s\n", e.Owner.Hex(), e.Spender.Hex(), e.Value.String())
			default:
				fmt.Printf("Unknown event type\n")
			}
		}
	}
}

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
