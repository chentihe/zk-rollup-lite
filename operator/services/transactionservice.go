package services

import (
	"context"
	"math"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-merkletree-sql/v2"
)

type TransactionService struct {
	accountService *AccountService
	accountTree    *tree.AccountTree
	redisCache     *cache.RedisCache
	ethClient      *ethclient.Client
	signer         *clients.Signer
	abi            *abi.ABI
	context        context.Context
}

func NewTransactionService(accountService *AccountService, tree *tree.AccountTree, cache *cache.RedisCache, ethClient *ethclient.Client, signer *clients.Signer, abi *abi.ABI, context context.Context) *TransactionService {
	return &TransactionService{
		accountService: accountService,
		accountTree:    tree,
		redisCache:     cache,
		ethClient:      ethClient,
		signer:         signer,
		abi:            abi,
		context:        context,
	}
}

// only verify the zkp since the signature is signed by user
func (service *TransactionService) Deposit(deposit *txmanager.DepositInfo) error {
	depositInputs := &circuits.DepositInputs{
		Root:          service.accountTree.GetRoot(),
		DepositAmount: deposit.DepositAmount,
	}

	var mtProof *merkletree.CircomProcessorProof

	account, err := service.accountService.GetAccountByIndex(deposit.AccountIndex)
	// add account into db & merkle tree if it's new account
	// event hanlder will update the rest info once the tx is on-chain
	if err == daos.ErrAccountNotFound {
		userIndex, err := service.accountService.GetCurrentAccountIndex()
		if err != nil {
			return err
		}

		account = &models.Account{
			AccountIndex: userIndex,
			PublicKey:    deposit.PublicKey,
			Balance:      big.NewInt(0),
			Nonce:        0,
		}

		if err := service.accountService.CreateAccount(account); err != nil {
			return err
		}

		leaf, err := tree.GenerateAccountLeaf(account)
		if err != nil {
			return err
		}

		mtProof, err = service.accountTree.AddAndGetCircomProof(userIndex, leaf)
		if err != nil {
			return err
		}
	} else {
		// mock update to get the circuit processor proof
		mtProof, err = service.accountTree.UpdateAccountTree(account)
		if err != nil {
			return err
		}
	}

	depositInputs.Account = account
	depositInputs.MTProof = mtProof

	circuitInput, err := depositInputs.InputsMarshal()
	if err != nil {
		return err
	}

	proof, err := circuits.GenerateGroth16Proof(circuitInput, circuitPath+"/deposit")
	if err != nil {
		return err
	}

	if err = circuits.VerifierGroth16(proof, circuitPath+"/deposit"); err != nil {
		return err
	}

	return nil
}

// only verify the zkp since the signature is signed by user
func (service *TransactionService) Withdraw(withdraw *txmanager.WithdrawInfo) error {
	if err := withdraw.VerifySignature(); err != nil {
		return err
	}

	withdrawInputs := &circuits.WithdrawInputs{
		Root:           service.accountTree.GetRoot(),
		Nullifier:      withdraw.Nullifier,
		Signature:      withdraw.Signature,
		WithdrawAmount: withdraw.WithdrawAmount,
	}

	account, err := service.accountService.GetAccountByIndex(withdraw.AccountIndex)
	if err != nil {
		return err
	}

	// mock update to get the circuit processor proof
	mtProof, err := service.accountTree.UpdateAccountTree(account)
	if err != nil {
		return err
	}

	withdrawInputs.Account = account
	withdrawInputs.MTProof = mtProof

	circuitInput, err := withdrawInputs.InputsMarshal()
	if err != nil {
		return err
	}

	proof, err := circuits.GenerateGroth16Proof(circuitInput, circuitPath+"/withdraw")
	if err != nil {
		return err
	}

	if err = circuits.VerifierGroth16(proof, circuitPath+"/withdraw"); err != nil {
		return err
	}

	var withdrawOutputs circuits.WithdrawOutputs
	if err = withdrawOutputs.OutputUnmarshal(proof); err != nil {
		return err
	}

	return nil
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

	if _, err := service.accountTree.UpdateAccountTree(fromAccount); err != nil {
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

	if _, err := service.accountTree.UpdateAccountTree(toAccount); err != nil {
		return math.MaxInt64, err
	}

	lastInsertedTx, err := service.redisCache.Get(service.context, lastInsertedKey, new(int))
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

	if err = service.redisCache.Set(service.context, lastInsertedKey, lastInsertedTx); err != nil {
		return math.MaxInt64, err
	}

	return lastInsertedTx.(int64), nil
}
