package services

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/chentihe/zk-rollup-lite/operator/txutils"
)

type TransactionService struct {
	accountService *AccountService
	accountTree    *tree.AccountTree
	redisCache     *cache.RedisCache
	context        context.Context
	keys           *config.Keys
	circuitPath    string
}

func NewTransactionService(context context.Context, accountService *AccountService, tree *tree.AccountTree, cache *cache.RedisCache, keys *config.Keys, circuitPath string) *TransactionService {
	return &TransactionService{
		accountService: accountService,
		accountTree:    tree,
		redisCache:     cache,
		context:        context,
		keys:           keys,
		circuitPath:    circuitPath + "/tx",
	}
}

func (service *TransactionService) SendTransaction(tx *txutils.TransactionInfo) error {
	// create a process tx inputs object to save into redis
	processTxInputs := &circuits.ProcessTxInputs{
		Tx:   tx,
		Root: service.accountTree.GetRoot(),
	}

	fromAccount, err := service.accountService.GetAccountByIndex(tx.From)
	if err != nil {
		return err
	}
	cpFrom := fromAccount.Copy()

	// validate txinfo
	if err = tx.Validate(fromAccount.Nonce); err != nil {
		return err
	}

	// validate signature
	if err = tx.VerifySignature(fromAccount.PublicKey); err != nil {
		return err
	}

	// set sender data into rollup tx
	senderPathElements, err := service.accountTree.GetPathByAccount(fromAccount)
	if err != nil {
		return err
	}
	processTxInputs.Sender = &circuits.AccountInfo{
		Account:      *cpFrom,
		PathElements: senderPathElements,
	}

	toAccount, err := service.accountService.GetAccountByIndex(tx.To)
	if err != nil {
		return err
	}
	cpTo := toAccount.Copy()

	// set recipient data into rollup tx
	recipientPathElements, err := service.accountTree.GetPathByAccount(toAccount)
	if err != nil {
		return err
	}
	processTxInputs.Recipient = &circuits.AccountInfo{
		Account:      *cpTo,
		PathElements: recipientPathElements,
	}

	// update sender balance & nonce
	fromAccount.Balance = fromAccount.Balance.Sub(fromAccount.Balance, tx.Amount)
	fromAccount.Balance = fromAccount.Balance.Sub(fromAccount.Balance, tx.Fee)
	fromAccount.Nonce++

	if err := service.accountService.UpdateAccount(fromAccount); err != nil {
		return err
	}

	if _, err := service.accountTree.UpdateAccount(fromAccount); err != nil {
		return err
	}

	// update intermediate tree info
	intermediateBalanceTreePathElements, err := service.accountTree.GetPathByAccount(toAccount)
	if err != nil {
		return err
	}
	processTxInputs.IntermediateBalanceTreePathElements = intermediateBalanceTreePathElements
	processTxInputs.IntermediateBalanceTreeRoot = service.accountTree.GetRoot()

	// update recipient balance
	toAccount.Balance = toAccount.Balance.Add(toAccount.Balance, tx.Amount)

	if err := service.accountService.UpdateAccount(toAccount); err != nil {
		return err
	}

	if _, err := service.accountTree.UpdateAccount(toAccount); err != nil {
		return err
	}

	circuitInput, err := processTxInputs.InputsMarshal()
	if err != nil {
		return err
	}

	// if errors occur during zkp process, roll back the account table & mt
	proof, err := circuits.GenerateGroth16Proof(circuitInput, service.circuitPath)
	if err != nil {
		service.rollback(cpFrom, cpTo)
		return err
	}

	if err = circuits.VerifierGroth16(proof, service.circuitPath); err != nil {
		service.rollback(cpFrom, cpTo)
		return err
	}

	// update the value of last inserted key
	value, err := service.redisCache.Get(service.context, service.keys.LastInsertedKey)
	if err != nil {
		return err
	}

	lastInsertedTx, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	lastInsertedTx++
	serializedTx, err := json.Marshal(processTxInputs)
	if err != nil {
		return err
	}

	// save tx into redis
	if err = service.redisCache.Set(service.context, strconv.Itoa(lastInsertedTx), serializedTx); err != nil {
		return err
	}

	if err = service.redisCache.Set(service.context, service.keys.LastInsertedKey, strconv.Itoa(lastInsertedTx)); err != nil {
		return err
	}

	return nil
}

func (service *TransactionService) rollback(from *models.AccountDto, to *models.AccountDto) {
	service.accountService.UpdateAccount(from)
	service.accountService.UpdateAccount(to)
	service.accountTree.UpdateAccount(from)
	service.accountTree.UpdateAccount(to)
}
