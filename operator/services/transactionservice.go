package services

import (
	"bytes"
	"context"
	"encoding/gob"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/redis/go-redis/v9"
)

type TransactionService struct {
	AccountService *AccountService
	AccountTree    *accounttree.AccountTree
	Cache          *redis.Client
}

func NewTransactionService(accountService *AccountService, accountTree *accounttree.AccountTree, cache *redis.Client) *TransactionService {
	return &TransactionService{
		AccountService: accountService,
		AccountTree:    accountTree,
		Cache:          cache,
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
	const lastInertedKey = "last-inserted"
	lastInsertedTx, err := service.Cache.Get(context, lastInertedKey).Int()
	if err != nil {
		return err
	}

	if lastInsertedTx == -1 {
		lastInsertedTx = 0
	}

	encodedBytes := new(bytes.Buffer)
	if err := gob.NewEncoder(encodedBytes).Encode(tx); err != nil {
		return err
	}

	service.Cache.Set(context, strconv.Itoa(lastInsertedTx), encodedBytes, 0)
	service.Cache.Set(context, lastInertedKey, lastInsertedTx+1, 0)

	return nil
}
