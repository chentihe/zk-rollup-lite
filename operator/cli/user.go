package cli

import (
	"encoding/hex"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

type User struct {
	privateKey *babyjub.PrivateKey
	PublicKey  *babyjub.PublicKey
}

func NewUser(privKey string) (*User, error) {
	var k babyjub.PrivateKey
	_, err := hex.Decode(k[:], []byte(privKey))
	if err != nil {
		return nil, err
	}
	return &User{
		privateKey: &k,
		PublicKey:  k.Public(),
	}, nil
}
