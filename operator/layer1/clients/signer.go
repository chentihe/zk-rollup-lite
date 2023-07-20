package clients

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"

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
	Signer     types.Signer
	Context    context.Context
}

func NewSigner(chainId *big.Int) (*Signer, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrPubKeyToECDSA
	}

	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Signer{
		privateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    signerAddress,
		Signer:     types.NewEIP155Signer(chainId),
		Context:    context.Background(),
	}, nil
}

func (signer *Signer) GetAuth(ethClient *ethclient.Client) (*bind.TransactOpts, error) {
	nonce, err := ethClient.PendingNonceAt(signer.Context, signer.Address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := ethClient.SuggestGasPrice(signer.Context)
	if err != nil {
		return nil, err
	}

	chainId, err := ethClient.ChainID(signer.Context)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(signer.privateKey, chainId)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	return auth, nil
}

func (signer *Signer) GenerateDynamicTx(ethClient *ethclient.Client, to common.Address, data []byte) (*types.Transaction, error) {
	chainId, err := ethClient.ChainID(signer.Context)
	if err != nil {
		return nil, err
	}

	nonce, err := ethClient.PendingNonceAt(signer.Context, signer.Address)
	if err != nil {
		return nil, err
	}

	tipCap, err := ethClient.SuggestGasTipCap(signer.Context)
	if err != nil {
		return nil, err
	}

	feeCap, err := ethClient.SuggestGasPrice(signer.Context)
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:      signer.Address,
		To:        &to,
		GasPrice:  nil,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Value:     big.NewInt(0),
		Data:      data,
	}

	gasLimit, err := ethClient.EstimateGas(signer.Context, callMsg)
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
			To:        &signer.Address,
			Value:     big.NewInt(0),
			Data:      data,
		})

	return tx, nil
}

func (signer *Signer) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	return types.SignTx(tx, signer.Signer, signer.privateKey)
}
