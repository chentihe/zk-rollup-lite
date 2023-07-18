package clients

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Signer struct {
	privateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

func NewSigner() (*Signer, error) {
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
	}, nil
}

func (signer *Signer) GetAuth(ethClient *ethclient.Client, context context.Context) (*bind.TransactOpts, error) {
	nonce, err := ethClient.PendingNonceAt(context, signer.Address)
	if err != nil {
		return nil, err
	}

	gasPrice, err := ethClient.SuggestGasPrice(context)
	if err != nil {
		return nil, err
	}

	chainId, err := ethClient.ChainID(context)
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
