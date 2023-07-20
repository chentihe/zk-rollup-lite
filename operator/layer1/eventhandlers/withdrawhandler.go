package eventhandlers

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type Withdraw struct {
	User User
}

func (e *EventHandler) afterWithdraw(vLog *types.Log) error {
	var withdraw Withdraw
	if err := e.abi.UnpackIntoInterface(&withdraw, "Withdraw", vLog.Data); err != nil {
		return err
	}

	user := withdraw.User

	account, err := e.accountService.GetAccountByIndex(user.Index)
	if err != nil {
		return err
	}

	account.Balance = user.Balance
	account.Nonce = user.Nonce

	if err := e.accountService.UpdateAccount(account); err != nil {
		return err
	}

	if _, err := e.accountTree.UpdateAccountTree(account); err != nil {
		return err
	}

	return nil
}
