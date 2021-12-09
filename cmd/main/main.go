/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-11
 */
package main

import (
	"os"
	"runtime"

	"github.com/0xb10c/memo/cmd/api"
	"github.com/0xb10c/memo/cmd/daemon"

	"github.com/0xb10c/memo/cmd/tools"

	"github.com/0xb10c/memo/cmd/utils"

	"github.com/urfave/cli"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "Memo CLI"
	app.Action = startMemo
	app.Copyright = "Copyright in 20/10/2021 The Memo Authors"
	app.Commands = []cli.Command{
		tools.SyncCommand,
		api.APICommand,
		daemon.DaemonCommand,
	}
	app.Flags = []cli.Flag{
		utils.FromFlag,
		utils.ToFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func startMemo(ctx *cli.Context) {
}
