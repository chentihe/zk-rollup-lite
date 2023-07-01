package services

import (
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/models"
)

type AccountService struct {
	AccountDao daos.AccountDao
}

func NewAccountService(accountDao *daos.AccountDao) *AccountService {
	return &AccountService{
		AccountDao: *accountDao,
	}
}

func (service *AccountService) GetAccountByIndex(index int64) (account *models.AccountModel, err error) {
	return service.AccountDao.GetAccountByIndex(index)
}

func (service *AccountService) CreateAccount(account *models.AccountModel) (err error) {
	return service.AccountDao.CreateAccount(account)
}

func (service *AccountService) UpdateAccount(account *models.AccountModel) (err error) {
	return service.AccountDao.UpdateAccount(account)
}
