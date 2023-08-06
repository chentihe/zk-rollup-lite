package circuits

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/types"
)

type RollupInputs struct {
	Txs []*ProcessTxInputs
}

type rollupCircuitInputs struct {
	BalanceTreeRoots                     []*merkletree.Hash    `json:"balanceTreeRoots"`
	TxData                               [][]string            `json:"txData"`
	TxSendersPublicKey                   [][2]string           `json:"txSendersPublicKey"`
	TxSendersBalance                     []string              `json:"txSendersBalance"`
	TxSendersNonce                       []string              `json:"txSendersNonce"`
	TxSendersPathElements                [][6]*merkletree.Hash `json:"txSendersPathElements"`
	TxRecipientsPublicKey                [][2]string           `json:"txRecipientsPublicKey"`
	TxRecipientsBalance                  []string              `json:"txRecipientsBalance"`
	TxRecipientsNonce                    []string              `json:"txRecipientsNonce"`
	TxRecipientsPathElements             [][6]*merkletree.Hash `json:"txRecipientsPathElements"`
	IntermediateBalanceTreeRoots         []*merkletree.Hash    `json:"intermediateBalanceTreeRoots"`
	IntermediateBalanceTreesPathElements [][6]*merkletree.Hash `json:"intermediateBalanceTreesPathElements"`
}

func (r *RollupInputs) InputsMarshal() ([]byte, error) {
	var circuitInputs rollupCircuitInputs

	len := len(r.Txs)

	for i := 0; i < len; i++ {
		tx := r.Txs[i]

		senderPublicKey, err := tree.StringifyPublicKey(tx.Sender.Account.PublicKey)
		if err != nil {
			return nil, err
		}

		recipientPublicKey, err := tree.StringifyPublicKey(tx.Recipient.Account.PublicKey)
		if err != nil {
			return nil, err
		}

		circuitInputs.BalanceTreeRoots = append(circuitInputs.BalanceTreeRoots, tx.Root)
		circuitInputs.TxData = append(circuitInputs.TxData, tx.Tx.ToArray())

		// sender
		circuitInputs.TxSendersPublicKey = append(circuitInputs.TxSendersPublicKey, *senderPublicKey)
		circuitInputs.TxSendersBalance = append(circuitInputs.TxSendersBalance, tx.Sender.Account.Balance.String())
		circuitInputs.TxSendersNonce = append(circuitInputs.TxSendersNonce, strconv.Itoa(int(tx.Sender.Account.Nonce)))
		circuitInputs.TxSendersPathElements = append(circuitInputs.TxSendersPathElements, ([6]*merkletree.Hash)(tx.Sender.PathElements))

		// recipient
		circuitInputs.TxRecipientsPublicKey = append(circuitInputs.TxRecipientsPublicKey, *recipientPublicKey)
		circuitInputs.TxRecipientsBalance = append(circuitInputs.TxRecipientsBalance, tx.Recipient.Account.Balance.String())
		circuitInputs.TxRecipientsNonce = append(circuitInputs.TxRecipientsNonce, strconv.Itoa(int(tx.Recipient.Account.Nonce)))
		circuitInputs.TxRecipientsPathElements = append(circuitInputs.TxRecipientsPathElements, ([6]*merkletree.Hash)(tx.Recipient.PathElements))

		// intermediate info
		circuitInputs.IntermediateBalanceTreeRoots = append(circuitInputs.IntermediateBalanceTreeRoots, tx.IntermediateBalanceTreeRoot)
		circuitInputs.IntermediateBalanceTreesPathElements = append(circuitInputs.IntermediateBalanceTreesPathElements, ([6]*merkletree.Hash)(tx.IntermediateBalanceTreePathElements))
	}

	return json.Marshal(circuitInputs)
}

type RollupOutputs struct {
	Proof         ProofData
	PublicSignals [19]*big.Int
}

func (r *RollupOutputs) OutputUnmarshal(proof *types.ZKProof) error {
	inputs, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return err
	}

	r.PublicSignals = ([19]*big.Int)(inputs)
	if err = r.Proof.ProofUnmarshal(proof.Proof); err != nil {
		return err
	}

	return nil
}
