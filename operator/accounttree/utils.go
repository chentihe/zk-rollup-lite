package accounttree

import (
	"encoding/hex"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
)

func GenerateAccountLeaf(account *models.Account) (*big.Int, error) {
	publicKey, err := DecodePublicKeyFromString(account.PublicKey)
	if err != nil {
		return nil, err
	}

	hashedLeaf, err := poseidon.Hash([]*big.Int{
		publicKey.X,
		publicKey.Y,
		account.Balance,
		big.NewInt(account.Nonce),
	})
	if err != nil {
		return nil, err
	}

	return hashedLeaf, nil
}

func DecodePublicKeyFromString(comp string) (*babyjub.PublicKey, error) {
	bytesPublicKey, err := hex.DecodeString(comp)

	publicKeyComp := babyjub.PublicKeyComp(bytesPublicKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := publicKeyComp.Decompress()
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func StringifyPublicKey(comp string) (*[2]string, error) {
	publicKey, err := DecodePublicKeyFromString(comp)
	if err != nil {
		return nil, err
	}
	x := publicKey.X.String()
	y := publicKey.Y.String()
	return &[2]string{x, y}, nil
}
