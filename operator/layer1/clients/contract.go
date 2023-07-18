package clients

import (
	"github.com/chentihe/zk-rollup-lite/operator/layer1/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewRollUp(ethClient *ethclient.Client) (*contracts.Rollup, error) {
	address := common.HexToAddress(rollupAddress)
	instance, err := contracts.NewRollup(address, ethClient)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
