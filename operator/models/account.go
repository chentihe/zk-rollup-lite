package models

import (
	"fmt"
	"math/big"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	AccountIndex int64  `gorm:"uniqueIndex"`
	PublicKey    string `gorm:"uniqueIndex"`
	L1Address    string
	Nonce        int64
	Balance      string
}

type AccountDto struct {
	AccountIndex int64
	PublicKey    string
	L1Address    string
	Nonce        int64
	Balance      *big.Int
}

func (dto *AccountDto) Copy() *AccountDto {
	balance := new(big.Int)
	balance.SetString(dto.Balance.String(), 10)
	return &AccountDto{
		AccountIndex: dto.AccountIndex,
		PublicKey:    dto.PublicKey,
		L1Address:    dto.L1Address,
		Nonce:        dto.Nonce,
		Balance:      balance,
	}
}

func (dto *AccountDto) ToModel() *Account {
	return &Account{
		AccountIndex: dto.AccountIndex,
		PublicKey:    dto.PublicKey,
		L1Address:    dto.L1Address,
		Nonce:        dto.Nonce,
		Balance:      dto.Balance.String(),
	}
}

func (model *Account) ToDto() (*AccountDto, error) {
	balance := new(big.Int)
	balance, ok := balance.SetString(model.Balance, 10)
	if !ok {
		return nil, fmt.Errorf("cannot convert string to big int for balance")
	}

	return &AccountDto{
		AccountIndex: model.AccountIndex,
		PublicKey:    model.PublicKey,
		L1Address:    model.L1Address,
		Nonce:        model.Nonce,
		Balance:      balance,
	}, nil
}
