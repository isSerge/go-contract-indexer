package loghandler

import (
	"context"
	"go-contract-indexer/db"
	"go-contract-indexer/parser"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type LogHandler struct {
	DB     db.DBInterface
	Logger logrus.FieldLogger
}

func NewLogHandler(db db.DBInterface, logger logrus.FieldLogger) *LogHandler {
	return &LogHandler{
		DB:     db,
		Logger: logger,
	}
}

// HandleLogs processes the logs received from the Ethereum client.
func (h *LogHandler) HandleLogs(ctx context.Context, logs chan types.Log, sub ethereum.Subscription) {
	for {
		select {
		case err := <-sub.Err():
			h.Logger.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			h.Logger.Debugf("Received log: %v", vLog)
			event, err := parser.UnpackLog(vLog)
			if err != nil {
				h.Logger.Printf("Failed to unpack log: %v", err)
				continue
			}

			switch e := event.(type) {
			case *parser.ERC20Transfer:
				h.handleTransferEvent(e, vLog)
			case *parser.ERC20Approval:
				h.handleApprovalEvent(e, vLog)
			default:
				h.Logger.Printf("Unknown event type")
			}
		case <-ctx.Done():
			h.Logger.Info("Shutting down log handling")
			sub.Unsubscribe()
			return
		}
	}
}

// handleTransferEvent handles the Transfer event logs.
func (h *LogHandler) handleTransferEvent(e *parser.ERC20Transfer, vLog types.Log) {
	from := e.From.Hex()
	to := e.To.Hex()
	h.Logger.Infof("Handling Transfer Event: From %s To %s Value %s", from, to, e.Value.String())
	err := h.DB.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Transfer", &from, &to, nil, nil, e.Value)
	if err != nil {
		h.Logger.Errorf("Failed to save transfer event: %v", err)
	}
}

// handleApprovalEvent handles the Approval event logs.
func (h *LogHandler) handleApprovalEvent(e *parser.ERC20Approval, vLog types.Log) {
	owner := e.Owner.Hex()
	spender := e.Spender.Hex()
	h.Logger.Infof("Handling Approval Event: Owner %s Spender %s Value %s", owner, spender, e.Value.String())
	err := h.DB.SaveEvent(vLog.BlockNumber, vLog.TxHash.Hex(), "Approval", nil, nil, &owner, &spender, e.Value)
	if err != nil {
		h.Logger.Errorf("Failed to save approval event: %v", err)
	}
}
