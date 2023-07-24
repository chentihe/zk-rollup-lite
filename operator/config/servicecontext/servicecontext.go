package servicecontext

import (
	"context"
	"fmt"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/controllers"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/db"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/eventhandler"
	"github.com/chentihe/zk-rollup-lite/operator/pubsubs"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type ServiceContext struct {
	PostgresDB            *gorm.DB
	Redis                 *cache.RedisCache
	AccountTree           *tree.AccountTree
	EthClient             *ethclient.Client
	Abi                   *abi.ABI
	AccountService        *services.AccountService
	AccountController     *controllers.AccountController
	TransactionController *controllers.TransactionController
	rollupPubSub          pubsubs.Subscriber
	txPubSub              pubsubs.Subscriber
	eventHandler          *eventhandler.EventHandler
}

func NewServiceContext(context context.Context, config *config.Config) *ServiceContext {
	db, err := db.InitializeDB(&config.Postgres)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize db, %v\n", err))
	}

	redis, err := cache.NewRedisCache(context, &config.Redis)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize cache, %v\n", err))
	}

	ethClient, err := clients.InitEthClient(&config.EthClient)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize eth client, %v\n", err))
	}

	chainId, err := ethClient.ChainID(context)
	if err != nil {
		panic(fmt.Sprintf("cannot get chain id, %v\n", err))
	}

	signer, err := clients.NewSigner(chainId, config.EthClient.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("cannot create signer, %v\n", err))
	}

	contractAbi, err := abi.JSON(strings.NewReader(config.SmartContract.Abi))
	if err != nil {
		panic(fmt.Sprintf("cannot parse abi, %v\n", err))
	}

	accountDao := daos.NewAccountDao(db)
	if err = accountDao.CreateAccountTable(); err != nil {
		panic(fmt.Sprintf("cannot create account table, %v\n", err))
	}
	accountService := services.NewAccountService(&accountDao)
	accountController := controllers.NewAccountController(accountService)
	accountTree, err := tree.InitAccountTree()
	if err != nil {
		panic(fmt.Sprintf("cannot create merkletree, %v\n", err))
	}

	eventHandler, err := eventhandler.NewEventHandler(context, accountService, accountTree, ethClient, &contractAbi, config.SmartContract.Address)
	if err != nil {
		panic(fmt.Sprintf("cannot create event handler, %v\n", err))
	}

	transctionService := services.NewTransactionService(accountService, accountTree, redis, ethClient, signer, &contractAbi, context)

	txPubSub := pubsubs.NewTxPubSub(context, redis, ethClient, "sendTxChannel")
	rollupPubSub := pubsubs.NewRollupPubSub(redis, signer, ethClient, &contractAbi, "rollupChannel", context)

	tracsactionController := controllers.NewTransactionController(transctionService, &txPubSub)

	return &ServiceContext{
		PostgresDB:            db,
		Redis:                 redis,
		AccountTree:           accountTree,
		EthClient:             ethClient,
		Abi:                   &contractAbi,
		AccountService:        accountService,
		AccountController:     accountController,
		TransactionController: tracsactionController,
		rollupPubSub:          rollupPubSub,
		txPubSub:              txPubSub,
		eventHandler:          eventHandler,
	}
}

func (svc *ServiceContext) StartDaemon() {
	svc.eventHandler.Listening()
	svc.rollupPubSub.Receive()
	svc.txPubSub.Receive()
}
