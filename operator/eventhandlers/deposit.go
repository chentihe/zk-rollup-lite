package eventhandlers

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/database"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
)

func initEthClient() (*ethclient.Client, error) {
	// TODO: change to yaml config
	client, err := ethclient.Dial(os.Getenv("WEBSOCKER_URL"))

	if err != nil {
		return nil, err
	}

	return client, nil
}

func updateMerletree(accountSerivce *services.AccountService, mt *merkletree.MerkleTree) error {
	ctx := context.Background()

	client, err := initEthClient()
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	contractAbi, err := abi.JSON(strings.NewReader(os.Getenv("CONTRACT_ABI")))
	if err != nil {
		return err
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	depositHash := crypto.Keccak256Hash([]byte("Deposit((uint256,uint256,uint256,uint256,uint256))"))

	withdrawHash := crypto.Keccak256Hash([]byte("Withdraw((uint256,uint256,uint256,uint256,uint256))"))

	fmt.Println("Listening to events...")

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			fmt.Println("vLog: ", vLog)

			switch vLog.Topics[0] {
			case depositHash:
				fmt.Println("Deposit Event")
				event := map[string]interface{}{}
				inerr := contractAbi.UnpackIntoMap(event, "Deposit", vLog.Data)
				if inerr != nil {
					return err
				}
				index := event["index"].(int64)
				publicKeyX := event["publicKeyX"].(*big.Int)
				publicKeyY := event["publicKeyY"].(*big.Int)
				balance := event["balance"].(*big.Int)
				nonce := event["nonce"].(int64)

				publicKey := babyjub.PublicKey(babyjub.Point{X: publicKeyX, Y: publicKeyY})

				account, err := accountSerivce.GetAccountByIndex(index)
				if err != nil {
					account = &models.AccountModel{
						AccountIndex: index,
						PublicKey:    publicKey.String(),
						Balance:      balance,
						Nonce:        nonce,
					}

					accountLeaf, err := database.GenerateAccountLeaf(account)
					if err != nil {
						return err
					}
					accountSerivce.CreateAccount(account)

					mt.Add(ctx, big.NewInt(index), accountLeaf)
				} else {
					account.Balance = balance
					account.Nonce = nonce
					accountLeaf, err := database.GenerateAccountLeaf(account)
					if err != nil {
						return err
					}
					accountSerivce.UpdateAccount(account)

					mt.Update(ctx, big.NewInt(index), accountLeaf)
				}
			case withdrawHash:
				fmt.Println("Withdraw Event")
				event := map[string]interface{}{}
				inerr := contractAbi.UnpackIntoMap(event, "Withdraw", vLog.Data)
				if inerr != nil {
					return err
				}
				index := event["index"].(int64)
				balance := event["balance"].(*big.Int)
				nonce := event["nonce"].(int64)

				account, err := accountSerivce.GetAccountByIndex(index)
				if err != nil {
					return fmt.Errorf("Invalid account")
				} else {
					account.Balance = balance
					account.Nonce = nonce
					accountLeaf, err := database.GenerateAccountLeaf(account)
					if err != nil {
						return err
					}
					accountSerivce.UpdateAccount(account)

					mt.Update(ctx, big.NewInt(index), accountLeaf)
				}
			}
		}
	}
}
