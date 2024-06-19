package loghandler

import (
	"context"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"go-contract-indexer/parser"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Initialize the environment variable for the ABI path
func init() {
	os.Setenv("ERC20_ABI_PATH", "../erc20/erc20.abi") // Adjust the path as necessary
}

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

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestHandleTransferLogs(t *testing.T) {
	// Initialize the ABI
	parser.Init()

	// Create mock logger
	mockLogger := logrus.New()

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
	mockDB.On("Close").Return(nil).Once() // Ensure it is expected once

	// Create a context and a cancel function to simulate graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Create LogHandler instance
	logHandler := NewLogHandler(mockDB, mockLogger)

	// Start handling logs
	wg.Add(1)
	go func() {
		defer wg.Done()
		logHandler.HandleLogs(ctx, logs, mockSub)
	}()

	// Allow some time for the log to be processed
	time.Sleep(1 * time.Second)

	// Simulate shutdown by canceling the context
	cancel()

	// Allow some time for the shutdown to complete
	time.Sleep(1 * time.Second)

	// Wait for handleLogs to return
	wg.Wait()

	// Ensure that the subscription was unsubscribed
	mockSub.AssertCalled(t, "Unsubscribe")

	// Close the mock database connection
	err := mockDB.Close()
	assert.NoError(t, err, "Expected no error while closing the database connection")
	mockDB.AssertExpectations(t)
}

func TestHandleApprovalLogs(t *testing.T) {
	// Initialize the ABI
	parser.Init()

	// Create mock logger
	mockLogger := logrus.New()

	// Create mock subscription and logs channel
	mockSub := new(MockSubscription)
	logs := make(chan types.Log, 1)

	// Mock the subscription error channel
	errChan := make(chan error)
	mockSub.On("Err").Return((<-chan error)(errChan))
	mockSub.On("Unsubscribe").Run(func(args mock.Arguments) {
		t.Log("Unsubscribe called")
	}).Return()

	// Simulate a log being received with correct data length for Approval event
	value := new(big.Int).SetInt64(1000) // Example value
	valueBytes := common.LeftPadBytes(value.Bytes(), 32)
	log := types.Log{
		Address: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Topics:  []common.Hash{parser.ApprovalEventSigHash, common.HexToHash("0x1234567890abcdef1234567890abcdef12345678"), common.HexToHash("0x1234567890abcdef1234567890abcdef12345679")},
		Data:    valueBytes,
	}
	logs <- log

	// Create a mock DB
	mockDB := new(MockDB)
	mockDB.On("SaveEvent", mock.AnythingOfType("uint64"), mock.AnythingOfType("string"), "Approval", (*string)(nil), (*string)(nil), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*big.Int")).Return(nil)
	mockDB.On("Close").Return(nil).Once() // Ensure it is expected once

	// Create a context and a cancel function to simulate graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Create LogHandler instance
	logHandler := NewLogHandler(mockDB, mockLogger)

	// Start handling logs
	wg.Add(1)
	go func() {
		defer wg.Done()
		logHandler.HandleLogs(ctx, logs, mockSub)
	}()

	// Allow some time for the log to be processed
	time.Sleep(1 * time.Second)

	// Simulate shutdown by canceling the context
	cancel()
	t.Log("Context cancelled")

	// Allow some time for the shutdown to complete
	time.Sleep(1 * time.Second)

	// Wait for handleLogs to return
	wg.Wait()
	t.Log("Wait group done")

	// Ensure that the subscription was unsubscribed
	mockSub.AssertCalled(t, "Unsubscribe")

	// Close the mock database connection
	err := mockDB.Close()
	assert.NoError(t, err, "Expected no error while closing the database connection")
	mockDB.AssertExpectations(t)
}
