package loghandler

import (
	"context"
	"go-contract-indexer/db"
	"go-contract-indexer/parser"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// handleLogs processes the logs received from the Ethereum client.
func HandleLogs(ctx context.Context, logs chan types.Log, sub ethereum.Subscription, db db.DBInterface, logger logrus.FieldLogger) {
	for {
		select {
		case err := <-sub.Err():
			logger.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			logger.Debugf("Received log: %v", vLog)
			event, err := parser.UnpackLog(vLog)
			if err != nil {
				logger.Printf("Failed to unpack log: %v", err)
				continue
			}

			switch e := event.(type) {
			case *parser.ERC20Transfer:
				handleTransferEvent(e, vLog, db, logger)
			case *parser.ERC20Approval:
				handleApprovalEvent(e, vLog, db, logger)
			default:
				logger.Printf("Unknown event type")
			}
		case <-ctx.Done():
			logger.Info("Shutting down log handling")
			sub.Unsubscribe()
			return
		}
	}
}

// handleTransferEvent handles the Transfer event logs.
func handleTransferEvent(e *parser.ERC20Transfer, vLog types.Log, db db.DBInterface, logger logrus.FieldLogger) {
	from := e.From.Hex()
	to := e.To.Hex()
	logger.Infof("Handling Transfer Event: From %s To %s Value %s", from, to, e.Value.String())
	err := db.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Transfer", &from, &to, nil, nil, e.Value)
	if err != nil {
		logger.Errorf("Failed to save transfer event: %v", err)
	}
}

// handleApprovalEvent handles the Approval event logs.
func handleApprovalEvent(e *parser.ERC20Approval, vLog types.Log, db db.DBInterface, logger logrus.FieldLogger) {
	owner := e.Owner.Hex()
	spender := e.Spender.Hex()
	logger.Infof("Handling Approval Event: Owner %s Spender %s Value %s", owner, spender, e.Value.String())
	err := db.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Approval", nil, nil, &owner, &spender, e.Value)
	if err != nil {
		logger.Errorf("Failed to save approval event: %v", err)
	}
}
