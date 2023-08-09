package contracts

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type Deployer struct {
	ethClient      *ethclient.Client
	signer         *clients.Signer
	smartContracts *config.SmartContracts
	context        context.Context
}

func NewDeployer(context context.Context, ethClient *ethclient.Client, signer *clients.Signer, smartContracts *config.SmartContracts) *Deployer {
	return &Deployer{
		ethClient:      ethClient,
		signer:         signer,
		smartContracts: smartContracts,
		context:        context,
	}
}

func (d *Deployer) Deploy() error {
	var (
		txVerifierAddress       = common.HexToAddress(d.smartContracts.TxVerifier.Address)
		withdrawVerifierAddress = common.HexToAddress(d.smartContracts.WithdrawVerifier.Address)
	)

	txVerifier := d.smartContracts.TxVerifier
	if !d.isDeployed(txVerifier.Address) {
		contractAbi, err := abi.JSON(strings.NewReader(txVerifier.Abi))
		if err != nil {
			return err
		}

		data, err := contractAbi.Pack("")
		if err != nil {
			return err
		}

		tx, err := d.signer.GenerateLegacyTx(nil, append(common.FromHex(txVerifier.ByteCode), data...), big.NewInt(0))
		if err != nil {
			return err
		}

		verifierTx, err := d.signer.SignTx(tx)
		if err != nil {
			return err
		}
		if err = d.ethClient.SendTransaction(d.context, verifierTx); err != nil {
			return err
		}
		d.ethClient.TransactionReceipt(d.context, verifierTx.Hash())

		txVerifierAddress = crypto.CreateAddress(d.signer.Address, verifierTx.Nonce())
		d.smartContracts.TxVerifier.Address = txVerifierAddress.String()
		log.Printf("Deployed tx verifier: %s\n", txVerifierAddress.String())
	}

	withdrawVerifier := d.smartContracts.WithdrawVerifier
	if !d.isDeployed(withdrawVerifier.Address) {
		contractAbi, err := abi.JSON(strings.NewReader(withdrawVerifier.Abi))
		if err != nil {
			return err
		}

		data, err := contractAbi.Pack("")
		if err != nil {
			return err
		}

		tx, err := d.signer.GenerateLegacyTx(nil, append(common.FromHex(withdrawVerifier.ByteCode), data...), big.NewInt(0))
		if err != nil {
			return err
		}

		withdrawTx, err := d.signer.SignTx(tx)
		if err != nil {
			return err
		}
		if err = d.ethClient.SendTransaction(d.context, withdrawTx); err != nil {
			return err
		}
		d.ethClient.TransactionReceipt(d.context, withdrawTx.Hash())

		withdrawVerifierAddress = crypto.CreateAddress(d.signer.Address, withdrawTx.Nonce())
		d.smartContracts.WithdrawVerifier.Address = withdrawVerifierAddress.String()
		log.Printf("Deployed withdraw verifier: %s\n", withdrawVerifierAddress.String())
	}

	rollup := d.smartContracts.Rollup
	if !d.isDeployed(rollup.Address) {
		contractAbi, err := abi.JSON(strings.NewReader(rollup.Abi))
		if err != nil {
			return err
		}

		// deploy contract uses contructor, just passing the constructor args
		data, err := contractAbi.Pack("", txVerifierAddress, withdrawVerifierAddress)
		if err != nil {
			return err
		}

		// call data consist of the contract bytecode and constructor args
		tx, err := d.signer.GenerateLegacyTx(nil, append(common.FromHex(rollup.ByteCode), data...), big.NewInt(0))
		if err != nil {
			return err
		}

		rollupTx, err := d.signer.SignTx(tx)
		if err != nil {
			return err
		}
		if err = d.ethClient.SendTransaction(d.context, rollupTx); err != nil {
			return err
		}
		d.ethClient.TransactionReceipt(d.context, rollupTx.Hash())

		// the result of the contract address is the same between receipt and crypto.CreateAddress
		rollupAddress := crypto.CreateAddress(d.signer.Address, rollupTx.Nonce())
		d.smartContracts.Rollup.Address = rollupAddress.String()
		log.Printf("Deployed rollup: %s\n", rollupAddress.String())
		// update the env.yaml for the deposit / withdraw cli
		viper.Set("smartcontracts", d.smartContracts)
		viper.WriteConfig()
	}
	return nil
}

func (d *Deployer) isDeployed(hexAddress string) bool {
	address := common.HexToAddress(hexAddress)
	bytecode, err := d.ethClient.CodeAt(d.context, address, nil)
	if err != nil {
		return false
	}

	return len(bytecode) > 0
}
