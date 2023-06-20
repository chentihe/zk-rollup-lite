package models

import "gorm.io/gorm"

type AccountModel struct {
	gorm.Model
	AccountIndex int64  `gorm:"uniqueIndex"`
	AccountName  string `gorm:"uniqueIndex"`
	PublicKey    string `gorm:"uniqueIndex"`
	// AccountNameHash string `gorm:"uniqueIndex"`
	L1Address string
	Nonce     string
}
