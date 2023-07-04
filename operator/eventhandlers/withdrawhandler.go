package eventhandlers

import (
	"context"
	"fmt"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/database"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-merkletree-sql/v2"
)

type Withdraw struct {
	User User
}

func AfterWithdraw(vLog *types.Log, accountService *services.AccountService, mt *merkletree.MerkleTree, contractAbi *abi.ABI, context context.Context) error {
	var withdraw Withdraw
	if err := contractAbi.UnpackIntoInterface(&withdraw, "Withdraw", vLog.Data); err != nil {
		return err
	}

	user := withdraw.User

	account, err := accountService.GetAccountByIndex(user.Index)
	if err != nil {
		return fmt.Errorf("Invalid account")
	}

	account.Balance = user.Balance
	account.Nonce = user.Nonce

	accountLeaf, err := database.GenerateAccountLeaf(account)
	if err != nil {
		return err
	}

	accountService.UpdateAccount(account)
	mt.Update(context, big.NewInt(user.Index), accountLeaf)

	return nil
}
