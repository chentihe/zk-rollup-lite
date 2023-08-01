package eventhandler

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/layer1"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/spf13/viper"
)

type EventHandler struct {
	context   context.Context
	ethClient *ethclient.Client
	query     ethereum.FilterQuery
	// logs           chan types.Log
	// sub            ethereum.Subscription
	abi            *abi.ABI
	accountService *services.AccountService
	accountTree    *tree.AccountTree
	config         *config.Config
}

func NewEventHandler(context context.Context, accountService *services.AccountService, accountTree *tree.AccountTree, ethClient *ethclient.Client, abi *abi.ABI, config *config.Config) (*EventHandler, error) {
	rollupAddress := common.HexToAddress(config.SmartContract.Address)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{rollupAddress},
	}

	return &EventHandler{
		context:        context,
		query:          query,
		abi:            abi,
		ethClient:      ethClient,
		accountService: accountService,
		accountTree:    accountTree,
		config:         config,
	}, nil
}

func (e *EventHandler) Listening() {
	fmt.Println("Listening to events...")

	logs := make(chan types.Log)

	sub, err := e.ethClient.SubscribeFilterLogs(e.context, e.query, logs)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Subscription err: %v\n", err)
			case vLog := <-logs:
				switch vLog.Topics[0] {
				case depositHash:
					log.Println("Deposit Event")
					if err := e.afterDeposit(&vLog); err != nil {
						log.Printf("Deposit event err: %v", err)
					}
				case withdrawHash:
					log.Println("Withdraw Event")
					if err := e.afterWithdraw(&vLog); err != nil {
						log.Printf("Withdraw event err: %v", err)
					}
				}
			}
		}
	}()
}

func (e *EventHandler) afterWithdraw(vLog *types.Log) error {
	var withdraw layer1.Withdraw
	if err := e.abi.UnpackIntoInterface(&withdraw, "Withdraw", vLog.Data); err != nil {
		return err
	}

	user := withdraw.User
	// user must be found in the db
	accountDto, err := e.accountService.GetAccountByIndex(user.Index.Int64())
	if err != nil {
		return err
	}

	accountDto.Balance = user.Balance
	accountDto.Nonce = user.Nonce.Int64()

	if err := e.accountService.UpdateAccount(accountDto); err != nil {
		return err
	}

	if _, err := e.accountTree.UpdateAccount(accountDto); err != nil {
		return err
	}

	log.Printf("Withdraw account: %#v\n", accountDto)

	return nil
}

func (e *EventHandler) afterDeposit(vLog *types.Log) error {
	var deposit layer1.Deposit
	if err := e.abi.UnpackIntoInterface(&deposit, "Deposit", vLog.Data); err != nil {
		return err
	}

	user := deposit.User
	accountDto, err := e.accountService.GetAccountByIndex(user.Index.Int64())
	switch err {
	case daos.ErrAccountNotFound:
		// retrieve sender address from tx
		tx, _, err := e.ethClient.TransactionByHash(e.context, vLog.TxHash)
		if err != nil {
			return err
		}

		sender, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
		if err != nil {
			return err
		}

		publicKey := babyjub.PublicKey(babyjub.Point{X: user.PublicKeyX, Y: user.PublicKeyY}).String()

		accountDto = &models.AccountDto{
			AccountIndex: user.Index.Int64(),
			PublicKey:    publicKey,
			Balance:      user.Balance,
			Nonce:        user.Nonce.Int64(),
			L1Address:    sender.Hex(),
		}

		if err = e.accountService.CreateAccount(accountDto); err != nil {
			return err
		}

		if _, err = e.accountTree.AddAccount(accountDto); err != nil {
			return err
		}

		// update account index on env.yaml
		for i, account := range e.config.Accounts {
			var k babyjub.PrivateKey
			_, err := hex.Decode(k[:], []byte(account.EddsaPrivKey))
			if err != nil {
				return err
			}

			if strings.Compare(k.Public().String(), publicKey) == 0 {
				e.config.Accounts[i].Index = user.Index.Int64()
			}
		}
		viper.Set("accounts", e.config.Accounts)
		viper.WriteConfig()
	default:
		accountDto.Balance = user.Balance
		accountDto.Nonce = user.Nonce.Int64()
		if err = e.accountService.UpdateAccount(accountDto); err != nil {
			return err
		}

		if _, err = e.accountTree.UpdateAccount(accountDto); err != nil {
			return err
		}
	}
	log.Printf("Deposit account: %#v\n", accountDto)

	return nil
}
