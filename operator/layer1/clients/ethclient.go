package clients

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func InitEthClient() (*ethclient.Client, error) {
	// TODO: change to yaml config
	client, err := ethclient.Dial(os.Getenv("WEBSOCKER_URL"))

	if err != nil {
		return nil, err
	}

	return client, nil
}
