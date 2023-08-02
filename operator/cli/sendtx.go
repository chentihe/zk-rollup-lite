package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chentihe/zk-rollup-lite/operator/cmd/flags"
	"github.com/chentihe/zk-rollup-lite/operator/config"
	"github.com/chentihe/zk-rollup-lite/operator/config/servicecontext"
	"github.com/chentihe/zk-rollup-lite/operator/txutils"
	"github.com/urfave/cli/v2"
)

type Accounts struct {
	Sender    *User
	Recipient *User
}

func SendTx(ctx *cli.Context, context context.Context, config *config.Config, svc *servicecontext.ServiceContext) error {
	accountIndex := ctx.Int64(flags.AccountIndexFlag.Name)

	var accounts Accounts
	for i, account := range config.Accounts {
		if i == int(accountIndex) {
			sender, err := NewUser(account)
			if err != nil {
				return err
			}
			accounts.Sender = sender
		} else {
			recipient, err := NewUser(account)
			if err != nil {
				return err
			}
			accounts.Recipient = recipient
		}
	}

	transferAmount := ToWei(ctx.String(flags.AmountFlag.Name), 18)

	accountDto, err := svc.AccountService.GetAccountByPublicKey(accounts.Sender.PublicKey.String())
	if err != nil {
		return err
	}

	tx := txutils.TransactionInfo{
		From:   accounts.Sender.Index,
		To:     accounts.Recipient.Index,
		Amount: transferAmount,
		Fee:    txutils.Fee,
		Nonce:  accountDto.Nonce + 1,
	}

	hashedMsg, err := tx.HashMsg()
	if err != nil {
		return err
	}
	tx.Signature = accounts.Sender.privateKey.SignMimc7(hashedMsg)

	requestBody, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8000/api/v1/send", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))

	return nil
}
