package parser

import (
	"bytes"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Global variable for parsed ABI
var parsedABI abi.ABI

// Load the ERC-20 ABI from file
func init() {
	abiFile := "erc20/erc20.abi"
	abiData, err := os.ReadFile(abiFile)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	parsedABI, err = abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}
}

// Helper function to unpack logs
func UnpackLog(log types.Log) (interface{}, error) {
	switch log.Topics[0] {
	case transferEventSigHash:
		event := new(ERC20Transfer)
		err := parsedABI.UnpackIntoInterface(event, "Transfer", log.Data)
		if err != nil {
			return nil, err
		}
		event.From = common.HexToAddress(log.Topics[1].Hex())
		event.To = common.HexToAddress(log.Topics[2].Hex())
		return event, nil
	case approvalEventSigHash:
		event := new(ERC20Approval)
		err := parsedABI.UnpackIntoInterface(event, "Approval", log.Data)
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
