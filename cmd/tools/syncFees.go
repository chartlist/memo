/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-09-17
 */
package tools

import (
	"log"
	"time"

	"github.com/0xb10c/memo/cmd/utils"

	"github.com/urfave/cli"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/0xb10c/memo/logger"

	"github.com/0xb10c/memo/config"

	rpcclient "github.com/0xb10c/memo/core-rpc"
)

func syncFees(ctx *cli.Context) {
	if ctx.NArg() < 2 {
		logger.Warning.Println("args is not enough")
		return
	}
	from := ctx.Uint64(utils.GetFlagName(utils.FromFlag))
	to := ctx.Uint64(utils.GetFlagName(utils.ToFlag))

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host: config.GetString("bitcoind.jsonrpc.host") + ":" + config.GetString("bitcoind.jsonrpc.port"),
		User: config.GetString("bitcoind.jsonrpc.username"),
		Pass: config.GetString("bitcoind.jsonrpc.password"),
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	for i := from; i < to; i++ {
		// Get the current block hash at special height
		bHash, err := client.GetBlockHash(int64(i))
		if err != nil {
			logger.Error.Printf("get block hash err at height=%d, err:%s", i, err.Error())
			return
		}
		logger.Info.Printf("at height:%d, block hash:%s", i, bHash.String())
		blk, err := client.GetBlockVerboseTx(bHash)
		if err != nil {
			logger.Error.Printf("get verbose tx in block hash %s err:%s", i, bHash, err.Error())
			return
		}
		logger.Info.Printf("total txs:%d", len(blk.Tx))
		var totalFee float64
		for _, tx := range blk.Tx {
			var inAmountSats, outAmountSats float64
			//fmt.Println("================ tx id =============:", tx.Txid)
			for _, vin := range tx.Vin {
				if vin.IsCoinBase() {
					continue
				}

				txId, err := chainhash.NewHashFromStr(vin.Txid)
				if err != nil {
					logger.Error.Printf("new hash from str %s err:%s", vin.Txid, err.Error())
					return
				}

				txRaw, err := client.GetRawTransactionVerbose(txId)
				if err != nil {
					logger.Error.Printf("get raw transaction(txid:%s) verbose err:%s", txId, err.Error())
					return
				}
				inAmountSats += txRaw.Vout[vin.Vout].Value

				//fmt.Println("tx.Vin txId:", txId, ",index:", vin.Vout, ",value:", txRaw.Vout[vin.Vout].Value)
			}
			for _, vout := range tx.Vout {
				outAmountSats += vout.Value
			}

			if inAmountSats-outAmountSats > 0 {
				totalFee = totalFee + inAmountSats - outAmountSats
			}
			logger.Info.Printf("tx id:%s, block height:%d, tx fee:%f, out amout:%f, time:%s", tx.Txid, i, inAmountSats-outAmountSats, outAmountSats, time.Unix(blk.Time, 0))
		}

	}
}
