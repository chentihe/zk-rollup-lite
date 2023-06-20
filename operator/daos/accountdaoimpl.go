package daos

import (
	"errors"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"gorm.io/gorm"
)

type AccountDaoImpl struct {
	DB *gorm.DB
}

func NewAccountDao(db *gorm.DB) AccountDao {
	return &AccountDaoImpl{
		DB: db,
	}
}

func (dao *AccountDaoImpl) GetAccountByIndex(index string) (account *models.AccountModel, err error) {
	dbTx := dao.DB.
		Where("account_index = ?", index).
		First(&account)
	if dbTx.Error != nil {
		return nil, errors.New("db error: sql operation")
	} else if dbTx.RowsAffected == 0 {
		return nil, errors.New("db error: not found")
	}

	return account, nil
}

func (dao *AccountDaoImpl) CreateAccount(account *models.AccountModel) (err error) {
	dbTx := dao.DB.Create(&account)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	return nil
}

func (dao *AccountDaoImpl) UpdateAccount(account *models.AccountModel) (err error) {
	dbTx := dao.DB.Model(&models.AccountModel{}).
		Where("account_index = ?", account.AccountIndex).
		Updates(&account)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	return nil
}
