/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-11
 */
package tools

import (
	"github.com/0xb10c/memo/cmd/utils"
	"github.com/urfave/cli"
)

var SyncCommand = cli.Command{
	Name:  "sync",
	Usage: "Display information about the chain",
	Subcommands: []cli.Command{
		{
			Action:    sync,
			Name:      "all",
			Usage:     "Display transaction information",
			ArgsUsage: "txHash",
			Flags: []cli.Flag{
				utils.FromFlag,
				utils.ToFlag,
			},
			Description: "Display transaction information",
		},
		{
			Action:      syncFees,
			Name:        "fee",
			Usage:       "Display transaction information",
			ArgsUsage:   "txHash",
			Flags:       []cli.Flag{},
			Description: "Display transaction information",
		},
		{
			Action:      syncBlocks,
			Name:        "block",
			Usage:       "Display transaction information",
			ArgsUsage:   "txHash",
			Flags:       []cli.Flag{},
			Description: "Display transaction information",
		},
		{
			Action:      syncTxs,
			Name:        "tx",
			Usage:       "Display transaction information",
			ArgsUsage:   "txHash",
			Flags:       []cli.Flag{},
			Description: "Display transaction information",
		},
	},
	Description: `Query information command can query information such as blocks, transactions, and transaction executions. 
You can use the ./Ontology info block --help command to view help information.`,
}
