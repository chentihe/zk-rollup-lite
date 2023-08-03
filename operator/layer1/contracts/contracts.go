package contracts

import (
	"context"
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
	d.ethClient.SendTransaction(context, verifierTx)
	d.ethClient.TransactionReceipt(context, verifierTx.Hash())

	txVerifierAddress := crypto.CreateAddress(d.signer.Address, verifierTx.Nonce())
	d.smartContracts.TxVerifier.Address = txVerifierAddress.String()

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
	d.ethClient.SendTransaction(context, withdrawTx)
	d.ethClient.TransactionReceipt(context, withdrawTx.Hash())

	withdrawVerifierAddress := crypto.CreateAddress(d.signer.Address, withdrawTx.Nonce())
	d.smartContracts.WithdrawVerifier.Address = withdrawVerifierAddress.String()

	rollup := d.smartContracts.Rollup

	contractAbi, err = abi.JSON(strings.NewReader(rollup.Abi))
	if err != nil {
		return err
	}

	data, err = contractAbi.Pack("", txVerifierAddress, withdrawVerifierAddress)
	if err != nil {
		return err
	}

	tx, err = d.signer.GenerateLegacyTx(nil, append(common.FromHex(rollup.ByteCode), data...), big.NewInt(0))
	if err != nil {
		return err
	}

	rollupTx, err := d.signer.SignTx(tx)
	if err != nil {
		return err
	}
	d.ethClient.SendTransaction(context, rollupTx)
	d.ethClient.TransactionReceipt(context, rollupTx.Hash())

	rollupAddress := crypto.CreateAddress(d.signer.Address, rollupTx.Nonce())
	d.smartContracts.Rollup.Address = rollupAddress.String()
	viper.Set("smartcontracts", d.smartContracts)
	viper.WriteConfig()
	return nil
}
