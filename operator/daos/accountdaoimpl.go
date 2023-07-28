package daos

import (
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"gorm.io/gorm"
)

const AccountTableName = `accounts`

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
	if dbTx.RowsAffected == 0 {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (dao *AccountDaoImpl) GetAccountByPublicKey(comp string) (account *models.Account, err error) {
	dbTx := dao.DB.
		Where("public_key = ?", comp).
		First(&account)
	if dbTx.RowsAffected == 0 {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (dao *AccountDaoImpl) GetCurrentAccountIndex() (amount int64, err error) {
	var count int64
	dbTx := dao.DB.Table(dao.table).Count(&count)
	if dbTx.Error != nil {
		return 0, ErrSqlOperation
	} else if dbTx.RowsAffected == 0 {
		return 0, ErrAccountNotFound
	}

	return count + 1, nil
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

func (dao *AccountDaoImpl) DeleteAccountByIndex(index int64) error {
	dbTx := dao.DB.
		Where("account_index = ?", index).
		Unscoped().
		Delete(&models.Account{})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	return nil
}
