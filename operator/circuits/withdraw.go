package circuits

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/types"
)

type WithdrawInputs struct {
	Account        *models.AccountDto
	Nullifier      *big.Int
	Signature      *babyjub.Signature
	WithdrawAmount *big.Int
	MTProof        *merkletree.CircomVerifierProof
}

type withdrawCircuitInputs struct {
	BalanceTreeRoot *merkletree.Hash    `json:"balanceTreeRoot"`
	Signature       [3]string           `json:"signature"`
	Nullifier       string              `json:"nullifier"`
	PublicKey       [2]string           `json:"publicKey"`
	Balance         string              `json:"balance"`
	Nonce           string              `json:"nonce"`
	PathElements    [6]*merkletree.Hash `json:"pathElements"`
	OldKey          *merkletree.Hash    `json:"oldKey"`
	OldValue        *merkletree.Hash    `json:"oldValue"`
	NewKey          *merkletree.Hash    `json:"newKey"`
}

func (w *WithdrawInputs) InputsMarshal() ([]byte, error) {
	circuitInputs := &withdrawCircuitInputs{
		BalanceTreeRoot: w.MTProof.Root,
		Nullifier:       w.Nullifier.String(),
		Balance:         w.Account.Balance.String(),
		Nonce:           strconv.Itoa(int(w.Account.Nonce)),
		PathElements:    ([6]*merkletree.Hash)(w.MTProof.Siblings),
		OldKey:          w.MTProof.OldKey,
		OldValue:        w.MTProof.OldValue,
		NewKey:          w.MTProof.Key,
	}

	signature := [3]string{w.Signature.R8.X.String(), w.Signature.R8.Y.String(), w.Signature.S.String()}
	circuitInputs.Signature = signature

	publicKey, err := tree.StringifyPublicKey(w.Account.PublicKey)
	if err != nil {
		return nil, err
	}
	circuitInputs.PublicKey = *publicKey

	return json.Marshal(circuitInputs)
}

type WithdrawOutputs struct {
	Proof         ProofData
	PublicSignals [3]*big.Int
}

func (w *WithdrawOutputs) OutputsUnmarshal(proof *types.ZKProof) error {
	inputs, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return err
	}

	w.PublicSignals = ([3]*big.Int)(inputs)
	if err = w.Proof.ProofUnmarshal(proof.Proof); err != nil {
		return err
	}

	return nil
}
