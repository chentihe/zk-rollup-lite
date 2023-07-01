package daos

import "github.com/chentihe/zk-rollup-lite/operator/models"

type AccountDao interface {
	GetAccountByIndex(index int64) (*models.AccountModel, error)
	CreateAccount(account *models.AccountModel) error
	UpdateAccount(account *models.AccountModel) error
}
