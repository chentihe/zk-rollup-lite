package txmanager

import (
	"math/big"

	"github.com/iden3/go-rapidsnark/types"
)

type DepositInfo struct {
	AccountIndex  int64          `json:"accountIndex"`
	PublicKey     string         `json:"publicKey"`
	DepositAmount *big.Int       `json:"depositAmount"`
	SignedTxHash  string         `json:"signTxHash"`
	ZkProof       *types.ZKProof `json:"zkProof"`
}
