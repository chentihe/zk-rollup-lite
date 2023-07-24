package txmanager

import "math/big"

type DepositInfo struct {
	AccountIndex  int64    `json:"accountIndex"`
	PublicKey     string   `json:"publicKey"`
	DepositAmount *big.Int `json:"depositAmount"`
	SignedTxHash  string   `json:"signTxHash"`
}
