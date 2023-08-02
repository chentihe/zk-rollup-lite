package daos

import "github.com/chentihe/zk-rollup-lite/operator/models"

type AccountDao interface {
	GetAccountByIndex(index int64) (*models.Account, error)
	GetAccountByPublicKey(comp string) (*models.Account, error)
	CreateAccount(account *models.Account) error
	UpdateAccount(account *models.Account) error
	DeleteAccountByIndex(index int64) error
	CreateAccountTable() error
	DropAccountTable() error
}
