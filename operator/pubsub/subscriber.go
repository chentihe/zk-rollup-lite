package pubsub

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/chentihe/zk-rollup-lite/operator/contracts"
	"github.com/chentihe/zk-rollup-lite/operator/dbcache"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/zeroknowledge"
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
	channel                = "pendingTx"
	rollUpCommand          = "execute roll up"
	zkeyFilePath           = "../prover/build/tx/circuit_final.zkey"
	wasmFilePath           = "../prover/build/tx/circuit.wasm"
	verficationKeyFilePath = "../prover/build/tx/verification_key.json"
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

func (sub *Subscriber) Receive(context context.Context, redisCache *dbcache.RedisCache, accountService *services.AccountService) error {
	ch := sub.pubsub.Channel()

	for msg := range ch {
		switch msg.String() {
		case rollUpCommand:

			// get signer
			privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
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

			// init rollup contract
			address := common.HexToAddress(rollupAddress)
			instance, err := contracts.NewRollup(address, sub.ethclient)
			if err != nil {
				return err
			}

			// get tx amounts from redis
			const lastInsertedKey = "last-inserted"
			lastInsertedTx, err := redisCache.Get(context, lastInsertedKey, new(int))
			if err != nil {
				return err
			}

			circuitInput, err := zeroknowledge.GenerateCircuitInput(lastInsertedTx.(int), redisCache, accountService)
			if err != nil {
				return err
			}

			proof, err := zeroknowledge.GenerateGroth16Proof(circuitInput)
			if err != nil {
				return err
			}

			if err := zeroknowledge.VerifierGroth16(proof); err != nil {
				return err
			}

			bigIntProof, err := zeroknowledge.ParseProofToBigInt(proof)
			if err != nil {
				return err
			}

			tx, err := instance.RollUp(auth, bigIntProof.Proof.A, bigIntProof.Proof.B, bigIntProof.Proof.C, bigIntProof.PublicSignals)
			if err != nil {
				return err
			}

			if err := redisCache.Set(context, lastInsertedKey, -1); err != nil {
				return err
			}

			fmt.Printf("Rollup finished: %v", tx)
		default:
			fmt.Printf("Error msg: %v", msg.String())
		}
	}

	return nil
}
