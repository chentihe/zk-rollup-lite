package clients

import (
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthClient(config *config.EthClient) (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.WSUrl)

	if err != nil {
		return nil, err
	}

	return client, nil
}
