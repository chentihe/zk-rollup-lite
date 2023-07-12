package pubsub

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/contracts"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	pubsub    *redis.PubSub
	ethclient *ethclient.Client
}

const (
	channel       = "pendingTx"
	rollUpCommand = "execute roll up"
)

func NewSubscriber(redisCache *dbcache.RedisCache, ethclient *ethclient.Client) *Subscriber {
	pubsub := redisCache.Subscribe(context.Background(), channel)

	return &Subscriber{
		pubsub:    pubsub,
		ethclient: ethclient,
	}
}

func (sub *Subscriber) Close() error {
	return sub.pubsub.Close()
}

const rollupAddress = ""

func (sub *Subscriber) Receive(context context.Context, prvKey string) error {
	ch := sub.pubsub.Channel()

	for msg := range ch {
		switch msg.String() {
		case rollUpCommand:

			privateKey, err := crypto.HexToECDSA(prvKey)
			if err != nil {
				return err
			}

			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return ErrPubKeyToECDSA
			}

			fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

			nonce, err := sub.ethclient.PendingNonceAt(context, fromAddress)
			if err != nil {
				return err
			}

			gasPrice, err := sub.ethclient.SuggestGasPrice(context)
			if err != nil {
				return err
			}

			chainId, err := sub.ethclient.ChainID(context)
			if err != nil {
				return err
			}

			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
			if err != nil {
				return err
			}

			auth.Nonce = big.NewInt(int64(nonce))
			auth.Value = big.NewInt(0)
			auth.GasLimit = uint64(300000)
			auth.GasPrice = gasPrice

			address := common.HexToAddress(rollupAddress)
			instance, err := contracts.NewRollup(address, sub.ethclient)
			if err != nil {
				return err
			}

			// TODO: Get txs from redis and generate zkp(iden3 go lib)
			tx, err := instance.RollUp(auth, a, b, c, input)
			if err != nil {
				return err
			}

		default:
			fmt.Printf("Error msg: %v", msg.String())
		}
	}

	return nil
}
