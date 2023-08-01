package cli

import (
	"encoding/hex"

	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type User struct {
	privateKey *babyjub.PrivateKey
	PublicKey  *babyjub.PublicKey
	Index      int64
}

func NewUser(account *config.Account) (*User, error) {
	var k babyjub.PrivateKey
	_, err := hex.Decode(k[:], []byte(account.EddsaPrivKey))
	if err != nil {
		return nil, err
	}
	return &User{
		privateKey: &k,
		PublicKey:  k.Public(),
		Index:      account.Index,
	}, nil
}
