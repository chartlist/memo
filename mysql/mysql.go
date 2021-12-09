/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-09-30
 */

package mysql

import (
	"fmt"
	"time"

	"github.com/0xb10c/memo/config"
	"github.com/0xb10c/memo/logger"
	"github.com/0xb10c/memo/mysql/tables"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	BTC_DB                  *gorm.DB // it is no need to mind closing action.
	SCRIPT_SIG_MONTH        = time.Month(0)
	VOUT_MONTH              = time.Month(0)
	FEE_MONTH               = time.Month(0)
	VIN_MONTH               = time.Month(0)
	TX_MONTH                = time.Month(0)
	SCRIPT_PUBKEY_RES_MONTH = time.Month(0)
)

func init() {
	//CREATE DATABASE IF NOT EXISTS bitcoin DEFAULT CHARSET utf8 COLLATE utf8_general_ci
	var err error
	conn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetString("mysql.user"),
		config.GetString("mysql.password"),
		config.GetString("mysql.host"),
		config.GetString("mysql.port"),
		config.GetString("mysql.db"),
	)
	//conn := "root:yihen382465@(localhost:3306)/bitcoin?charset=utf8mb4&parseTime=True&loc=Local"
	BTC_DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		logger.Error.Println("create mysql connection err:", err.Error())
		return
	} else {
		logger.Info.Println("create mysql connection successed.")
	}

	if BTC_DB.Migrator().HasTable(&tables.Config{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.Config{}); err != nil {
			panic("craete config table err:" + err.Error())
		}

		conf := &tables.Config{
			Id:                   1,
			ScriptSigMonth:       0,
			VoutMonth:            0,
			FeeMonth:             0,
			VinMonth:             0,
			TxMonth:              0,
			ScriptPubKeyResMonth: 0,
		}

		if err := BTC_DB.Create(&conf).Error; err != nil {
			panic("create config err:" + err.Error())
		}
	}

	conf := tables.Config{}
	if err := BTC_DB.First(&conf).Error; err == nil {
		FEE_MONTH = time.Month(conf.FeeMonth)
		VOUT_MONTH = time.Month(conf.VoutMonth)
		VIN_MONTH = time.Month(conf.VinMonth)
		SCRIPT_SIG_MONTH = time.Month(conf.ScriptSigMonth)
		SCRIPT_PUBKEY_RES_MONTH = time.Month(conf.ScriptPubKeyResMonth)
		TX_MONTH = time.Month(conf.TxMonth)
	} else {
		panic("get config err:" + err.Error())
	}
	return
}

func CreateBlockTables() {
	if BTC_DB.Migrator().HasTable(&tables.Block{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.Block{}); err != nil {
			panic("create block table err:" + err.Error())
		}
	}
}

func CreateScriptSigTable(timestamp int64) {
	if BTC_DB.Migrator().HasTable(&tables.ScriptSig{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.ScriptSig{}); err != nil {
			panic("create scriptsig table err:" + err.Error())
		} else {
			return
		}
	}

	t := time.Unix(timestamp, 0)
	if t.Month() != SCRIPT_SIG_MONTH {
		if err := BTC_DB.Migrator().RenameTable("script_sigs", fmt.Sprintf("script_sigs_%d_%d", t.Year(), t.Month())); err != nil {
			panic(fmt.Sprintf("rename script sigs err, year:%d, month:%d", t.Year(), t.Month()))
		}

		if err := BTC_DB.Migrator().CreateTable(&tables.ScriptSig{}); err != nil {
			panic(fmt.Sprintf("re-create script sigs table err:%s, year:%d, month:%d", err.Error(), t.Year(), t.Month()))
		}

		SCRIPT_SIG_MONTH = t.Month()
		if err := BTC_DB.Model(&tables.Config{}).Where("id=?", 1).Update("script_sig_month", int32(t.Month())).Error; err != nil {
			panic("update script sig month err:" + err.Error())
		}
	}
}

func CreateVoutTable(timestamp int64) {
	if BTC_DB.Migrator().HasTable(&tables.Vout{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.Vout{}); err != nil {
			panic("create vout table err:" + err.Error())
		} else {
			return
		}
	}

	t := time.Unix(timestamp, 0)
	if t.Month() != VOUT_MONTH {
		if err := BTC_DB.Migrator().RenameTable("vouts", fmt.Sprintf("vouts_%d_%d", t.Year(), t.Month())); err != nil {
			panic(fmt.Sprintf("rename fees err, year:%d, month:%d", t.Year(), t.Month()))
		}

		if err := BTC_DB.Migrator().CreateTable(&tables.Vout{}); err != nil {
			panic(fmt.Sprintf("re-create vouts table err:%s, year:%d, month:%d", err.Error(), t.Year(), t.Month()))
		}

		VOUT_MONTH = t.Month()
		if err := BTC_DB.Model(&tables.Config{}).Where("id=?", 1).Update("vout_month", int32(t.Month())).Error; err != nil {
			panic("update vout month err:" + err.Error())
		}
	}
}

