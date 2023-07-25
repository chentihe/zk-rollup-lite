package clients

import (
	"context"
	"crypto/ecdsa"
	"math/big"

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

func NewSigner(chainId *big.Int, priv string) (*Signer, error) {
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

	return &Signer{
		privateKey: privateKey,
		PublicKey:  publicKeyECDSA,
		Address:    signerAddress,
		Signer:     types.NewLondonSigner(chainId),
		Context:    context.Background(),
	}, nil
}

func (signer *Signer) GetAuth(ethClient *ethclient.Client) (*bind.TransactOpts, error) {
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

	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(800000)
	auth.GasPrice = gasPrice

	return auth, nil
}

func (signer *Signer) GenerateDynamicTx(ethClient *ethclient.Client, to *common.Address, data []byte, value *big.Int) (*types.Transaction, error) {
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

	// callMsg := ethereum.CallMsg{
	// 	From:      signer.Address,
	// 	To:        to,
	// 	Gas:       math.MaxInt64,
	// 	GasTipCap: tipCap,
	// 	GasFeeCap: feeCap,
	// 	Value:     value,
	// 	Data:      data,
	// }

	// gasLimit, err := ethClient.EstimateGas(signer.Context, callMsg)
	// if err != nil {
	// 	return nil, err
	// }

	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:   chainId,
			Nonce:     nonce,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       uint64(800000),
			To:        to,
			Value:     value,
			Data:      data,
		})

	return tx, nil
}

func (signer *Signer) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	return types.SignTx(tx, signer.Signer, signer.privateKey)
}
