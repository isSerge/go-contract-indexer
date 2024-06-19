package main

import (
	"math/big"
	"testing"
	"time"

	"go-contract-indexer/parser"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

// MockSubscription is a mock of ethereum.Subscription interface
type MockSubscription struct {
	mock.Mock
}

func (m *MockSubscription) Unsubscribe() {
	m.Called()
}

func (m *MockSubscription) Err() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

// MockDB is a mock of the db package's SaveEvent function
type MockDB struct {
	mock.Mock
}

func (m *MockDB) SaveEvent(blockNumber uint64, txHash, eventType string, from, to, owner, spender *string, value *big.Int) error {
	args := m.Called(blockNumber, txHash, eventType, from, to, owner, spender, value)
	return args.Error(0)
}

func TestHandleLogs(t *testing.T) {
	// Initialize the ABI
	parser.Init()

	// Create mock subscription and logs channel
	mockSub := new(MockSubscription)
	logs := make(chan types.Log, 1)

	// Mock the subscription error channel
	errChan := make(chan error)
	mockSub.On("Err").Return((<-chan error)(errChan))

	// Simulate a log being received with correct data length for Transfer event
	value := new(big.Int).SetInt64(1000) // Example value
	valueBytes := common.LeftPadBytes(value.Bytes(), 32)
	log := types.Log{
		Address: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Topics:  []common.Hash{parser.TransferEventSigHash, common.HexToHash("0x1234567890abcdef1234567890abcdef12345678"), common.HexToHash("0x1234567890abcdef1234567890abcdef12345679")},
		Data:    valueBytes,
	}
	logs <- log

	// Create a mock DB
	mockDB := new(MockDB)
	mockDB.On("SaveEvent", mock.AnythingOfType("uint64"), mock.AnythingOfType("string"), "Transfer", mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), (*string)(nil), (*string)(nil), mock.AnythingOfType("*big.Int")).Return(nil)

	// Start handling logs
	go handleLogs(logs, mockSub, mockDB)

	// Allow some time for the log to be processed
	time.Sleep(1 * time.Second)

	// Verify that the log was processed and saved
	mockDB.AssertExpectations(t)
}
