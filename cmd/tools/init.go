/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-12
 */
package tools

import (
	"github.com/0xb10c/memo/config"
	rpcclient "github.com/0xb10c/memo/core-rpc"
)

var RPCClient *rpcclient.Client

func init() {

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	var err error
	connCfg := &rpcclient.ConnConfig{
		Host: config.GetString("bitcoind.jsonrpc.host") + ":" + config.GetString("bitcoind.jsonrpc.port"),
		User: config.GetString("bitcoind.jsonrpc.username"),
		Pass: config.GetString("bitcoind.jsonrpc.password"),
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	RPCClient, err = rpcclient.New(connCfg)
	if err != nil {
		panic(err.Error())
	}
}
