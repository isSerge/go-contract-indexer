package main

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"go-contract-indexer/parser"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

func initTestConfig() {
	viper.SetConfigName("config_test") // Use a test-specific config file
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.ReadInConfig()
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
	mockSub.On("Unsubscribe").Run(func(args mock.Arguments) {
		t.Log("Unsubscribe called")
	}).Return()

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

	// Create a context and a cancel function to simulate graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start handling logs
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleLogs(ctx, logs, mockSub, mockDB)
	}()

	// Allow some time for the log to be processed
	time.Sleep(1 * time.Second)

	// Verify that the log was processed and saved
	mockDB.AssertExpectations(t)

	// Simulate shutdown by canceling the context
	cancel()

	// Allow some time for the shutdown to complete
	time.Sleep(1 * time.Second)

	// Wait for handleLogs to return
	wg.Wait()

	// Ensure that the subscription was unsubscribed
	mockSub.AssertCalled(t, "Unsubscribe")
}

func TestLoadConfig(t *testing.T) {
	initTestConfig()

	rpcURL := viper.GetString("RPC_URL")
	contractAddress := viper.GetString("CONTRACT_ADDRESS")
	dbConnStr := viper.GetString("DB_CONN_STR")

	assert.NotEmpty(t, rpcURL, "RPC_URL should not be empty")
	assert.NotEmpty(t, contractAddress, "CONTRACT_ADDRESS should not be empty")
	assert.NotEmpty(t, dbConnStr, "DB_CONN_STR should not be empty")
}

func TestLoadConfig_Error(t *testing.T) {
	// Clear any previously set environment variables
	viper.Reset()

	// Initialize without setting environment variables to test error case
	initTestConfig()

	rpcURL := viper.GetString("RPC_URL")
	contractAddress := viper.GetString("CONTRACT_ADDRESS")
	dbConnStr := viper.GetString("DB_CONN_STR")

	assert.Empty(t, rpcURL, "RPC_URL should be empty")
	assert.Empty(t, contractAddress, "CONTRACT_ADDRESS should be empty")
	assert.Empty(t, dbConnStr, "DB_CONN_STR should be empty")
}
