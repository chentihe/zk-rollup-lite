package daos

import (
	"errors"

	"github.com/chentihe/zk-rollup-lite/operator/models"
	"gorm.io/gorm"
)

const AccountTableName = `account`

type AccountDaoImpl struct {
	table string
	DB    *gorm.DB
}

func NewAccountDao(db *gorm.DB) AccountDao {
	return &AccountDaoImpl{
		table: AccountTableName,
		DB:    db,
	}
}

func (dao *AccountDaoImpl) CreateAccountTable() error {
	return dao.DB.AutoMigrate(models.Account{})
}

func (dao *AccountDaoImpl) DropAccountTable() error {
	return dao.DB.Migrator().DropTable(dao.table)
}

func (dao *AccountDaoImpl) GetAccountByIndex(index int64) (account *models.Account, err error) {
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

func (dao *AccountDaoImpl) CreateAccount(account *models.Account) (err error) {
	dbTx := dao.DB.Create(&account)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	return nil
}

func (dao *AccountDaoImpl) UpdateAccount(account *models.Account) (err error) {
	dbTx := dao.DB.Model(&models.Account{}).
		Where("account_index = ?", account.AccountIndex).
		Updates(&account)
	if dbTx.Error != nil {
		return dbTx.Error
	}

	return nil
}
