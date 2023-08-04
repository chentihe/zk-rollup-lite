package cli

import (
	"context"
	"log"
	"math/big"

	"github.com/chentihe/zk-rollup-lite/operator/cmd/flags"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/config/servicecontext"
	"github.com/chentihe/zk-rollup-lite/operator/daos"
	"github.com/chentihe/zk-rollup-lite/operator/layer1/clients"
	"github.com/chentihe/zk-rollup-lite/operator/models"
	"github.com/chentihe/zk-rollup-lite/operator/txutils"
	"github.com/urfave/cli/v2"
)

func Deposit(ctx *cli.Context, context context.Context, config *config.Config, svc *servicecontext.ServiceContext) error {
	accountIndex := ctx.Int64(flags.AccountIndexFlag.Name)
	account := config.Accounts[accountIndex]
	signer, err := clients.NewSigner(context, account.EcdsaPrivKey, svc.EthClient)
	if err != nil {
		return err
	}

	user, err := NewUser(account)
	if err != nil {
		return err
	}

	depositAmount := txutils.ToWei(ctx.String(flags.AmountFlag.Name), 18)

	comp := user.PublicKey.String()

	accountDto, err := svc.AccountService.GetAccountByPublicKey(comp)
	if err == daos.ErrAccountNotFound {
		// index will be updated once the deposit event emits
		accountDto = &models.AccountDto{
			PublicKey: comp,
			Balance:   depositAmount,
			Nonce:     0,
		}
	} else {
		accountDto.Balance = new(big.Int).Add(accountDto.Balance, depositAmount)
	}

	data, err := svc.Abi.Pack("deposit", user.PublicKey.X, user.PublicKey.Y)
	if err != nil {
		log.Printf("Cannot pack rollup call data: %v\n", err)
	}

	tx, err := signer.GenerateLegacyTx(svc.RollUpAddress, data, depositAmount)
	if err != nil {
		return err
	}

	signTx, err := signer.SignTx(tx)
	if err != nil {
		return err
	}

	if err = svc.EthClient.SendTransaction(context, signTx); err != nil {
		return err
	}
	log.Printf("Deposit success: %s\n", signTx.Hash().Hex())

	return nil
}
