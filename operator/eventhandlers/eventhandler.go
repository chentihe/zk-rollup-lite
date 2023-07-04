package eventhandlers

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-merkletree-sql/v2"
)

type User struct {
	Index      int64    `json:"index"`
	PublicKeyX *big.Int `json:"publicKeyX"`
	PublicKeyY *big.Int `json:"publicKeyY"`
	Balance    *big.Int `json:"balance"`
	Nonce      int64    `json:"nonce"`
}

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

	contractAbi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return err
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	fmt.Println("Listening to events...")

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			switch vLog.Topics[0] {
			case depositHash:
				fmt.Println("Deposit Event")
				if err := AfterDeposit(&vLog, accountSerivce, mt, &contractAbi, ctx, client); err != nil {
					return err
				}
			case withdrawHash:
				fmt.Println("Withdraw Event")
				if err := AfterWithdraw(&vLog, accountSerivce, mt, &contractAbi, ctx); err != nil {
					return err
				}
			}
		}
	}
}
