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
	Signer     types.Signer
	Context    context.Context
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
		Signer:     types.NewLondonSigner(chainId),
		Context:    context,
	}, nil
}

func (signer *Signer) GetAuth() (*bind.TransactOpts, error) {
	gasPrice, err := signer.ethClient.SuggestGasPrice(signer.Context)
	if err != nil {
		return nil, err
	}

	chainId, err := signer.ethClient.ChainID(signer.Context)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(signer.privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Value = big.NewInt(0)
	// auth.GasLimit = uint64(800000)
	auth.GasPrice = gasPrice

	return auth, nil
}

func (signer *Signer) GenerateDynamicTx(to *common.Address, data []byte, value *big.Int) (*types.Transaction, error) {
	chainId, err := signer.ethClient.ChainID(signer.Context)
	if err != nil {
		return nil, err
	}

	nonce, err := signer.ethClient.PendingNonceAt(signer.Context, signer.Address)
	if err != nil {
		return nil, err
	}

	tipCap, err := signer.ethClient.SuggestGasTipCap(signer.Context)
	if err != nil {
		return nil, err
	}
	feeCap, err := signer.ethClient.SuggestGasPrice(signer.Context)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:  signer.Address,
		To:    to,
		Value: value,
		Data:  data,
	}

	gasLimit, err := signer.ethClient.EstimateGas(signer.Context, callMsg)
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainId,
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
	return types.SignTx(tx, signer.Signer, signer.privateKey)
}
