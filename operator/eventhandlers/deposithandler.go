package eventhandlers

import (
	"context"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/accounttree"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/services"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
)

type Deposit struct {
	User User
}

func AfterDeposit(vLog *types.Log, accountService *services.AccountService, mt *merkletree.MerkleTree, contractAbi *abi.ABI, context context.Context, client *ethclient.Client) error {
	var deposit Deposit
	if err := contractAbi.UnpackIntoInterface(&deposit, "Deposit", vLog.Data); err != nil {
		return err
	}

	user := deposit.User

	publicKey := babyjub.PublicKey(babyjub.Point{X: user.PublicKeyX, Y: user.PublicKeyY})

	account, err := accountService.GetAccountByIndex(user.Index)

	// if err occurs, this account is a new user
	if err != nil {
		// retrieve sender address from tx
		tx, _, err := client.TransactionByHash(context, vLog.TxHash)
		if err != nil {
			return err
		}

		sender, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
		if err != nil {
			return err
		}

		account = &models.AccountModel{
			AccountIndex: user.Index,
			PublicKey:    publicKey.String(),
			Balance:      user.Balance,
			Nonce:        user.Nonce,
			L1Address:    sender.Hex(),
		}

		accountLeaf, err := accounttree.GenerateAccountLeaf(account)
		if err != nil {
			return err
		}

		accountService.CreateAccount(account)
		mt.Add(context, big.NewInt(user.Index), accountLeaf)
	} else {
		account.Balance = user.Balance
		account.Nonce = user.Nonce

		accountLeaf, err := accounttree.GenerateAccountLeaf(account)
		if err != nil {
			return err
		}

		accountService.UpdateAccount(account)
		mt.Update(context, big.NewInt(user.Index), accountLeaf)
	}

	return nil
}
