package circuits

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/chentihe/zk-rollup-lite/operator/txutils"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/types"
)

type ProcessTxInputs struct {
	Tx                                  *txutils.TransactionInfo
	Root                                *merkletree.Hash
	Sender                              *AccountInfo
	Recipient                           *AccountInfo
	IntermediateBalanceTreeRoot         *merkletree.Hash
	IntermediateBalanceTreePathElements []*merkletree.Hash
}

type AccountInfo struct {
	Account      models.AccountDto
	PathElements []*merkletree.Hash
}

type processTxCircuitInputs struct {
	BalanceTreeRoot                     *merkletree.Hash    `json:"balanceTreeRoot"`
	TxData                              []string            `json:"txData"`
	TxSenderPublicKey                   [2]string           `json:"txSenderPublicKey"`
	TxSenderBalance                     string              `json:"txSenderBalance"`
	TxSenderNonce                       string              `json:"txSenderNonce"`
	TxSenderPathElements                [6]*merkletree.Hash `json:"txSenderPathElements"`
	TxRecipientPublicKey                [2]string           `json:"txRecipientPublicKey"`
	TxRecipientBalance                  string              `json:"txRecipientBalance"`
	TxRecipientNonce                    string              `json:"txRecipientNonce"`
	TxRecipientPathElements             [6]*merkletree.Hash `json:"txRecipientPathElements"`
	IntermediateBalanceTreeRoot         *merkletree.Hash    `json:"intermediateBalanceTreeRoot"`
	IntermediateBalanceTreePathElements [6]*merkletree.Hash `json:"intermediateBalanceTreePathElements"`
}

func (p *ProcessTxInputs) InputsMarshal() ([]byte, error) {
	var circuitInputs processTxCircuitInputs

	senderPublicKey, err := tree.StringifyPublicKey(p.Sender.Account.PublicKey)
	if err != nil {
		return nil, err
	}

	recipientPublicKey, err := tree.StringifyPublicKey(p.Recipient.Account.PublicKey)
	if err != nil {
		return nil, err
	}

	circuitInputs.BalanceTreeRoot = p.Root
	circuitInputs.TxData = p.Tx.ToArray()

	// sender
	circuitInputs.TxSenderPublicKey = *senderPublicKey
	circuitInputs.TxSenderBalance = p.Sender.Account.Balance.String()
	circuitInputs.TxSenderNonce = strconv.Itoa(int(p.Sender.Account.Nonce))
	circuitInputs.TxSenderPathElements = ([6]*merkletree.Hash)(p.Sender.PathElements)

	// recipient
	circuitInputs.TxRecipientPublicKey = *recipientPublicKey
	circuitInputs.TxRecipientBalance = p.Recipient.Account.Balance.String()
	circuitInputs.TxRecipientNonce = strconv.Itoa(int(p.Recipient.Account.Nonce))
	circuitInputs.TxRecipientPathElements = ([6]*merkletree.Hash)(p.Recipient.PathElements)

	// intermediate info
	circuitInputs.IntermediateBalanceTreeRoot = p.IntermediateBalanceTreeRoot
	circuitInputs.IntermediateBalanceTreePathElements = ([6]*merkletree.Hash)(p.IntermediateBalanceTreePathElements)

	return json.Marshal(circuitInputs)
}

type ProcessTxOutputs struct {
	Proof         ProofData
	PublicSignals [1]*big.Int
}

func (p *ProcessTxOutputs) OutputUnmarshal(proof *types.ZKProof) error {
	inputs, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return err
	}

	p.PublicSignals = ([1]*big.Int)(inputs)
	if err = p.Proof.ProofUnmarshal(proof.Proof); err != nil {
		return err
	}

	return nil
}
