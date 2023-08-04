package contracts

import (
	"context"
	"math/big"
	"strings"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type Deployer struct {
	ethClient      *ethclient.Client
	signer         *clients.Signer
	smartContracts *config.SmartContracts
}

func NewDeployer(ethClient *ethclient.Client, signer *clients.Signer, smartContracts *config.SmartContracts) *Deployer {
	return &Deployer{
		ethClient:      ethClient,
		signer:         signer,
		smartContracts: smartContracts,
	}
}

func (d *Deployer) Deploy() error {
	context := context.Background()

	txVerifier := d.smartContracts.TxVerifier

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
	if err = d.ethClient.SendTransaction(context, verifierTx); err != nil {
		return err
	}
	receipt, err := d.ethClient.TransactionReceipt(context, verifierTx.Hash())
	if err != nil {
		return err
	}

	// txVerifierAddress := crypto.CreateAddress(d.signer.Address, verifierTx.Nonce())
	d.smartContracts.TxVerifier.Address = receipt.ContractAddress.String()

	withdrawVerifier := d.smartContracts.WithdrawVerifier

	contractAbi, err = abi.JSON(strings.NewReader(withdrawVerifier.Abi))
	if err != nil {
		return err
	}

	data, err = contractAbi.Pack("")
	if err != nil {
		return err
	}

	tx, err = d.signer.GenerateLegacyTx(nil, append(common.FromHex(withdrawVerifier.ByteCode), data...), big.NewInt(0))
	if err != nil {
		return err
	}

	withdrawTx, err := d.signer.SignTx(tx)
	if err != nil {
		return err
	}
	if err = d.ethClient.SendTransaction(context, withdrawTx); err != nil {
		return err
	}
	receipt, err = d.ethClient.TransactionReceipt(context, withdrawTx.Hash())
	if err != nil {
		return err
	}

	// withdrawVerifierAddress := crypto.CreateAddress(d.signer.Address, withdrawTx.Nonce())
	d.smartContracts.WithdrawVerifier.Address = receipt.ContractAddress.String()

	rollup := d.smartContracts.Rollup

	contractAbi, err = abi.JSON(strings.NewReader(rollup.Abi))
	if err != nil {
		return err
	}

	// deploy contract uses contructor, just passing the constructor args
	data, err = contractAbi.Pack("", d.smartContracts.TxVerifier.Address, d.smartContracts.WithdrawVerifier.Address)
	if err != nil {
		return err
	}

	// call data consist of the contract bytecode and constructor args
	tx, err = d.signer.GenerateLegacyTx(nil, append(common.FromHex(rollup.ByteCode), data...), big.NewInt(0))
	if err != nil {
		return err
	}

	rollupTx, err := d.signer.SignTx(tx)
	if err != nil {
		return err
	}
	if err = d.ethClient.SendTransaction(context, rollupTx); err != nil {
		return err
	}
	receipt, err = d.ethClient.TransactionReceipt(context, rollupTx.Hash())
	if err != nil {
		return err
	}

	// the result of the contract address is the same between receipt and crypto.CreateAddress
	// rollupAddress := crypto.CreateAddress(d.signer.Address, rollupTx.Nonce())
	d.smartContracts.Rollup.Address = receipt.ContractAddress.String()

	// update the env.yaml for the deposit / withdraw cli
	viper.Set("smartcontracts", d.smartContracts)
	viper.WriteConfig()
	return nil
}