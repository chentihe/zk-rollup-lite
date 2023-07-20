package eventhandlers

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/contracts"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type User struct {
	Index      int64    `json:"index"`
	PublicKeyX *big.Int `json:"publicKeyX"`
	PublicKeyY *big.Int `json:"publicKeyY"`
	Balance    *big.Int `json:"balance"`
	Nonce      int64    `json:"nonce"`
}

type EventHandler struct {
	context        context.Context
	ethClient      *ethclient.Client
	query          ethereum.FilterQuery
	logs           chan types.Log
	sub            ethereum.Subscription
	abi            abi.ABI
	accountService *services.AccountService
	accountTree    *tree.AccountTree
}

func NewEventHandler(context context.Context, svc *config.ServiceContext) (*EventHandler, error) {
	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	contractAbi, err := abi.JSON(strings.NewReader(contracts.RollupABI))
	if err != nil {
		return nil, err
	}

	logs := make(chan types.Log)

	sub, err := svc.EthClient.SubscribeFilterLogs(context, query, logs)
	if err != nil {
		return nil, err
	}

	return &EventHandler{
		context:        context,
		query:          query,
		logs:           logs,
		sub:            sub,
		abi:            contractAbi,
		accountService: svc.AccountService,
		accountTree:    svc.AccountTree,
	}, nil
}

func (e *EventHandler) Listening() {
	fmt.Println("Listening to events...")
	for {
		select {
		case err := <-e.sub.Err():
			fmt.Printf("Subscription err: %v", err)
		case vLog := <-e.logs:
			switch vLog.Topics[0] {
			case depositHash:
				fmt.Println("Deposit Event")
				if err := e.afterDeposit(&vLog); err != nil {
					fmt.Printf("Deposit event err: %v", err)
				}
			case withdrawHash:
				fmt.Println("Withdraw Event")
				if err := e.afterWithdraw(&vLog); err != nil {
					fmt.Printf("Withdraw event err: %v", err)
				}
			}
		}
	}
}
