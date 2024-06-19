package main

import (
	"math/big"
	"testing"

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
