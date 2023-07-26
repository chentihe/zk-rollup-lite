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
	Root           *merkletree.Hash
	WithdrawAmount *big.Int
	MTProof        *merkletree.CircomProcessorProof
}

type withdrawCircuitInputs struct {
	PublicKey       [2]string
	Signature       [3]string
	Nullifier       string
	BalanceTreeRoot string
	Balance         string
	Nonce           string
	PathElements    [6]string
	OldKey          string
	OldValue        string
	IsOld0          string
	NewKey          string
	Func            [2]string
}

func (w *WithdrawInputs) InputsMarshal() ([]byte, error) {
	circuitInputs := &withdrawCircuitInputs{
		Nullifier:       w.Nullifier.String(),
		BalanceTreeRoot: w.Root.String(),
		Balance:         w.Account.Balance.String(),
		Nonce:           strconv.Itoa(int(w.Account.Nonce)),
		OldKey:          w.MTProof.OldKey.String(),
		OldValue:        w.MTProof.OldValue.String(),
		NewKey:          w.MTProof.NewKey.String(),
	}

	signature := [3]string{w.Signature.R8.X.String(), w.Signature.R8.Y.String(), w.Signature.S.String()}
	circuitInputs.Signature = signature

	publicKey, err := tree.StringifyPublicKey(w.Account.PublicKey)
	if err != nil {
		return nil, err
	}
	circuitInputs.PublicKey = *publicKey

	circuitInputs.PathElements = tree.StringifyPath(w.MTProof.Siblings)

	if w.MTProof.IsOld0 {
		circuitInputs.IsOld0 = "1"
	} else {
		circuitInputs.IsOld0 = "0"
	}

	var op [2]string
	switch w.MTProof.Fnc {
	case INSERT:
		op = [2]string{"1", "0"}
	case UPDATE:
		op = [2]string{"0", "1"}
	case NOP, DELETE:
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
