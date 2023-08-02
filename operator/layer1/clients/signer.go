package clients

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Signer struct {
	privateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
	ethClient  *ethclient.Client
	chainId    *big.Int
	context    context.Context
}

func NewSigner(context context.Context, priv string, ethClient *ethclient.Client) (*Signer, error) {
	privateKey, err := crypto.HexToECDSA(priv)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrPubKeyToECDSA
	}

	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainId, err := ethClient.ChainID(context)
	if err != nil {
		return nil, err
	}

	return &Signer{
		privateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    signerAddress,
		ethClient:  ethClient,
		chainId:    chainId,
		context:    context,
	}, nil
}

func (signer *Signer) GetAuth() (*bind.TransactOpts, error) {
	gasPrice, err := signer.ethClient.SuggestGasPrice(signer.context)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(signer.privateKey, signer.chainId)
	if err != nil {
		return nil, err
	}

	auth.Value = big.NewInt(0)
	auth.GasPrice = gasPrice

	return auth, nil
}

func (signer *Signer) GenerateLegacyTx(to *common.Address, data []byte, value *big.Int) (*types.Transaction, error) {
	nonce, err := signer.ethClient.PendingNonceAt(signer.context, signer.Address)
	if err != nil {
		return nil, err
	}

	feeCap, err := signer.ethClient.SuggestGasPrice(signer.context)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From: signer.Address,
		To:   to,
		Data: data,
	}

	if value != nil {
		callMsg.Value = value
	}

	gasLimit, err := signer.ethClient.EstimateGas(signer.context, callMsg)
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: feeCap,
			Gas:      gasLimit,
			To:       to,
			Value:    value,
			Data:     data,
		})

	return tx, nil
}

func (signer *Signer) GenerateDynamicTx(to *common.Address, data []byte, value *big.Int) (*types.Transaction, error) {
	nonce, err := signer.ethClient.PendingNonceAt(signer.context, signer.Address)
	if err != nil {
		return nil, err
	}

	tipCap, err := signer.ethClient.SuggestGasTipCap(signer.context)
	if err != nil {
		return nil, err
	}
	feeCap, err := signer.ethClient.SuggestGasPrice(signer.context)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:  signer.Address,
		To:    to,
		Value: value,
		Data:  data,
	}

	gasLimit, err := signer.ethClient.EstimateGas(signer.context, callMsg)
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   signer.chainId,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        to,
			Value:     value,
			Data:      data,
		})

	return tx, nil
}

func (signer *Signer) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	return types.SignTx(tx, types.NewLondonSigner(signer.chainId), signer.privateKey)
}
