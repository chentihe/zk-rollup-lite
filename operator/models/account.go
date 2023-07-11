package models

import (
	"math/big"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	AccountIndex int64  `gorm:"uniqueIndex"`
	PublicKey    string `gorm:"uniqueIndex"`
	L1Address    string
	Nonce        int64
	Balance      *big.Int
}
