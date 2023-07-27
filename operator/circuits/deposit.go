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
	BalanceTreeRoot string    `json:"balanceTreeRoot"`
	PublicKey       [2]string `json:"publicKey"`
	Balance         string    `json:"balance"`
	Nonce           string    `json:"nonce"`
	PathElements    [6]string `json:"pathElements"`
	OldKey          string    `json:"oldKey"`
	OldValue        string    `json:"oldValue"`
	IsOld0          string    `json:"isOld0"`
	NewKey          string    `json:"newKey"`
	Func            [2]string `json:"func"`
}

func (d *DepositInputs) InputsMarshal() ([]byte, error) {
	circuitInputs := &depositCircuitInputs{
		BalanceTreeRoot: d.MTProof.OldRoot.BigInt().String(),
		Balance:         d.Account.Balance.String(),
		Nonce:           strconv.Itoa(int(d.Account.Nonce)),
		OldKey:          d.MTProof.OldKey.BigInt().String(),
		OldValue:        d.MTProof.OldValue.BigInt().String(),
		NewKey:          d.MTProof.NewKey.BigInt().String(),
	}

	publicKey, err := tree.StringifyPublicKey(d.Account.PublicKey)
	if err != nil {
		return nil, err
	}
	circuitInputs.PublicKey = *publicKey

	circuitInputs.PathElements = tree.StringifyPath(d.MTProof.Siblings)

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
