package txmanager

import (
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-rapidsnark/types"
)

type WithdrawInfo struct {
	AccountIndex   int64
	PublicKey      string
	Nullifier      *big.Int
	Signature      *babyjub.Signature
	WithdrawAmount *big.Int
	SignedTxHash   string
	ZkProof        *types.ZKProof
}

func (w *WithdrawInfo) VerifySignature() error {
	publicKey, err := tree.DecodePublicKeyFromString(w.PublicKey)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	if !publicKey.VerifyMimc7(w.Nullifier, w.Signature) {
		return ErrInvalidSignature
	}

	return nil
}
