package dbcache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RollupSubscriber struct {
	redisCache *RedisCache
	signer     *clients.Signer
	ethclient  *ethclient.Client
	abi        *abi.ABI
	channel    string
	context    context.Context
}

func NewRollupSubscriber(redisCache *RedisCache, signer *clients.Signer, ethclient *ethclient.Client, abi *abi.ABI, channel string, context context.Context) Subscriber {
	return &RollupSubscriber{
		redisCache: redisCache,
		signer:     signer,
		ethclient:  ethclient,
		abi:        abi,
		channel:    channel,
		context:    context,
	}
}

func (sub *RollupSubscriber) Publish(msg interface{}) {
	sub.redisCache.client.Publish(sub.context, sub.channel, msg)
}

func (sub *RollupSubscriber) Receive() {
	pubsub := sub.redisCache.client.Subscribe(context.Background(), sub.channel)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			switch msg.String() {
			case rollUpCommand:
				// get tx amounts from redis
				lastInsertedTx, err := sub.redisCache.Get(sub.context, lastInsertedKey, new(int))
				if err != nil {
					fmt.Printf("Get tx num err: %v", err)
				}

				var rollupInputs circuits.RollupInputs
				for i := 0; i < lastInsertedTx.(int); i++ {
					object, err := sub.redisCache.Get(sub.context, strconv.Itoa(i), new(circuits.RollupTx))
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

				data, err := sub.abi.Pack("rollup", rollupOutputs.Proof.A, rollupOutputs.Proof.B, rollupOutputs.Proof.C, rollupOutputs.PublicSignals)
				if err != nil {
					fmt.Printf("Cannot pack rollup call data: %v", err)
				}

				tx, err := sub.signer.GenerateDynamicTx(sub.ethclient, common.HexToAddress(rollupAddress), data)
				if err != nil {
					fmt.Printf("Send tx err: %v", err)
				}

				signTx, err := sub.signer.SignTx(tx)
				if err != nil {
					fmt.Printf("Sign tx err: %v", err)
				}

				if err := sub.ethclient.SendTransaction(sub.context, signTx); err != nil {
					fmt.Printf("Send tx err: %v", err)
				}

				if err := sub.redisCache.Set(sub.context, lastInsertedKey, -1); err != nil {
					fmt.Printf("Update redis err: %v", err)
				}

				fmt.Printf("Rollup finished: %v", tx)
			default:
				fmt.Printf("Error msg: %v", msg.String())
			}
		}
	}()
}
