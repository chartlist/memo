/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-11
 */
package tools

import (
	"github.com/0xb10c/memo/core-rpc/btcjson"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/mysql"
	"github.com/0xb10c/memo/mysql/tables"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/urfave/cli"
)

func syncBlocks(ctx *cli.Context) {
	logger.Info.Println("sync blocks...")
	mysql.CreateBlockTables()
}

func calculateValue(txs []btcjson.TxRawResult) float64 {
	var sum float64
	for _, tx := range txs {
		/*
			if tx.Vin[0].IsCoinBase() == true {
				continue
			}
		*/
		for _, vout := range tx.Vout {
			sum += vout.Value
		}
	}
	return sum
}
func calculateBlockFee(blk *btcjson.GetBlockVerboseTxResult) float64 {

	var totalFee float64
	for _, tx := range blk.Tx {
		var inAmountSats, outAmountSats float64
		for _, vin := range tx.Vin {
			if vin.IsCoinBase() {
				continue
			}

			txId, err := chainhash.NewHashFromStr(vin.Txid)
			if err != nil {
				logger.Error.Printf("new hash from str %s err:%s", vin.Txid, err.Error())
				return 0
			}

			txRaw, err := RPCClient.GetRawTransactionVerbose(txId)
			if err != nil {
				logger.Error.Printf("in calc block fee, get raw transaction(txid:%s) verbose err:%s", txId, err.Error())
				return 0
			}
			inAmountSats += txRaw.Vout[vin.Vout].Value
		}

		for _, vout := range tx.Vout {
			outAmountSats += vout.Value
		}

		if inAmountSats-outAmountSats > 0 {
			totalFee = totalFee + inAmountSats - outAmountSats
		}
	}
	return totalFee
}

func insertBlockRecord(blk *btcjson.GetBlockVerboseTxResult) error {
	mysql.CreateBlockTables()
	block := tables.Block{
		Hash:          blk.Hash,
		Confirmations: blk.Confirmations,
		StrippedSize:  blk.StrippedSize,
		Size:          blk.Size,
		Weight:        blk.Weight,
		Height:        blk.Height,
		Version:       blk.Version,
		VersionHex:    blk.VersionHex,
		MerkleRoot:    blk.MerkleRoot,
		Time:          blk.Time,
		Nonce:         blk.Nonce,
		Bits:          blk.Bits,
		Difficulty:    blk.Difficulty,
		PreviousHash:  blk.PreviousHash,
		NTx:           uint32(len(blk.Tx)),
		Fee:           calculateBlockFee(blk),
		Value:         calculateValue(blk.Tx),
		NextHash:      blk.NextHash,
	}
	return mysql.BTC_DB.Create(&block).Error
}
