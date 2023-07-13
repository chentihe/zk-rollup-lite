package zeroknowledge

import (
	"context"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/chentihe/zk-rollup-lite/operator/services"
)

type CircuitInput struct {
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

func GenerateCircuitInput(numbersOfTx int, redisCache *dbcache.RedisCache, accountService *services.AccountService) (*CircuitInput, error) {
	circuitInput := CircuitInput{}

	for i := 0; i < numbersOfTx; i++ {
		object, err := redisCache.Get(context.Background(), strconv.Itoa(i), new(services.RedisTxInfo))
		if err != nil {
			return nil, err
		}

		tx, ok := object.(services.RedisTxInfo)
		if !ok {
			return nil, ErrTx
		}

		senderPublicKey, err := accounttree.StringifyPublicKey(tx.Sender.Account.PublicKey)
		if err != nil {
			return nil, err
		}

		recipientPublicKey, err := accounttree.StringifyPublicKey(tx.Recipient.Account.PublicKey)
		if err != nil {
			return nil, err
		}

		circuitInput.BalanceTreeRoots = append(circuitInput.BalanceTreeRoots, tx.Root.String())
		circuitInput.TxData = append(circuitInput.TxData, tx.Tx.ToArray())

		// sender
		circuitInput.TxSendersPublicKey = append(circuitInput.TxSendersPublicKey, *senderPublicKey)
		circuitInput.TxSendersBalance = append(circuitInput.TxSendersBalance, tx.Sender.Account.Balance.String())
		circuitInput.TxSendersNonce = append(circuitInput.TxSendersNonce, strconv.Itoa(int(tx.Sender.Account.Nonce)))
		circuitInput.TxSendersPathElements = append(circuitInput.TxSendersPathElements, stringifyPathElements(tx.Sender.PathElements))

		// recipient
		circuitInput.TxRecipientsPublicKey = append(circuitInput.TxRecipientsPublicKey, *recipientPublicKey)
		circuitInput.TxRecipientsBalance = append(circuitInput.TxRecipientsBalance, tx.Recipient.Account.Balance.String())
		circuitInput.TxRecipientsNonce = append(circuitInput.TxRecipientsNonce, strconv.Itoa(int(tx.Recipient.Account.Nonce)))
		circuitInput.TxRecipientsPathElements = append(circuitInput.TxRecipientsPathElements, stringifyPathElements(tx.Sender.PathElements))

		// intermediate info
		circuitInput.IntermediateBalanceTreeRoots = append(circuitInput.IntermediateBalanceTreeRoots, tx.IntermediateBalanceTreeRoot.String())
		circuitInput.IntermediateBalanceTreesPathElements = append(circuitInput.IntermediateBalanceTreesPathElements, stringifyPathElements(tx.IntermediateBalanceTreePathElements))
	}

	return &circuitInput, nil
}

func stringifyPathElements(pathElements []*big.Int) []string {
	path := []string{}
	for _, pathElement := range pathElements {
		path = append(path, pathElement.String())
	}
	return path
}
