package flags

import "github.com/urfave/cli/v2"

var (
	AmountFlag = &cli.StringFlag{
		Name:    "amount",
		Aliases: []string{"a"},
		Usage:   "the deposit/withdraw amount",
	}
	AccountIndexFlag = &cli.Int64Flag{
		Name:    "account",
		Aliases: []string{"i"},
		Usage:   "the account index",
	}
	NodeFlag = &cli.StringFlag{
		Name:        "node",
		Aliases:     []string{"n"},
		Usage:       "the executing node",
		Value:       "anvil",
		DefaultText: "default node is anvil from foundry",
	}
)
