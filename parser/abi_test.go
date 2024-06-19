package parser

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

const mockedABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

// Test cases for UnpackLog function
func TestUnpackLog(t *testing.T) {
	// Set the mocked ABI
	SetABI(mockedABI)

	// Transfer event test
	transferLog := types.Log{
		Topics: []common.Hash{
			TransferEventSigHash,
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123"), // From address
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000456"), // To address
		},
		Data: hexToBytes("00000000000000000000000000000000000000000000000000000000000003e8"), // 1000 in hex
	}

	transferEvent, err := UnpackLog(transferLog)
	assert.Nil(t, err)
	assert.NotNil(t, transferEvent)

	transfer, ok := transferEvent.(*ERC20Transfer)
	assert.True(t, ok)
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000123"), transfer.From)
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000456"), transfer.To)
	assert.Equal(t, big.NewInt(1000), transfer.Value)

	// Approval event test
	approvalLog := types.Log{
		Topics: []common.Hash{
			ApprovalEventSigHash,
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123"), // Owner address
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000456"), // Spender address
		},
		Data: hexToBytes("00000000000000000000000000000000000000000000000000000000000001f4"), // 500 in hex
	}

	approvalEvent, err := UnpackLog(approvalLog)
	assert.Nil(t, err)
	assert.NotNil(t, approvalEvent)

	approval, ok := approvalEvent.(*ERC20Approval)
	assert.True(t, ok)
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000123"), approval.Owner)
	assert.Equal(t, common.HexToAddress("0x0000000000000000000000000000000000000456"), approval.Spender)
	assert.Equal(t, big.NewInt(500), approval.Value)
}

// Helper function to convert hex string to bytes
func hexToBytes(hexString string) []byte {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}
	return bytes
}
