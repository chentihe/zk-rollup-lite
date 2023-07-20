package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/controllers"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/db"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/eventhandler"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type ServiceContext struct {
	PostgresDB            *gorm.DB
	Redis                 *dbcache.RedisCache
	AccountTree           *tree.AccountTree
	EthClient             *ethclient.Client
	AccountService        *services.AccountService
	AccountController     *controllers.AccountController
	TransactionController *controllers.TransactionController
	rollupSubscriber      dbcache.Subscriber
	sendTxPubSub          dbcache.Subscriber
	eventHandler          *eventhandler.EventHandler
}

func NewServiceContext(context context.Context) *ServiceContext {
	db, err := db.InitializeDB("")
	if err != nil {
		panic(fmt.Sprintf("cannot initialize db, %v\n", err))
	}

	redis, err := dbcache.NewRedisCache(context)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize cache, %v\n", err))
	}

	ethClient, err := clients.InitEthClient()
	if err != nil {
		panic(fmt.Sprintf("cannot initialize eth client, %v\n", err))
	}

	chainId, err := ethClient.ChainID(context)
	if err != nil {
		panic(fmt.Sprintf("cannot get chain id, %v\n", err))
	}

	signer, err := clients.NewSigner(chainId)
	if err != nil {
		panic(fmt.Sprintf("cannot create signer, %v\n", err))
	}

	contractAbi, err := abi.JSON(strings.NewReader(os.Getenv("ROLLUP_ABI")))
	if err != nil {
		panic(fmt.Sprintf("cannot parse abi, %v\n", err))
	}

	accountDao := daos.NewAccountDao(db)
	accountService := services.NewAccountService(&accountDao)
	accountController := controllers.NewAccountController(accountService)
	accountTree, err := tree.InitAccountTree()
	if err != nil {
		panic(fmt.Sprintf("cannot create merkletree, %v\n", err))
	}

	eventHandler, err := eventhandler.NewEventHandler(context, accountService, accountTree, ethClient, &contractAbi)
	if err != nil {
		panic(fmt.Sprintf("cannot create event handler, %v\n", err))
	}

	transctionService := services.NewTransactionService(accountService, accountTree, redis, ethClient, signer, &contractAbi, context)

	sendTxSubscriber := dbcache.NewSendTxPubSub(context, redis, ethClient, "sendTxChannel")
	rollupSubscriber := dbcache.NewRollupSubscriber(redis, signer, ethClient, &contractAbi, "rollupChannel", context)

	tracsactionController := controllers.NewTransactionController(transctionService, &sendTxSubscriber)

	return &ServiceContext{
		PostgresDB:            db,
		Redis:                 redis,
		AccountTree:           accountTree,
		EthClient:             ethClient,
		AccountService:        accountService,
		AccountController:     accountController,
		TransactionController: tracsactionController,
		rollupSubscriber:      rollupSubscriber,
		sendTxPubSub:          sendTxSubscriber,
		eventHandler:          eventHandler,
	}
}

func (svc *ServiceContext) StartDaemon() {
	svc.eventHandler.Listening()
	svc.rollupSubscriber.Receive()
	svc.sendTxPubSub.Receive()
}
