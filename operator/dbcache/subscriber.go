package dbcache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/contracts"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	redisCache *RedisCache
	pubsub     *redis.PubSub
	signer     *clients.Signer
	ethclient  *ethclient.Client
	contract   *contracts.Rollup
}

func NewSubscriber(redisCache *RedisCache, signer *clients.Signer, ethclient *ethclient.Client, contract *contracts.Rollup) *Subscriber {
	pubsub := redisCache.Subscribe(context.Background(), channel)

	return &Subscriber{
		redisCache: redisCache,
		pubsub:     pubsub,
		signer:     signer,
		ethclient:  ethclient,
		contract:   contract,
	}
}

func (sub *Subscriber) Close() error {
	return sub.pubsub.Close()
}

func (sub *Subscriber) Receive(context context.Context) {
	ch := sub.pubsub.Channel()

	for msg := range ch {
		switch msg.String() {
		case rollUpCommand:
			// get tx amounts from redis
			lastInsertedTx, err := sub.redisCache.Get(context, lastInsertedKey, new(int))
			if err != nil {
				fmt.Printf("Get tx num err: %v", err)
			}

			var rollupInputs circuits.RollupInputs
			for i := 0; i < lastInsertedTx.(int); i++ {
				object, err := sub.redisCache.Get(context, strconv.Itoa(i), new(circuits.RollupTx))
				if err != nil {
					fmt.Printf("Get rollup tx err: %v", err)
				}

				tx, ok := object.(circuits.RollupTx)
				if !ok {
					fmt.Printf("Casting err: %v", circuits.ErrTx)
				}
				rollupInputs.Txs = append(rollupInputs.Txs, &tx)
			}

			circuitInput, err := rollupInputs.InputsMarshal()
			if err != nil {
				fmt.Printf("Circuit inputs marshal err: %v", err)
			}

			proof, err := circuits.GenerateGroth16Proof(circuitInput, circuitPath+"/tx")
			if err != nil {
				fmt.Printf("Generate proof err: %v", err)
			}

			if err = circuits.VerifierGroth16(proof, circuitPath+"/tx"); err != nil {
				fmt.Printf("Verify proof err: %v", err)
			}

			var rollupOutputs circuits.RollupOutputs
			if err = rollupOutputs.OutputUnmarshal(proof); err != nil {
				fmt.Printf("Circuit ouputs unmarshal err: %v", err)
			}

			auth, err := sub.signer.GetAuth(sub.ethclient, context)
			if err != nil {
				fmt.Printf("Get signer auth err: %v", err)
			}

			tx, err := sub.contract.RollUp(auth, rollupOutputs.Proof.A, rollupOutputs.Proof.B, rollupOutputs.Proof.C, rollupOutputs.PublicSignals)
			if err != nil {
				fmt.Printf("Send tx err: %v", err)
			}

			if err := sub.redisCache.Set(context, lastInsertedKey, -1); err != nil {
				fmt.Printf("Update redis err: %v", err)
			}

			fmt.Printf("Rollup finished: %v", tx)
		default:
			fmt.Printf("Error msg: %v", msg.String())
		}
	}
}
