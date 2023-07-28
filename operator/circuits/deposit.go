package circuits

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/types"
)

type DepositInputs struct {
	Account       *models.AccountDto
	DepositAmount *big.Int
	MTProof       *merkletree.CircomProcessorProof
}

type depositCircuitInputs struct {
	BalanceTreeRoot *merkletree.Hash    `json:"balanceTreeRoot"`
	PublicKey       [2]string           `json:"publicKey"`
	Balance         string              `json:"balance"`
	Nonce           string              `json:"nonce"`
	PathElements    [6]*merkletree.Hash `json:"pathElements"`
	OldKey          *merkletree.Hash    `json:"oldKey"`
	OldValue        *merkletree.Hash    `json:"oldValue"`
	IsOld0          string              `json:"isOld0"`
	NewKey          *merkletree.Hash    `json:"newKey"`
	Func            [2]string           `json:"func"`
}

func (d *DepositInputs) InputsMarshal() ([]byte, error) {
	circuitInputs := &depositCircuitInputs{
		BalanceTreeRoot: d.MTProof.OldRoot,
		Balance:         d.Account.Balance.String(),
		Nonce:           strconv.Itoa(int(d.Account.Nonce)),
		NewKey:          d.MTProof.NewKey,
		PathElements:    ([6]*merkletree.Hash)(d.MTProof.Siblings),
	}

	publicKey, err := tree.StringifyPublicKey(d.Account.PublicKey)
	if err != nil {
		return nil, err
	}
	circuitInputs.PublicKey = *publicKey

	if d.MTProof.IsOld0 {
		circuitInputs.IsOld0 = "1"
		circuitInputs.OldKey = &merkletree.HashZero
		circuitInputs.OldValue = &merkletree.HashZero
	} else {
		circuitInputs.IsOld0 = "0"
		circuitInputs.OldKey = d.MTProof.OldKey
		circuitInputs.OldValue = d.MTProof.OldValue
	}

	var op [2]string
	switch d.MTProof.Fnc {
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

type DepositOutputs struct {
	Proof         ProofData
	PublicSignals [5]*big.Int
}

func (d *DepositOutputs) OutputsUnmarshal(proof *types.ZKProof) error {
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
