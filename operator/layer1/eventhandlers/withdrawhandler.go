package eventhandlers

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/chentihe/zk-rollup-lite/operator/tree"
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
		return err
	}

	account.Balance = user.Balance
	account.Nonce = user.Nonce

	if err := accountService.UpdateAccount(account); err != nil {
		return err
	}

	accountLeaf, err := tree.GenerateAccountLeaf(account)
	if err != nil {
		return err
	}

	if _, err := mt.Update(context, big.NewInt(user.Index), accountLeaf); err != nil {
		return err
	}

	return nil
}
