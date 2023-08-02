package txmanger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/chentihe/zk-rollup-lite/operator/cache"
	"github.com/chentihe/zk-rollup-lite/operator/circuits"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxManager struct {
	redisCache      *cache.RedisCache
	signer          *clients.Signer
	ethclient       *ethclient.Client
	contractAddress *common.Address
	abi             *abi.ABI
	context         context.Context
	circuitPath     string
	keys            *config.Keys
	commands        *config.Commands
}

const batchSize = 2

func NewTxManager(context context.Context, redisCache *cache.RedisCache, signer *clients.Signer, ethclient *ethclient.Client, abi *abi.ABI, address *common.Address, circuitPath string, config *config.Redis) *TxManager {
	return &TxManager{
		redisCache:      redisCache,
		signer:          signer,
		ethclient:       ethclient,
		contractAddress: address,
		abi:             abi,
		context:         context,
		circuitPath:     circuitPath + "/tx",
		keys:            &config.Keys,
		commands:        &config.Commands,
	}
}

func (txManager *TxManager) Listening() {
	fmt.Println("Listening to txs...")

	// tx manager will roll up txs every 10s
	timer := time.NewTicker(time.Second * 10)
	go func() {
		for {
			<-timer.C
			// 1.1 get finished tx amounts from redis
			value, err := txManager.redisCache.Get(txManager.context, txManager.keys.LastInsertedKey)
			if err != nil {
				log.Printf("Get last inserted tx num err: %v\n", err)
			}

			lastInsertedTx, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("Get last inserted tx num err: %v\n", err)
			}

			// 1.2 get rolluped tx amounts from redis
			value, err = txManager.redisCache.Get(txManager.context, txManager.keys.RollupedTxsKey)
			if err != nil {
				log.Printf("Get rolluped tx num err: %v\n", err)
			}

			rollupedTxs, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("Get rolluped tx num err: %v\n", err)
			}

			// do rollup if the pending txs are more than batch size
			// rollupedTxs starts from 0
			// 1st round rollup no.1 & no.2 txs
			// next round starts from no.3
			// if last inserted tx is odd or no more txs need to roll up
			// break the loop
			if lastInsertedTx-rollupedTxs >= batchSize {
				var rolluped []string
				for i := rollupedTxs; i <= lastInsertedTx; {
					txManager.Rollup(&rolluped, i+1, i+batchSize)
					i += batchSize
					if lastInsertedTx <= i || lastInsertedTx-i == 1 {
						break
					}
				}
				// delete the rolluped txs from redis
				if err = txManager.redisCache.Del(txManager.context, rolluped); err != nil {
					log.Printf("Del redis keys err: %v\n", err)
				}
			}
		}
	}()
}

func (txManager *TxManager) Rollup(rolluped *[]string, start int, end int) {
	var rollupInputs circuits.RollupInputs
	for i := start; i <= end; i++ {
		var tx circuits.RollupTx
		object, err := txManager.redisCache.Get(txManager.context, strconv.Itoa(i))
		if err != nil {
			log.Printf("Get rollup tx err: %v", err)
		}

		json.Unmarshal([]byte(object), &tx)
		rollupInputs.Txs = append(rollupInputs.Txs, &tx)
		*rolluped = append(*rolluped, strconv.Itoa(i))
	}

	circuitInput, err := rollupInputs.InputsMarshal()
	if err != nil {
		log.Printf("Circuit inputs marshal err: %v", err)
	}

	proof, err := circuits.GenerateGroth16Proof(circuitInput, txManager.circuitPath)
	if err != nil {
		log.Printf("Generate proof err: %v", err)
	}

	if err = circuits.VerifierGroth16(proof, txManager.circuitPath); err != nil {
		log.Printf("Verify proof err: %v", err)
	}

	var rollupOutputs circuits.RollupOutputs
	if err = rollupOutputs.OutputUnmarshal(proof); err != nil {
		log.Printf("Circuit ouputs unmarshal err: %v", err)
	}

	data, err := txManager.abi.Pack("rollUp", rollupOutputs.Proof.A, rollupOutputs.Proof.B, rollupOutputs.Proof.C, rollupOutputs.PublicSignals)
	if err != nil {
		log.Printf("Cannot pack rollup call data: %v\n", err)
	}

	tx, err := txManager.signer.GenerateLegacyTx(txManager.contractAddress, data, big.NewInt(0))
	if err != nil {
		log.Printf("Send tx err: %v\n", err)
	}

	signTx, err := txManager.signer.SignTx(tx)
	if err != nil {
		log.Printf("Sign tx err: %v\n", err)
	}

	if err = txManager.ethclient.SendTransaction(txManager.context, signTx); err != nil {
		log.Printf("Send tx err: %v\n", err)
	}

	if err = txManager.redisCache.Set(txManager.context, txManager.keys.RollupedTxsKey, strconv.Itoa(end)); err != nil {
		log.Printf("Update redis err: %v\n", err)
	}

	log.Printf("Rollup success: %v\n", signTx.Hash().String())
}
