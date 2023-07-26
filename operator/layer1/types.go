package layer1

import "math/big"

type User struct {
	Index      *big.Int `json:"index"`
	PublicKeyX *big.Int `json:"publicKeyX"`
	PublicKeyY *big.Int `json:"publicKeyY"`
	Balance    *big.Int `json:"balance"`
	Nonce      *big.Int `json:"nonce"`
}

type Withdraw struct {
	User User
}

type Deposit struct {
	User User
}
