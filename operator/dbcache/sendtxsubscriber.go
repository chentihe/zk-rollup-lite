package dbcache

import (
	"context"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type SendTxSubscriber struct {
	redisCache *RedisCache
	ethClient  *ethclient.Client
	channel    string
	context    context.Context
}

func NewSendTxPubSub(context context.Context, redisCache *RedisCache, ethClient *ethclient.Client, channel string) Subscriber {
	return &SendTxSubscriber{
		redisCache: redisCache,
		ethClient:  ethClient,
		channel:    channel,
	}
}

func (sub *SendTxSubscriber) Publish(msg interface{}) {
	sub.redisCache.client.Publish(sub.context, sub.channel, msg)
}

func (sub *SendTxSubscriber) Receive() {
	pubsub := sub.redisCache.client.Subscribe(sub.context, sub.channel)
	ch := pubsub.Channel()

	var forever chan struct{}

	go func() {
		for msg := range ch {

			txHash := msg.String()
			log.Printf("Received a new signed tx: %s", txHash)

			rawTxBytes, err := hex.DecodeString(txHash)
			if err != nil {
				log.Printf("Decode tx hash error: %s", err)
			}

			tx := new(types.Transaction)
			rlp.DecodeBytes(rawTxBytes, &tx)
			err = sub.ethClient.SendTransaction(sub.context, tx)
			if err != nil {
				log.Printf("Send tx error: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
