package services

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/layer1/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractService struct {
	ethClient *ethclient.Client
	contract  *contracts.Rollup
	address   *common.Address
	context   context.Context
}

func NewContractService(context context.Context, ethClient *ethclient.Client, address *common.Address) (*ContractService, error) {
	rollup, err := contracts.NewRollup(*address, ethClient)
	if err != nil {
		return nil, err
	}

	return &ContractService{
		ethClient: ethClient,
		contract:  rollup,
		address:   address,
		context:   context,
	}, nil
}

func (s *ContractService) GetUserByIndex(index *big.Int) (contracts.RollupUser, error) {
	return s.contract.GetUserByIndex(nil, index)
}

func (s *ContractService) GetContractBalance() (*big.Int, error) {
	return s.ethClient.BalanceAt(s.context, *s.address, nil)
}

func (s *ContractService) GetStateRoot() (*big.Int, error) {
	return s.contract.BalanceTreeRoot(nil)
}
