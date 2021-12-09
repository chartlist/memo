/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-12
 */
package tools

import (
	"github.com/0xb10c/memo/cmd/utils"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/mysql"
	"github.com/urfave/cli"
)

func sync(ctx *cli.Context) {
	if ctx.NArg() < 2 {
		logger.Warning.Println("args is not enough")
		//return
	}
	from := ctx.Uint64(utils.GetFlagName(utils.FromFlag))
	to := ctx.Uint64(utils.GetFlagName(utils.ToFlag))
	logger.Info.Println(utils.GetFlagName(utils.FromFlag))
	logger.Info.Println(utils.GetFlagName(utils.ToFlag))
	logger.Info.Printf("from:%d, to:%d", from, to)

	for i := from; i < to; i++ {

		bHash, err := RPCClient.GetBlockHash(int64(i))
		if err != nil {
			logger.Error.Printf("get block hash err at height=%d, err:%s", i, err.Error())
			return
		}

		blk, err := RPCClient.GetBlockVerboseTx(bHash)
		if err != nil {
			logger.Error.Printf("get verbose tx in block hash %s err:%s", i, bHash, err.Error())
			return
		}

		TX := mysql.BTC_DB.Begin()

		if err := insertBlockRecord(blk); err != nil {
			TX.Rollback()
			panic(err.Error())
		}

		for _, tx := range blk.Tx {
			fee, err := insertTxRecord(tx, int64(i), blk.Time)
			if err != nil {
				TX.Rollback()
				panic(err.Error())
			}
			if err := insertFeeRecord(tx.Txid, fee, blk.Time); err != nil {
				TX.Rollback()
				panic(err.Error())
			}
			for _, vin := range tx.Vin {
				if err := insertVin(vin, tx.Txid, blk.Time, tx.Vout[0].Value); err != nil {
					TX.Rollback()
					panic(err.Error())
				}
				if vin.ScriptSig != nil {
					if err := insertScriptSig(*vin.ScriptSig, blk.Time); err != nil {
						TX.Rollback()
						panic(err.Error())
					}
				}
			}

			for _, vout := range tx.Vout {
				if err := insertVout(vout, tx.Txid, blk.Time); err != nil {
					TX.Rollback()
					panic(err.Error())
				}
				if err := insertScriptPubKeyResult(vout.ScriptPubKey, blk.Time); err != nil {
					TX.Rollback()
					panic(err.Error())
				}
			}

		}

		if err := TX.Commit().Error; err != nil {
			panic(err.Error())
		}

		logger.Info.Printf("block hash:%s, height:%d, tx cnt:%d, finished.", bHash, i, len(blk.Tx))
	}
}
