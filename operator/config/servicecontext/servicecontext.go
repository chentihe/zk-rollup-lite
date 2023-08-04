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
	"github.com/chentihe/zk-rollup-lite/operator/layer1/contracts"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/eventhandler"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/chentihe/zk-rollup-lite/operator/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type ServiceContext struct {
	PostgresDB            *gorm.DB
	Redis                 *cache.RedisCache
	AccountTree           *tree.AccountTree
	EthClient             *ethclient.Client
	RollUpAddress         *common.Address
	Abi                   *abi.ABI
	AccountService        *services.AccountService
	AccountController     *controllers.AccountController
	TransactionController *controllers.TransactionController
	ContractController    *controllers.ContractController
	Deployer              *contracts.Deployer
	txManager             *txmanager.TxManager
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

	// ethClient for the interaction with EVM use
	// wsClient for event handler use
	ethClient, wsClient, err := clients.InitEthClient(&config.EthClient)
	if err != nil {
		panic(fmt.Sprintf("cannot initialize eth clients, %v\n", err))
	}

	// this signer is to rollup layer2 tx and deploy contracts
	signer, err := clients.NewSigner(context, config.EthClient.PrivateKey, ethClient)
	if err != nil {
		panic(fmt.Sprintf("cannot create signer, %v\n", err))
	}

	// deploy tx verifier, withdraw verifier, rollup contracts
	deployer := contracts.NewDeployer(ethClient, signer, &config.SmartContracts)

	contractAddress := common.HexToAddress(config.SmartContracts.Rollup.Address)
	contractAbi, err := abi.JSON(strings.NewReader(config.SmartContracts.Rollup.Abi))
	if err != nil {
		panic(fmt.Sprintf("cannot parse abi, %v\n", err))
	}

	accountDao := daos.NewAccountDao(db)
	if err = accountDao.CreateAccountTable(); err != nil {
		panic(fmt.Sprintf("cannot create account table, %v\n", err))
	}
	accountService := services.NewAccountService(&accountDao)
	accountController := controllers.NewAccountController(accountService)
	accountTree, err := tree.InitAccountTree(context, ethClient, &contractAbi, &contractAddress, &config.Postgres)
	if err != nil {
		panic(fmt.Sprintf("cannot create merkletree, %v\n", err))
	}

	eventHandler, err := eventhandler.NewEventHandler(context, accountService, accountTree, wsClient, &contractAbi, &contractAddress, config.Accounts)
	if err != nil {
		panic(fmt.Sprintf("cannot create event handler, %v\n", err))
	}

	transctionService := services.NewTransactionService(context, accountService, accountTree, redis, &config.Redis.Keys)
	transactionController := controllers.NewTransactionController(transctionService)
	txManager := txmanager.NewTxManager(context, redis, signer, ethClient, &contractAbi, &contractAddress, config.Circuit.Path, &config.Redis)

	contractService, err := services.NewContractService(context, ethClient, &contractAddress)
	if err != nil {
		panic(fmt.Sprintf("cannot create contract service, %v\n", err))
	}
	contractController := controllers.NewContractController(contractService)

	return &ServiceContext{
		PostgresDB:            db,
		Redis:                 redis,
		AccountTree:           accountTree,
		EthClient:             ethClient,
		RollUpAddress:         &contractAddress,
		Abi:                   &contractAbi,
		AccountService:        accountService,
		AccountController:     accountController,
		TransactionController: transactionController,
		ContractController:    contractController,
		Deployer:              deployer,
		txManager:             txManager,
		eventHandler:          eventHandler,
	}
}

func (svc *ServiceContext) StartDaemon() {
	svc.eventHandler.Listening()
	svc.txManager.Listening()
}
