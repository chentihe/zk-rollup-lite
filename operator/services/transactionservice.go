package services

import (
	"context"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
)

type TransactionService struct {
	AccountService *AccountService
	AccountTree    *accounttree.AccountTree
	RedisCache     *dbcache.RedisCache
}

func NewTransactionService(accountService *AccountService, accountTree *accounttree.AccountTree, cache *dbcache.RedisCache) *TransactionService {
	return &TransactionService{
		AccountService: accountService,
		AccountTree:    accountTree,
		RedisCache:     cache,
	}
}

func (service *TransactionService) SendTransaction(tx *txmanager.TransactionInfo) error {
	context := context.Background()

	fromAccount, err := service.AccountService.GetAccountByIndex(tx.From)
	if err != nil {
		return err
	}

	toAccount, err := service.AccountService.GetAccountByIndex(tx.To)
	if err != nil {
		return err
	}

	if err = tx.Validate(fromAccount.Nonce); err != nil {
		return err
	}

	if err = tx.VerifySignature(fromAccount.PublicKey); err != nil {
		return err
	}

	fromAccount.Balance = fromAccount.Balance.Sub(fromAccount.Balance, tx.Amount)
	fromAccount.Balance = fromAccount.Balance.Sub(fromAccount.Balance, tx.Fee)
	toAccount.Balance = toAccount.Balance.Add(toAccount.Balance, tx.Amount)
	fromAccount.Nonce++

	if err := service.AccountService.UpdateAccount(fromAccount); err != nil {
		return err
	}

	if err := service.AccountService.UpdateAccount(toAccount); err != nil {
		return err
	}

	if err := service.AccountTree.UpdateAccountTree(fromAccount); err != nil {
		return err
	}

	if err := service.AccountTree.UpdateAccountTree(toAccount); err != nil {
		return err
	}

	// TODO: update the sent tx into redis
	// call rollup func once the tx amount is reaching 2
	const lastInsertedKey = "last-inserted"
	lastInsertedTx, err := service.RedisCache.Get(context, lastInsertedKey, new(int))
	if err != nil {
		return err
	}

	// no pending transactions to roll up
	if lastInsertedTx == -1 {
		lastInsertedTx = 0
	}

	// encodedBytes := new(bytes.Buffer)
	// if err := gob.NewEncoder(encodedBytes).Encode(tx); err != nil {
	// 	return err
	// }
	service.RedisCache.Set(context, lastInsertedTx.(string), tx)

	lastInsertedTx = lastInsertedTx.(int) + 1

	const rollUpCommand = "execute roll up"

	// TODO: add a subscriber to receive the rollup command
	// and execute rollup to L1
	if lastInsertedTx == 2 {
		lastInsertedTx = -1
		service.RedisCache.Publish(context, rollUpCommand)
	}
	service.RedisCache.Set(context, lastInsertedKey, lastInsertedTx)

	return nil
}
