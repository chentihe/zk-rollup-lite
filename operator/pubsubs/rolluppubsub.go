package pubsubs

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RollupPubSub struct {
	redisCache      *cache.RedisCache
	signer          *clients.Signer
	ethclient       *ethclient.Client
	contractAddress *common.Address
	abi             *abi.ABI
	channel         string
	context         context.Context
	circuitPath     string
	keys            *config.Keys
	commands        *config.Commands
}

func NewRollupPubSub(redisCache *cache.RedisCache, signer *clients.Signer, ethclient *ethclient.Client, abi *abi.ABI, channel string, context context.Context, contractAddress string, circuitPath string, config *config.Redis) Subscriber {
	rollupContract := common.HexToAddress(contractAddress)
	return &RollupPubSub{
		redisCache:      redisCache,
		signer:          signer,
		ethclient:       ethclient,
		contractAddress: &rollupContract,
		abi:             abi,
		channel:         channel,
		context:         context,
		circuitPath:     circuitPath,
		keys:            &config.Keys,
		commands:        &config.Commands,
	}
}

func (pubsub *RollupPubSub) Publish(msg interface{}) {
	pubsub.redisCache.Client.Publish(pubsub.context, pubsub.channel, msg)
}

func (pubsub *RollupPubSub) Receive() {
	sub := pubsub.redisCache.Client.Subscribe(pubsub.context, pubsub.channel)
	ch := sub.Channel()

	go func() {
		for msg := range ch {
			switch msg.String() {
			case pubsub.commands.RollupCommand:
				// get tx amounts from redis
				lastInsertedTx, err := pubsub.redisCache.Get(pubsub.context, pubsub.keys.LastInsertedKey, new(int))
				if err != nil {
					fmt.Printf("Get tx num err: %v", err)
				}

				var rollupInputs circuits.RollupInputs
				for i := 0; i < lastInsertedTx.(int); i++ {
					object, err := pubsub.redisCache.Get(pubsub.context, strconv.Itoa(i), new(circuits.RollupTx))
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

				proof, err := circuits.GenerateGroth16Proof(circuitInput, pubsub.circuitPath+"/tx")
				if err != nil {
					fmt.Printf("Generate proof err: %v", err)
				}

				if err = circuits.VerifierGroth16(proof, pubsub.circuitPath+"/tx"); err != nil {
					fmt.Printf("Verify proof err: %v", err)
				}

				var rollupOutputs circuits.RollupOutputs
				if err = rollupOutputs.OutputUnmarshal(proof); err != nil {
					fmt.Printf("Circuit ouputs unmarshal err: %v", err)
				}

				data, err := pubsub.abi.Pack("rollup", rollupOutputs.Proof.A, rollupOutputs.Proof.B, rollupOutputs.Proof.C, rollupOutputs.PublicSignals)
				if err != nil {
					fmt.Printf("Cannot pack rollup call data: %v", err)
				}

				tx, err := pubsub.signer.GenerateDynamicTx(pubsub.contractAddress, data, big.NewInt(0))
				if err != nil {
					fmt.Printf("Send tx err: %v", err)
				}

				signTx, err := pubsub.signer.SignTx(tx)
				if err != nil {
					fmt.Printf("Sign tx err: %v", err)
				}

				if err = pubsub.ethclient.SendTransaction(pubsub.context, signTx); err != nil {
					fmt.Printf("Send tx err: %v", err)
				}

				if err = pubsub.redisCache.Set(pubsub.context, pubsub.keys.LastInsertedKey, -1); err != nil {
					fmt.Printf("Update redis err: %v", err)
				}

				fmt.Printf("Rollup finished: %v", tx)
			default:
				fmt.Printf("Error msg: %v", msg.String())
			}
		}
	}()
}
