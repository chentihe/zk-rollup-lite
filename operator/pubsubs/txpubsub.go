package pubsubs

import (
	"context"
	"log"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/layer1"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxPubSub struct {
	redisCache *cache.RedisCache
	ethClient  *ethclient.Client
	channel    string
	context    context.Context
}

func NewTxPubSub(context context.Context, redisCache *cache.RedisCache, ethClient *ethclient.Client, channel string) Subscriber {
	return &TxPubSub{
		redisCache: redisCache,
		ethClient:  ethClient,
		channel:    channel,
	}
}

func (pubsub *TxPubSub) Publish(msg interface{}) {
	pubsub.redisCache.Client.Publish(pubsub.context, pubsub.channel, msg)
}

func (pubsub *TxPubSub) Receive() {
	sub := pubsub.redisCache.Client.Subscribe(pubsub.context, pubsub.channel)
	ch := sub.Channel()

	var forever chan struct{}

	go func() {
		for msg := range ch {

			txHash := msg.String()
			log.Printf("Received a new signed tx: %s", txHash)

			tx, err := layer1.DecodeTxHash(txHash)
			if err != nil {
				log.Printf("Decode tx error: %s", err)
			}

			err = pubsub.ethClient.SendTransaction(pubsub.context, tx)
			if err != nil {
				log.Printf("Send tx error: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
