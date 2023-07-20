package config

import (
	"context"
	"fmt"

	"github.com/chentihe/zk-rollup-lite/operator/controllers"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/db"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type ServiceContext struct {
	PostgresDB            *gorm.DB
	Redis                 *dbcache.RedisCache
	AccountTree           *tree.AccountTree
	Subscriber            *dbcache.Subscriber
	EthClient             *ethclient.Client
	AccountService        *services.AccountService
	AccountController     *controllers.AccountController
	TransactionController *controllers.TransactionController
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

	if err != nil {
		panic(fmt.Sprintf("cannot creat redis cache, %v\n", err))
	}

	ethClient, err := clients.InitEthClient()
	if err != nil {
		panic(fmt.Sprintf("cannot initialize eth client, %v\n", err))
	}

	signer, err := clients.NewSigner()
	if err != nil {
		panic(fmt.Sprintf("cannot create signer, %v\n", err))
	}

	contract, err := clients.NewRollUp(ethClient)
	if err != nil {
		panic(fmt.Sprintf("cannot create contract, %v\n", err))
	}

	accountDao := daos.NewAccountDao(db)
	accountService := services.NewAccountService(&accountDao)
	accountController := controllers.NewAccountController(accountService)

	accountTree, err := tree.InitAccountTree()
	if err != nil {
		panic(fmt.Sprintf("cannot create merkletree, %v\n", err))
	}

	subscriber := dbcache.NewSubscriber(redis, signer, ethClient, contract)
	// TODO: use go routine to subscribe the topic
	subscriber.Receive(context)

	// TODO: account service should be removed from tx service
	transctionService := services.NewTransactionService(accountService, accountTree, redis, ethClient, signer, contract, context)
	tracsactionController := controllers.NewTransactionController(transctionService)

	return &ServiceContext{
		PostgresDB:            db,
		Redis:                 redis,
		AccountTree:           accountTree,
		Subscriber:            subscriber,
		EthClient:             ethClient,
		AccountService:        accountService,
		AccountController:     accountController,
		TransactionController: tracsactionController,
	}
}
