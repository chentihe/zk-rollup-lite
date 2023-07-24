package daos

import "github.com/chentihe/zk-rollup-lite/operator/models"

type AccountDao interface {
	GetAccountByIndex(index int64) (*models.Account, error)
	GetCurrentAccountIndex() (int64, error)
	CreateAccount(account *models.Account) error
	UpdateAccount(account *models.Account) error
	CreateAccountTable() error
	DropAccountTable() error
}
