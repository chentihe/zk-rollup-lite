package clients

import (
	"os"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthClient(config *config.EthClient) (*ethclient.Client, *ethclient.Client, error) {
	url := os.Getenv("ANVIL_URL")
	if url != "" {
		config.RPCUrl = "http://" + url + ":8545"
		config.WSUrl = "ws://" + url + ":8545"
	}

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
