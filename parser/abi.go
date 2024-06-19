package parser

import (
	"bytes"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Global variable for parsed ABI and a mutex for safe concurrent access
var (
	ParsedABI abi.ABI
	once      sync.Once
)

// Load the ERC-20 ABI from file
func loadABI() {
	abiPath := os.Getenv("ERC20_ABI_PATH")
	if abiPath == "" {
		abiPath = "erc20/erc20.abi" // Default path
	}

	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	ParsedABI, err = abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}
}

// Init initializes the ABI, safe for concurrent use
func Init() {
	once.Do(loadABI)
}

// SetABI allows setting a mocked ABI for testing
func SetABI(mockedABI string) {
	var err error
	ParsedABI, err = abi.JSON(bytes.NewReader([]byte(mockedABI)))
	if err != nil {
		log.Fatalf("Failed to parse mocked ABI: %v", err)
	}
}

// UnpackLog unpacks logs into their respective events.
func UnpackLog(log types.Log) (interface{}, error) {
	switch log.Topics[0] {
	case TransferEventSigHash:
		event := new(ERC20Transfer)
		err := ParsedABI.UnpackIntoInterface(event, "Transfer", log.Data)
		if err != nil {
			return nil, err
		}
		event.From = common.HexToAddress(log.Topics[1].Hex())
		event.To = common.HexToAddress(log.Topics[2].Hex())
		return event, nil
	case ApprovalEventSigHash:
		event := new(ERC20Approval)
		err := ParsedABI.UnpackIntoInterface(event, "Approval", log.Data)
		if err != nil {
			return nil, err
		}
		event.Owner = common.HexToAddress(log.Topics[1].Hex())
		event.Spender = common.HexToAddress(log.Topics[2].Hex())
		return event, nil
	default:
		return nil, nil
	}
}
