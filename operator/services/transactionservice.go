package services

import (
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/txhandlers"
	"github.com/iden3/go-merkletree-sql/v2"
)

type TransactionService struct {
	AccountService AccountService
	MerkleTree     *merkletree.MerkleTree
}

func NewTransactionService(accountService *AccountService, merkleTree *merkletree.MerkleTree) *TransactionService {
	return &TransactionService{
		AccountService: *accountService,
		MerkleTree:     merkleTree,
	}
}

func (service *TransactionService) SendTransaction(tx *txhandlers.TransactionInfo) error {
	fromAccount, err := service.AccountService.GetAccountByIndex(tx.From)
	if err != nil {
		return err
	}

	toAccount, err := service.AccountService.GetAccountByIndex(tx.To)
	if err != nil {
		return err
	}

	fromAccount.Balance = new(big.Int).Sub(fromAccount.Balance, tx.Amount)
	fromAccount.Balance = new(big.Int).Sub(fromAccount.Balance, tx.Fee)
	toAccount.Balance = new(big.Int).Add(toAccount.Balance, tx.Amount)
	fromAccount.Nonce++

	if err = tx.Validate(fromAccount.Nonce); err != nil {
		return err
	}

	service.AccountService.UpdateAccount(fromAccount)
	service.AccountService.UpdateAccount(toAccount)

	// TODO: update merkle tree after the transaction

	return nil
}
