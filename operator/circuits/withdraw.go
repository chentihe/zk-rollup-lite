package circuits

import (
	"encoding/json"
	"fmt"
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
	MTProof        *merkletree.CircomProcessorProof
}

type withdrawCircuitInputs struct {
	PublicKey       [2]string
	Signature       [3]string
	Nullifier       string
	BalanceTreeRoot *merkletree.Hash
	Balance         string
	Nonce           string
	PathElements    [6]*merkletree.Hash
	OldKey          *merkletree.Hash
	OldValue        *merkletree.Hash
	IsOld0          string
	NewKey          *merkletree.Hash
	Func            [2]string
}

func (w *WithdrawInputs) InputsMarshal() ([]byte, error) {
	circuitInputs := &withdrawCircuitInputs{
		Nullifier:       w.Nullifier.String(),
		BalanceTreeRoot: w.MTProof.OldRoot,
		Balance:         w.Account.Balance.String(),
		Nonce:           strconv.Itoa(int(w.Account.Nonce)),
		PathElements:    ([6]*merkletree.Hash)(w.MTProof.Siblings),
		OldKey:          w.MTProof.OldKey,
		OldValue:        w.MTProof.OldValue,
		IsOld0:          "0",
		NewKey:          w.MTProof.NewKey,
	}

	signature := [3]string{w.Signature.R8.X.String(), w.Signature.R8.Y.String(), w.Signature.S.String()}
	circuitInputs.Signature = signature

	publicKey, err := tree.StringifyPublicKey(w.Account.PublicKey)
	if err != nil {
		return nil, err
	}
	circuitInputs.PublicKey = *publicKey

	// withdraw only update mt
	var op [2]string
	switch w.MTProof.Fnc {
	case UPDATE:
		op = [2]string{"0", "1"}
	case INSERT, NOP, DELETE:
		return nil, fmt.Errorf("Should not indicate these functions")
	default:
		return nil, fmt.Errorf("Invalid function")
	}

	circuitInputs.Func = op

	return json.Marshal(circuitInputs)
}

type WithdrawOutputs struct {
	Proof         ProofData
	PublicSignals [5]*big.Int
}

func (w *WithdrawOutputs) OutputsUnmarshal(proof *types.ZKProof) error {
	inputs, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return err
	}

	w.PublicSignals = ([5]*big.Int)(inputs)
	if err = w.Proof.ProofUnmarshal(proof.Proof); err != nil {
		return err
	}

	return nil
}
