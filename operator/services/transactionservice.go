package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/contracts"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/pubsub"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-merkletree-sql/v2"
)

type TransactionService struct {
	AccountService *AccountService
	AccountTree    *accounttree.AccountTree
	RedisCache     *dbcache.RedisCache
	EthClient      *ethclient.Client
}

func NewTransactionService(accountService *AccountService, accountTree *accounttree.AccountTree, cache *dbcache.RedisCache, ethClient *ethclient.Client) *TransactionService {
	return &TransactionService{
		AccountService: accountService,
		AccountTree:    accountTree,
		RedisCache:     cache,
		EthClient:      ethClient,
	}
}

// user can deposit on their own or via our app
func (service *TransactionService) Deposit(deposit txmanager.DepositInfo) error {
	context := context.Background()

	var depositInputs circuits.DepositInputs
	depositInputs.Root = service.AccountTree.GetRoot()

	var mtProof *merkletree.CircomProcessorProof

	account, err := service.AccountService.GetAccountByIndex(deposit.AccountIndex)
	// add account into db & merkle tree if it's new account
	// event hanlder will update the rest info once the tx is on-chain
	if err == daos.ErrAccountNotFound {
		userIndex, err := service.AccountService.GetCurrentAccountIndex()
		if err != nil {
			return err
		}

		account = &models.Account{
			AccountIndex: account.AccountIndex,
			PublicKey:    deposit.PublicKey,
			Balance:      big.NewInt(0),
			Nonce:        0,
		}

		if err := service.AccountService.CreateAccount(account); err != nil {
			return err
		}

		leaf, err := accounttree.GenerateAccountLeaf(account)
		if err != nil {
			return err
		}

		mtProof, err = service.AccountTree.AddAndGetCircomProof(userIndex, leaf)
		if err != nil {
			return err
		}
	} else {
		// mock update to get the circuit processor proof
		mtProof, err = service.AccountTree.UpdateAccountTree(account)
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

	var depositOutputs circuits.DepositOutputs
	if err = depositOutputs.OutputUnmarshal(proof); err != nil {
		return err
	}

	// init smart contract
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return pubsub.ErrPubKeyToECDSA
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := service.EthClient.PendingNonceAt(context, fromAddress)
	if err != nil {
		return err
	}

	gasPrice, err := service.EthClient.SuggestGasPrice(context)
	if err != nil {
		return err
	}

	chainId, err := service.EthClient.ChainID(context)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// init rollup contract
	address := common.HexToAddress(rollupAddress)
	instance, err := contracts.NewRollup(address, service.EthClient)
	if err != nil {
		return err
	}

	tx, err := instance.Deposit(auth, depositOutputs.Proof.A, depositOutputs.Proof.B, depositOutputs.Proof.C, depositOutputs.PublicSignals)

	if err != nil {
		return err
	}

	fmt.Printf("Rollup finished: %v", tx)

	return nil
}

func (service *TransactionService) SendTransaction(tx *txmanager.TransactionInfo) error {
	context := context.Background()
	redisTxInfo := circuits.RollupTx{Root: service.AccountTree.GetRoot()}

	fromAccount, err := service.AccountService.GetAccountByIndex(tx.From)
	if err != nil {
		return err
	}

	redisTxInfo.Sender.Account = fromAccount
	senderPathElements, err := service.AccountTree.GetPathByAccount(fromAccount)
	if err != nil {
		return err
	}
	redisTxInfo.Sender.PathElements = senderPathElements

	toAccount, err := service.AccountService.GetAccountByIndex(tx.To)
	if err != nil {
		return err
	}

	redisTxInfo.Recipient.Account = toAccount
	recipientPathElements, err := service.AccountTree.GetPathByAccount(toAccount)
	if err != nil {
		return err
	}
	redisTxInfo.Recipient.PathElements = recipientPathElements

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

	if _, err := service.AccountTree.UpdateAccountTree(fromAccount); err != nil {
		return err
	}

	// update intermediate tree info
	intermediateBalanceTreePathElements, err := service.AccountTree.GetPathByAccount(toAccount)
	if err != nil {
		return err
	}
	redisTxInfo.IntermediateBalanceTreePathElements = intermediateBalanceTreePathElements
	redisTxInfo.IntermediateBalanceTreeRoot = service.AccountTree.GetRoot()

	if err := service.AccountService.UpdateAccount(toAccount); err != nil {
		return err
	}

	if _, err := service.AccountTree.UpdateAccountTree(toAccount); err != nil {
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
	service.RedisCache.Set(context, lastInsertedTx.(string), redisTxInfo)

	lastInsertedTx = lastInsertedTx.(int) + 1

	const rollUpCommand = "execute roll up"

	// TODO: move to config yaml
	const channel = "pendingTx"

	// TODO: add a subscriber to receive the rollup command
	// and execute rollup to L1
	if lastInsertedTx == 2 {
		lastInsertedTx = -1
		service.RedisCache.Publish(context, channel, rollUpCommand)
	}
	service.RedisCache.Set(context, lastInsertedKey, lastInsertedTx)

	return nil
}