func CreateFeeTable(timestamp int64) {
	if BTC_DB.Migrator().HasTable(&tables.Fee{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.Fee{}); err != nil {
			panic("create fee table err:" + err.Error())
		} else {
			return
		}
	}

	t := time.Unix(timestamp, 0)
	if t.Month() != FEE_MONTH {
		if err := BTC_DB.Migrator().RenameTable("fees", fmt.Sprintf("fees_%d_%d", t.Year(), t.Month())); err != nil {
			panic(fmt.Sprintf("rename fees err, year:%d, month:%d", t.Year(), t.Month()))
		}

		if err := BTC_DB.Migrator().CreateTable(&tables.Fee{}); err != nil {
			panic(fmt.Sprintf("re-create fees table err:%s, year:%d, month:%d", err.Error(), t.Year(), t.Month()))
		}

		FEE_MONTH = t.Month()
		if err := BTC_DB.Model(&tables.Config{}).Where("id=?", 1).Update("fee_month", int32(t.Month())).Error; err != nil {
			panic("update fee month err:" + err.Error())
		}
	}
}

func CreateVinTable(timestamp int64) {
	if BTC_DB.Migrator().HasTable(&tables.Vin{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.Vin{}); err != nil {
			panic("create fee table err:" + err.Error())
		} else {
			return
		}
	}

	t := time.Unix(timestamp, 0)
	if t.Month() != VIN_MONTH {
		if err := BTC_DB.Migrator().RenameTable("vins", fmt.Sprintf("vins_%d_%d", t.Year(), t.Month())); err != nil {
			panic(fmt.Sprintf("rename vins err, year:%d, month:%d", t.Year(), t.Month()))
		}

		if err := BTC_DB.Migrator().CreateTable(&tables.Vin{}); err != nil {
			panic(fmt.Sprintf("re-create vins table err:%s, year:%d, month:%d", err.Error(), t.Year(), t.Month()))
		}

		VIN_MONTH = t.Month()
		if err := BTC_DB.Model(&tables.Config{}).Where("id=?", 1).Update("vin_month", int32(t.Month())).Error; err != nil {
			panic("update vin month err:" + err.Error())
		}
	}
}

func CreateScriptPubKeyResult(timestamp int64) {
	if BTC_DB.Migrator().HasTable(&tables.ScriptPubKeyResult{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.ScriptPubKeyResult{}); err != nil {
			panic("create script pubkey result table err:" + err.Error())
		} else {
			return
		}
	}

	t := time.Unix(timestamp, 0)
	if t.Month() != SCRIPT_PUBKEY_RES_MONTH {
		if err := BTC_DB.Migrator().RenameTable("script_pub_key_results", fmt.Sprintf("script_pub_key_results_%d_%d", t.Year(), t.Month())); err != nil {
			panic(fmt.Sprintf("rename script_pub_key_results err, year:%d, month:%d", t.Year(), t.Month()))
		}

		if err := BTC_DB.Migrator().CreateTable(&tables.ScriptPubKeyResult{}); err != nil {
			panic(fmt.Sprintf("re-create script_pub_key_results table err:%s, year:%d, month:%d", err.Error(), t.Year(), t.Month()))
		}

		SCRIPT_PUBKEY_RES_MONTH = t.Month()
		if err := BTC_DB.Model(&tables.Config{}).Where("id=?", 1).Update("script_pub_key_res_month", int32(t.Month())).Error; err != nil {
			panic("update script pub key res month err:" + err.Error())
		}
	}
}

func CreateTxTables(timestamp int64) {
	if BTC_DB.Migrator().HasTable(&tables.Transaction{}) == false {
		if err := BTC_DB.Migrator().CreateTable(&tables.Transaction{}); err != nil {
			panic("create tx table err:" + err.Error())
		} else {
			return
		}
	}

	t := time.Unix(timestamp, 0)
	if t.Month() != TX_MONTH {
		if err := BTC_DB.Migrator().RenameTable("transactions", fmt.Sprintf("transactions_%d_%d", t.Year(), t.Month())); err != nil {
			panic(fmt.Sprintf("rename transactions err, year:%d, month:%d", t.Year(), t.Month()))
		}

		if err := BTC_DB.Migrator().CreateTable(&tables.Transaction{}); err != nil {
			panic(fmt.Sprintf("re-create tx table err:%s, year:%d, month:%d", err.Error(), t.Year(), t.Month()))
		}

		TX_MONTH = t.Month()
		if err := BTC_DB.Model(&tables.Config{}).Where("id=?", 1).Update("tx_month", int32(t.Month())).Error; err != nil {
			panic("update tx month err:" + err.Error())
		}
	}
}
