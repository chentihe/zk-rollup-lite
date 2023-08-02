package services

import (
	"context"
	"encoding/json"
	"math"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
)

type TransactionService struct {
	accountService *AccountService
	accountTree    *tree.AccountTree
	redisCache     *cache.RedisCache
	context        context.Context
	keys           *config.Keys
}

func NewTransactionService(context context.Context, accountService *AccountService, tree *tree.AccountTree, cache *cache.RedisCache, keys *config.Keys) *TransactionService {
	return &TransactionService{
		accountService: accountService,
		accountTree:    tree,
		redisCache:     cache,
		context:        context,
		keys:           keys,
	}
}

func (service *TransactionService) SendTransaction(tx *txmanager.TransactionInfo) (int64, error) {
	// create a rollup tx object to save into redis
	redisTx := circuits.RollupTx{
		Tx:   tx,
		Root: service.accountTree.GetRoot(),
	}

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
	senderPathElements, err := service.accountTree.GetPathByAccount(fromAccount)
	if err != nil {
		return math.MaxInt64, err
	}
	redisTx.Sender = &circuits.AccountInfo{
		Account:      *fromAccount.Copy(),
		PathElements: senderPathElements,
	}

	toAccount, err := service.accountService.GetAccountByIndex(tx.To)
	if err != nil {
		return math.MaxInt64, err
	}

	recipientPathElements, err := service.accountTree.GetPathByAccount(toAccount)
	if err != nil {
		return math.MaxInt64, err
	}
	redisTx.Recipient = &circuits.AccountInfo{
		Account:      *toAccount.Copy(),
		PathElements: recipientPathElements,
	}

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

	value, err := service.redisCache.Get(service.context, service.keys.LastInsertedKey)
	if err != nil {
		return math.MaxInt64, err
	}

	lastInsertedTx, err := strconv.Atoi(value)
	if err != nil {
		return math.MaxInt64, err
	}

	// no pending transactions to roll up
	lastInsertedTx++
	serializedTx, err := json.Marshal(redisTx)
	if err != nil {
		return math.MaxInt64, err
	}

	if err = service.redisCache.Set(service.context, strconv.Itoa(lastInsertedTx), serializedTx); err != nil {
		return math.MaxInt64, err
	}

	if err = service.redisCache.Set(service.context, service.keys.LastInsertedKey, strconv.Itoa(lastInsertedTx)); err != nil {
		return math.MaxInt64, err
	}

	return (int64)(lastInsertedTx), nil
}
