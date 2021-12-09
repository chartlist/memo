/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-10-18
 */
package tables

type Config struct {
	Id                   int32 `json:"id"`
	ScriptSigMonth       int32 `json:"scriptsigmonth"`
	VoutMonth            int32 `json:"voutmonth"`
	FeeMonth             int32 `json:"feemonth"`
	VinMonth             int32 `json:"vinmonth"`
	TxMonth              int32 `json:"txmonth"`
	ScriptPubKeyResMonth int32 `json:"scriptpubkeyresmonth"`
}
