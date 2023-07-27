package cli

import (
	"github.com/chentihe/zk-rollup-lite/operator/cmd/flags"
	"github.com/chentihe/zk-rollup-lite/operator/config/servicecontext"
	"github.com/urfave/cli/v2"
)

func DeleteNode(ctx *cli.Context, svc *servicecontext.ServiceContext) error {
	accountIndex := ctx.Int64(flags.AccountIndexFlag.Name)
	if err := svc.AccountService.DeleteAccountByIndex(accountIndex); err != nil {
		return err
	}
	return svc.AccountTree.Delete(accountIndex)
}
