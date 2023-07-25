package clients

import (
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthClient(config *config.EthClient) (*ethclient.Client, *ethclient.Client, error) {
	rpcClient, err := ethclient.Dial(config.RPCUrl)
	if err != nil {
		return nil, nil, err
	}

	wsClient, err := ethclient.Dial(config.WSUrl)
	if err != nil {
		return nil, nil, err
	}

	return rpcClient, wsClient, nil
}
