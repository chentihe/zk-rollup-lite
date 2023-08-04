package txutils

import (
	"math/big"
)

const (
	minAccountIndex int64 = 0
	maxAccountIndex int64 = (1 << 32) - 1

	minNonce int64 = 0
)

var (
	minFeeAmount = big.NewInt(0)
	maxFeeAmount = new(big.Int).Mul(big.NewInt(2047), new(big.Int).Exp(big.NewInt(10), big.NewInt(31), nil))

	minAmount = big.NewInt(0)
	maxAmount = new(big.Int).Mul(big.NewInt(34359738367), new(big.Int).Exp(big.NewInt(10), big.NewInt(31), nil))

	Fee = ToWei(0.5, 18)
)
