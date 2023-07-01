package txhandlers

import (
	"math/big"
)

type TransactionInfo struct {
	From      int64
	To        int64
	Amount    *big.Int
	Fee       *big.Int
	Nonce     int64
	Signature []byte
}

func (txInfo *TransactionInfo) Validate(fromAccountNonce int64) error {
	if txInfo.From < minAccountIndex {
		return ErrFromAccountIndexTooLow
	}

	if txInfo.From > maxAccountIndex {
		return ErrFromAccountIndexTooHigh
	}

	if txInfo.To < minAccountIndex {
		return ErrToAccountIndexTooLow
	}

	if txInfo.To > maxAccountIndex {
		return ErrToAccountIndexTooHigh
	}

	if txInfo.Amount == nil {
		return ErrAmountNil
	}

	if txInfo.Amount.Cmp(minAmount) < 0 {
		return ErrAmountTooLow
	}

	if txInfo.Amount.Cmp(maxAmount) > 0 {
		return ErrAmountTooHigh
	}

	if txInfo.Fee == nil {
		return ErrFeeAmountNil
	}

	if txInfo.Fee.Cmp(minFeeAmount) < 0 {
		return ErrFeeAmountTooLow
	}

	if txInfo.Fee.Cmp(maxFeeAmount) > 0 {
		return ErrFeeAmountTooHigh
	}

	if txInfo.Nonce < minNonce {
		return ErrNonceTooLow
	}

	if txInfo.Nonce != fromAccountNonce {
		return ErrInvalidNonce
	}

	return nil
}

// TODO: check how to using iden3 MiMC lib to verify the signature
func (txInfo *TransactionInfo) VerifySignature(pubKey string) error {
	return nil
}
