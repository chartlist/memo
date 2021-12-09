/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-11
 */
package tools

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xb10c/memo/core-rpc/btcjson"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/mysql"
	"github.com/0xb10c/memo/mysql/tables"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/urfave/cli"
)

func syncTxs(ctx *cli.Context) {
	logger.Info.Println("sync txs")
	mysql.CreateTxTables(time.Now().Unix())
	mysql.CreateFeeTable(time.Now().Unix())
	mysql.CreateScriptPubKeyResult(time.Now().Unix())
	mysql.CreateScriptSigTable(time.Now().Unix())
	mysql.CreateVinTable(time.Now().Unix())
	mysql.CreateVoutTable(time.Now().Unix())
}

func calculateTransferValue(tx btcjson.TxRawResult) float64 {
	var sum float64
	for _, vout := range tx.Vout {
		sum += vout.Value
	}
	return sum
}

func calculateTxFee(tx btcjson.TxRawResult) float64 {

	var inAmountSats, outAmountSats float64
	for _, vin := range tx.Vin {
		if vin.IsCoinBase() {
			return 0
		}

		txId, err := chainhash.NewHashFromStr(vin.Txid)
		if err != nil {
			logger.Error.Printf("new hash from str %s err:%s", vin.Txid, err.Error())
			return 0
		}

		txRaw, err := RPCClient.GetRawTransactionVerbose(txId)
		if err != nil {
			logger.Error.Printf("in calc tx fee, get raw transaction(txid:%s) verbose err:%s", txId, err.Error())
			return 0
		}
		inAmountSats += txRaw.Vout[vin.Vout].Value
	}

	for _, vout := range tx.Vout {
		outAmountSats += vout.Value
	}

	return inAmountSats - outAmountSats
}

func vins(tx btcjson.TxRawResult) string {
	if tx.Vin[0].IsCoinBase() {
		return ""
	}

	var res string
	for _, vin := range tx.Vin {
		res = res + fmt.Sprintf("%s-%d;", vin.Txid, vin.Vout)
	}
	return res[0 : len(res)-1]
}

func vouts(tx btcjson.TxRawResult) string {
	var res string
	for _, vout := range tx.Vout {
		res = res + fmt.Sprintf("%s-%d;", tx.Txid, vout.N)
	}
	return res[0 : len(res)-1]
}

func insertTxRecord(tx btcjson.TxRawResult, height int64, timestamp int64) (float64, error) {
	mysql.CreateTxTables(timestamp)

	fee := calculateTxFee(tx)
	Tx := tables.Transaction{
		Coinbase:      tx.Vin[0].IsCoinBase(),
		Hex:           tx.Hex,
		Txid:          tx.Txid,
		Hash:          tx.Hash,
		Size:          tx.Size,
		Vsize:         tx.Vsize,
		Weight:        tx.Weight,
		Version:       tx.Version,
		LockTime:      tx.LockTime,
		Value:         calculateTransferValue(tx),
		Vins:          vins(tx),
		Vouts:         vouts(tx),
		BlockHight:    height,
		BlockHash:     tx.BlockHash,
		Confirmations: tx.Confirmations,
		Fee:           fee,
		Time:          tx.Time,
		Blocktime:     tx.Blocktime,
	}

	err := mysql.BTC_DB.Create(&Tx).Error
	return fee, err
}

func insertFeeRecord(txid string, fee float64, timestamp int64) error {
	mysql.CreateFeeTable(timestamp)
	_fee := tables.Fee{
		Txid:      txid,
		Fee:       fee,
		SpentTime: timestamp,
	}
	return mysql.BTC_DB.Create(&_fee).Error
}

func scriptSigId(sig *btcjson.ScriptSig) string {
	if sig == nil {
		return ""
	}
	script, err := json.Marshal(sig)
	if err != nil {
		panic(err.Error())
	}
	sha1Inst := sha1.New()
	sha1Inst.Write(script)
	return fmt.Sprintf("%x", sha1Inst.Sum([]byte("")))
}

func witnesses(wits []string) string {
	if len(wits) == 0 {
		return ""
	}
	var res string
	for _, wit := range wits {
		res = res + wit + ";"
	}
	return res[0 : len(res)-1]
}

func insertVin(vin btcjson.Vin, txid string, timestamp int64, value float64) error {
	mysql.CreateVinTable(timestamp)
	txId, err := chainhash.NewHashFromStr(vin.Txid)
	if err != nil {
		logger.Error.Printf("new hash from str %s err:%s", vin.Txid, err.Error())
	}

	if vin.Txid != "" {
		txRaw, err := RPCClient.GetRawTransactionVerbose(txId)
		if err != nil {
			logger.Error.Printf("in insert vin, vin.Txid:%s, get raw transaction(txid:%s) verbose err:%s", vin.Txid, txId, err.Error())
			return err
		}
		value = txRaw.Vout[vin.Vout].Value
	}

	in := tables.Vin{
		Txid:        txid,
		Coinbase:    vin.Coinbase,
		FromTxid:    vin.Txid,
		Vout:        vin.Vout,
		Value:       value,
		ScriptSigId: scriptSigId(vin.ScriptSig),
		Sequence:    vin.Sequence,
		Witnesses:   witnesses(vin.Witness),
		VinTime:     timestamp,
	}

	return mysql.BTC_DB.Create(&in).Error
}

func scriptPubKeyId(scriptPK btcjson.ScriptPubKeyResult) string {
	script, err := json.Marshal(scriptPK)
	if err != nil {
		panic(err.Error())
	}
	sha1Inst := sha1.New()
	sha1Inst.Write(script)
	return fmt.Sprintf("%x", sha1Inst.Sum([]byte("")))
}

func insertVout(vout btcjson.Vout, txid string, timestamp int64) error {
	mysql.CreateVoutTable(timestamp)
	out := tables.Vout{
		Txid:           txid,
		Value:          vout.Value,
		N:              vout.N,
		ScriptPubKeyId: scriptPubKeyId(vout.ScriptPubKey),
		VoutTime:       timestamp,
		Spent:          false,
		SpentTime:      0,
	}

	return mysql.BTC_DB.Create(&out).Error
}

func addresses(addrs []string) string {
	var res string
	if len(addrs) == 0 {
		return ""
	}
	for _, addr := range addrs {
		res = res + addr + ";"
	}
	return res[0 : len(res)-1]
}

func insertScriptPubKeyResult(spkRes btcjson.ScriptPubKeyResult, timestamp int64) error {
	mysql.CreateScriptPubKeyResult(timestamp)
	spk := tables.ScriptPubKeyResult{
		ScriptPubKeyId: scriptPubKeyId(spkRes),
		Asm:            spkRes.Asm,
		Hex:            spkRes.Hex,
		ReqSigs:        spkRes.ReqSigs,
		Type:           spkRes.Type,
		Addresses:      addresses(spkRes.Addresses),
		SPKTime:        timestamp,
	}

	return mysql.BTC_DB.Create(&spk).Error
}

func insertScriptSig(script btcjson.ScriptSig, timestamp int64) error {
	mysql.CreateScriptSigTable(timestamp)
	ss := tables.ScriptSig{
		ScriptSigId: scriptSigId(&script),
		Asm:         script.Asm,
		Hex:         script.Hex,
		ScriptTime:  timestamp,
	}
	return mysql.BTC_DB.Create(&ss).Error
}
