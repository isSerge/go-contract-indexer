package parser

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
)

var (
	transferEventSignature = []byte("Transfer(address,address,uint256)")
	approvalEventSignature = []byte("Approval(address,address,uint256)")

	TransferEventSigHash = crypto.Keccak256Hash(transferEventSignature)
	ApprovalEventSigHash = crypto.Keccak256Hash(approvalEventSignature)
)

type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}
