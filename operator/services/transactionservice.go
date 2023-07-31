package services

import (
	"context"
	"math"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type TransactionService struct {
	accountService *AccountService
	accountTree    *tree.AccountTree
	redisCache     *cache.RedisCache
	signer         *clients.Signer
	abi            *abi.ABI
	context        context.Context
	circuitPath    string
	keys           *config.Keys
}

func NewTransactionService(accountService *AccountService, tree *tree.AccountTree, cache *cache.RedisCache, signer *clients.Signer, abi *abi.ABI, context context.Context, circuitPath string, keys *config.Keys) *TransactionService {
	return &TransactionService{
		accountService: accountService,
		accountTree:    tree,
		redisCache:     cache,
		signer:         signer,
		abi:            abi,
		context:        context,
		circuitPath:    circuitPath,
		keys:           keys,
	}
}

func (service *TransactionService) SendTransaction(tx *txmanager.TransactionInfo) (int64, error) {
	// create a rollup tx object to save into redis
	redisTx := circuits.RollupTx{Root: service.accountTree.GetRoot()}

	fromAccount, err := service.accountService.GetAccountByIndex(tx.From)
	if err != nil {
		return math.MaxInt64, err
	}

	// validate txinfo
	if err = tx.Validate(fromAccount.Nonce); err != nil {
		return math.MaxInt64, err
	}

	// validate signature
	if err = tx.VerifySignature(fromAccount.PublicKey); err != nil {
		return math.MaxInt64, err
	}

	// set sender data into rollup tx
	redisTx.Sender.Account = fromAccount
	senderPathElements, err := service.accountTree.GetPathByAccount(fromAccount)
	if err != nil {
		return math.MaxInt64, err
	}
	redisTx.Sender.PathElements = senderPathElements

	toAccount, err := service.accountService.GetAccountByIndex(tx.To)
	if err != nil {
		return math.MaxInt64, err
	}

	redisTx.Recipient.Account = toAccount
	recipientPathElements, err := service.accountTree.GetPathByAccount(toAccount)
	if err != nil {
		return math.MaxInt64, err
	}
	redisTx.Recipient.PathElements = recipientPathElements

	// update sender balance & nonce
	fromAccount.Balance = fromAccount.Balance.Sub(fromAccount.Balance, tx.Amount)
	fromAccount.Balance = fromAccount.Balance.Sub(fromAccount.Balance, tx.Fee)
	fromAccount.Nonce++

	if err := service.accountService.UpdateAccount(fromAccount); err != nil {
		return math.MaxInt64, err
	}

	if _, err := service.accountTree.UpdateAccount(fromAccount); err != nil {
		return math.MaxInt64, err
	}

	// update intermediate tree info
	intermediateBalanceTreePathElements, err := service.accountTree.GetPathByAccount(toAccount)
	if err != nil {
		return math.MaxInt64, err
	}
	redisTx.IntermediateBalanceTreePathElements = intermediateBalanceTreePathElements
	redisTx.IntermediateBalanceTreeRoot = service.accountTree.GetRoot()

	// update recipient balance
	toAccount.Balance = toAccount.Balance.Add(toAccount.Balance, tx.Amount)

	if err := service.accountService.UpdateAccount(toAccount); err != nil {
		return math.MaxInt64, err
	}

	if _, err := service.accountTree.UpdateAccount(toAccount); err != nil {
		return math.MaxInt64, err
	}

	lastInsertedTx, err := service.redisCache.Get(service.context, service.keys.LastInsertedKey, new(int))
	if err != nil {
		return math.MaxInt64, err
	}

	// no pending transactions to roll up
	if lastInsertedTx == -1 {
		lastInsertedTx = 0
	}

	// encodedBytes := new(bytes.Buffer)
	// if err := gob.NewEncoder(encodedBytes).Encode(tx); err != nil {
	// 	return err
	// }
	if err = service.redisCache.Set(service.context, lastInsertedTx.(string), redisTx); err != nil {
		return math.MaxInt64, err
	}

	lastInsertedTx = lastInsertedTx.(int) + 1

	// send rollup command to channel
	// reset the inserted tx to -1
	if lastInsertedTx == 2 {
		lastInsertedTx = -1

		if err = service.redisCache.Set(service.context, lastInsertedTx.(string), redisTx); err != nil {
			return math.MaxInt64, err
		}
	}

	if err = service.redisCache.Set(service.context, service.keys.LastInsertedKey, lastInsertedTx); err != nil {
		return math.MaxInt64, err
	}

	return lastInsertedTx.(int64), nil
}
