package flags

import "github.com/urfave/cli/v2"

var (
	DepositAmountFlag = &cli.StringFlag{
		Name:    "deposit",
		Aliases: []string{"d"},
		Usage:   "the deposit amount",
	}
	WithdrawAmountFlag = &cli.StringFlag{
		Name:    "withdraw",
		Aliases: []string{"w"},
		Usage:   "the withdraw amount",
	}
	AccountIndexFlag = &cli.Int64Flag{
		Name:    "account",
		Aliases: []string{"a"},
		Usage:   "the account index",
	}
)
