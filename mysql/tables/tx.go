/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-09-30
 */
package tables

type Transaction struct {
	Coinbase      bool    `json:"coinbase"`
	Hex           string  `json:"hex"`
	Txid          string  `json:"txid" gorm:"index:idx_txid"`
	Hash          string  `json:"hash,omitempty"`
	Size          int32   `json:"size,omitempty"`
	Vsize         int32   `json:"vsize,omitempty"`
	Weight        int32   `json:"weight,omitempty"`
	Version       uint32  `json:"version"`
	LockTime      uint32  `json:"locktime"`
	Value         float64 `json:"value"`
	Vins          string  `json:"vins"`  //format: txid-N
	Vouts         string  `json:"vouts"` //format: txid-N
	BlockHight    int64   `json:"blockhight"`
	BlockHash     string  `json:"blockhash,omitempty"`
	Confirmations uint64  `json:"confirmations,omitempty"`
	Fee           float64 `json:"fee"`
	Time          int64   `json:"time,omitempty"`
	Blocktime     int64   `json:"blocktime,omitempty"`
}
type Fee struct {
	Txid      string  `json:"txid" gorm:"index:idx_txid"`
	Fee       float64 `json:"fee"`
	SpentTime int64   `json:"spenttime"`
}
type ScriptSig struct {
	ScriptSigId string `json:"scriptsigid" gorm:"index:idx_scriptsigid"`
	Asm         string `json:"asm"`
	Hex         string `json:"hex"`
	ScriptTime  int64  `json:"scripttime"`
}

type Vin struct {
	Txid        string  `json:"txid" gorm:"index:idx_txid"`
	Coinbase    string  `json:"coinbase"`
	FromTxid    string  `json:"fromtxid"`
	Vout        uint32  `json:"vout"`
	Value       float64 `json:"value"`
	ScriptSigId string  `json:"scriptSig"`
	Sequence    uint32  `json:"sequence"`
	Witnesses   string  `json:"txinwitnesses"`
	VinTime     int64   `json:"vintime"`
}

type ScriptPubKeyResult struct {
	ScriptPubKeyId string `json:"scriptpubkeyid" gorm:"index:idx_scriptpubkeyid"`
	Asm            string `json:"asm"`
	Hex            string `json:"hex,omitempty"`
	ReqSigs        int32  `json:"reqSigs,omitempty"`
	Type           string `json:"type"`
	Addresses      string `json:"addresses,omitempty"`
	SPKTime        int64  `json:"spktime"`
}

type Vout struct {
	Txid           string  `json:"txid" gorm:"index:idx_txid"`
	Value          float64 `json:"value"`
	N              uint32  `json:"n"`
	ScriptPubKeyId string  `json:"scriptpubkeyid"`
	VoutTime       int64   `json:"vouttime"`
	Spent          bool    `json:"spent"`
	SpentTime      int64   `json:"spenttime"`
}

func CreateTxTable(suffix string) {

}
