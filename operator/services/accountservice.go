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

func (service *AccountService) GetAccountByIndex(index int64) (accountDto *models.AccountDto, err error) {
	account, err := service.AccountDao.GetAccountByIndex(index)
	if err != nil {
		return nil, err
	}

	accountDto, err = account.ToDto()
	if err != nil {
		return nil, err
	}

	return accountDto, nil
}

func (service *AccountService) GetCurrentAccountIndex() (amount int64, err error) {
	return service.AccountDao.GetCurrentAccountIndex()
}

func (service *AccountService) CreateAccount(accountDto *models.AccountDto) (err error) {
	return service.AccountDao.CreateAccount(accountDto.ToModel())
}

func (service *AccountService) UpdateAccount(accountDto *models.AccountDto) (err error) {
	return service.AccountDao.UpdateAccount(accountDto.ToModel())
}
