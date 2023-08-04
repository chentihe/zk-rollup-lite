package clients

import (
	"fmt"
	"os"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthClient(config *config.EthClient) (*ethclient.Client, *ethclient.Client, error) {
	// for docker image to retrieve anvil host
	url := os.Getenv("ANVIL_URL")
	if url != "" {
		config.RPCUrl = fmt.Sprintf("http://%s:8545", url)
		config.WSUrl = fmt.Sprintf("ws://%s:8545", url)
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
