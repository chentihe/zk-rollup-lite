package services

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
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

func (service *TransactionService) SendTransaction(tx *txmanager.TransactionInfo) error {
	ctx := context.Background()

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

	fromAccount.Balance = new(big.Int).Sub(fromAccount.Balance, tx.Amount)
	fromAccount.Balance = new(big.Int).Sub(fromAccount.Balance, tx.Fee)
	toAccount.Balance = new(big.Int).Add(toAccount.Balance, tx.Amount)
	fromAccount.Nonce++

	service.AccountService.UpdateAccount(fromAccount)
	service.AccountService.UpdateAccount(toAccount)

	if err = tx.VerifySignature(fromAccount.PublicKey); err != nil {
		return err
	}

	fromLeaf, err := accounttree.GenerateAccountLeaf(fromAccount)
	if err != nil {
		return err
	}

	toLeaf, err := accounttree.GenerateAccountLeaf(toAccount)
	if err != nil {
		return err
	}

	if _, err := service.MerkleTree.Update(
		ctx,
		big.NewInt(fromAccount.AccountIndex),
		fromLeaf,
	); err != nil {
		return err
	}

	if _, err := service.MerkleTree.Update(
		ctx,
		big.NewInt(toAccount.AccountIndex),
		toLeaf,
	); err != nil {
		return err
	}

	// TODO: update the sent tx into redis
	// call rollup func once the tx amount is reaching 2

	return nil
}
