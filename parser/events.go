package parser

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
)

var (
	transferEventSignature = []byte("Transfer(address,address,uint256)")
	approvalEventSignature = []byte("Approval(address,address,uint256)")

	// TransferEventSigHash is the signature hash for the Transfer event.
	TransferEventSigHash = crypto.Keccak256Hash(transferEventSignature)
	// ApprovalEventSigHash is the signature hash for the Approval event.
	ApprovalEventSigHash = crypto.Keccak256Hash(approvalEventSignature)
)

// ERC20Transfer represents the Transfer event.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ERC20Approval represents the Approval event.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}
