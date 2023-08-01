package txmanager

import (
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

type TransactionInfo struct {
	From      int64              `json:"from"`
	To        int64              `json:"to"`
	Amount    *big.Int           `json:"amount"`
	Fee       *big.Int           `json:"fee"`
	Nonce     int64              `json:"nonce"`
	Signature *babyjub.Signature `json:"signature"`
}

func (txInfo *TransactionInfo) ToArray() []string {
	stringFrom := strconv.Itoa(int(txInfo.From))
	stringTo := strconv.Itoa(int(txInfo.To))
	stringAmount := txInfo.Amount.String()
	stringFee := txInfo.Fee.String()
	stringNonce := strconv.Itoa(int(txInfo.Nonce))
	stringR8X := txInfo.Signature.R8.X.String()
	stringR8Y := txInfo.Signature.R8.Y.String()
	stringS := txInfo.Signature.S.String()
	return []string{stringFrom, stringTo, stringAmount, stringFee, stringNonce, stringR8X, stringR8Y, stringS}
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

	// fix the fee
	if txInfo.Fee != Fee {
		txInfo.Fee = Fee
	}

	if txInfo.Nonce < minNonce {
		return ErrNonceTooLow
	}

	if txInfo.Nonce != fromAccountNonce+1 {
		return ErrInvalidNonce
	}

	return nil
}

func (txInfo *TransactionInfo) HashMsg() (*big.Int, error) {
	return poseidon.Hash([]*big.Int{
		big.NewInt(txInfo.From),
		big.NewInt(txInfo.To),
		txInfo.Amount,
		txInfo.Fee,
		big.NewInt(txInfo.Nonce),
	})
}

func (txInfo *TransactionInfo) VerifySignature(comp string) error {
	publicKey, err := tree.DecodePublicKeyFromString(comp)
	if err != nil {
		return err
	}

	hashedMsg, err := txInfo.HashMsg()
	if err != nil {
		return err
	}

	if !publicKey.VerifyMimc7(hashedMsg, txInfo.Signature) {
		return ErrInvalidSignature
	}

	return nil
}
