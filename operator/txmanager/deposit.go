package txmanager

import "math/big"

type DepositInfo struct {
	AccountIndex  int64
	PublicKey     string
	DepositAmount *big.Int
	SignedTxHash  string
}
