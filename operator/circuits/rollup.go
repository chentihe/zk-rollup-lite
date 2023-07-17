package circuits

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/types"
)

type RollupInputs struct {
	Txs []*RollupTx
}

type RollupTx struct {
	Tx                                  *txmanager.TransactionInfo
	Root                                *merkletree.Hash
	Sender                              *AccountInfo
	Recipient                           *AccountInfo
	IntermediateBalanceTreeRoot         *merkletree.Hash
	IntermediateBalanceTreePathElements []*merkletree.Hash
}

type AccountInfo struct {
	Account      *models.Account
	PathElements []*merkletree.Hash
}

type rollupCircuitInputs struct {
	BalanceTreeRoots                     []string
	TxData                               [][]string
	TxSendersPublicKey                   [][2]string
	TxSendersBalance                     []string
	TxSendersNonce                       []string
	TxSendersPathElements                [][]string
	TxRecipientsPublicKey                [][2]string
	TxRecipientsBalance                  []string
	TxRecipientsNonce                    []string
	TxRecipientsPathElements             [][]string
	IntermediateBalanceTreeRoots         []string
	IntermediateBalanceTreesPathElements [][]string
}

func (r *RollupInputs) InputsMarshal() ([]byte, error) {
	var circuitInputs rollupCircuitInputs

	len := len(r.Txs)

	for i := 0; i < len; i++ {
		tx := r.Txs[i]

		senderPublicKey, err := accounttree.StringifyPublicKey(tx.Sender.Account.PublicKey)
		if err != nil {
			return nil, err
		}

		recipientPublicKey, err := accounttree.StringifyPublicKey(tx.Recipient.Account.PublicKey)
		if err != nil {
			return nil, err
		}

		circuitInputs.BalanceTreeRoots = append(circuitInputs.BalanceTreeRoots, tx.Root.String())
		circuitInputs.TxData = append(circuitInputs.TxData, tx.Tx.ToArray())

		// sender
		circuitInputs.TxSendersPublicKey = append(circuitInputs.TxSendersPublicKey, *senderPublicKey)
		circuitInputs.TxSendersBalance = append(circuitInputs.TxSendersBalance, tx.Sender.Account.Balance.String())
		circuitInputs.TxSendersNonce = append(circuitInputs.TxSendersNonce, strconv.Itoa(int(tx.Sender.Account.Nonce)))
		circuitInputs.TxSendersPathElements = append(circuitInputs.TxSendersPathElements, accounttree.StringifyPath(tx.Sender.PathElements))

		// recipient
		circuitInputs.TxRecipientsPublicKey = append(circuitInputs.TxRecipientsPublicKey, *recipientPublicKey)
		circuitInputs.TxRecipientsBalance = append(circuitInputs.TxRecipientsBalance, tx.Recipient.Account.Balance.String())
		circuitInputs.TxRecipientsNonce = append(circuitInputs.TxRecipientsNonce, strconv.Itoa(int(tx.Recipient.Account.Nonce)))
		circuitInputs.TxRecipientsPathElements = append(circuitInputs.TxRecipientsPathElements, accounttree.StringifyPath(tx.Sender.PathElements))

		// intermediate info
		circuitInputs.IntermediateBalanceTreeRoots = append(circuitInputs.IntermediateBalanceTreeRoots, tx.IntermediateBalanceTreeRoot.String())
		circuitInputs.IntermediateBalanceTreesPathElements = append(circuitInputs.IntermediateBalanceTreesPathElements, accounttree.StringifyPath(tx.IntermediateBalanceTreePathElements))
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
