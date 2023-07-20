package eventhandlers

import (
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type Deposit struct {
	User User
}

func (e *EventHandler) afterDeposit(vLog *types.Log) error {
	var deposit Deposit
	if err := e.abi.UnpackIntoInterface(&deposit, "Deposit", vLog.Data); err != nil {
		return err
	}

	user := deposit.User

	publicKey := babyjub.PublicKey(babyjub.Point{X: user.PublicKeyX, Y: user.PublicKeyY})

	tx, _, err := e.ethClient.TransactionByHash(e.context, vLog.TxHash)
	if err != nil {
		return err
	}

	sender, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		return err
	}

	account, err := e.accountService.GetAccountByIndex(user.Index)
	switch err {
	case daos.ErrAccountNotFound:
		// retrieve sender address from tx
		account = &models.Account{
			AccountIndex: user.Index,
			PublicKey:    publicKey.String(),
			Balance:      user.Balance,
			Nonce:        user.Nonce,
			L1Address:    sender.Hex(),
		}

		if err := e.accountService.CreateAccount(account); err != nil {
			return err
		}

		accountLeaf, err := tree.GenerateAccountLeaf(account)
		if err != nil {
			return err
		}

		if err := e.accountTree.Add(user.Index, accountLeaf); err != nil {
			return err
		}
	case daos.ErrSqlOperation:
		return err
	default:
		account.Balance = user.Balance
		account.Nonce = user.Nonce
		account.L1Address = sender.Hex()

		if err := e.accountService.UpdateAccount(account); err != nil {
			return err
		}

		if _, err := e.accountTree.UpdateAccountTree(account); err != nil {
			return err
		}
	}

	return nil
}
