package circuits

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/types"
)

type DepositInputs struct {
	Account       *models.Account
	Root          *merkletree.Hash
	DepositAmount *big.Int
	MTProof       *merkletree.CircomProcessorProof
}

type depositCircuitInputs struct {
	BalanceTreeRoot string
	PublicKey       [2]string
	Balance         string
	Nonce           string
	PathElements    []string
	OldKey          string
	OldValue        string
	IsOld0          string
	NewKey          string
	Func            [2]string
}

const (
	NOP = iota
	UPDATE
	INSERT
	DELETE
)

func (d *DepositInputs) InputsMarshal() ([]byte, error) {
	circuitInputs := &depositCircuitInputs{
		BalanceTreeRoot: d.Root.String(),
		Balance:         d.Account.Balance.String(),
		Nonce:           strconv.Itoa(int(d.Account.Nonce)),
		OldKey:          d.MTProof.OldKey.String(),
		OldValue:        d.MTProof.OldValue.String(),
		NewKey:          d.MTProof.NewKey.String(),
	}

	publicKey, err := accounttree.StringifyPublicKey(d.Account.PublicKey)
	if err != nil {
		return nil, err
	}
	circuitInputs.PublicKey = *publicKey

	circuitInputs.PathElements = accounttree.StringifyPath(d.MTProof.Siblings)

	if d.MTProof.IsOld0 {
		circuitInputs.IsOld0 = "1"
	} else {
		circuitInputs.IsOld0 = "0"
	}

	var op [2]string
	switch d.MTProof.Fnc {
	case INSERT:
		op = [2]string{"1", "0"}
	case UPDATE:
		op = [2]string{"0", "1"}
	case NOP, DELETE:
		return nil, fmt.Errorf("Should not indicate these function")
	default:
		return nil, fmt.Errorf("Invalid function")
	}

	circuitInputs.Func = op

	return json.Marshal(circuitInputs)
}

type DepositOutputs struct {
	Proof         ProofData
	PublicSignals [5]*big.Int
}

func (d *DepositOutputs) OutputUnmarshal(proof *types.ZKProof) error {
	inputs, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return err
	}

	d.PublicSignals = ([5]*big.Int)(inputs)
	if err = d.Proof.ProofUnmarshal(proof.Proof); err != nil {
		return err
	}

	return nil
}
