/**
 * Description:
 * Author: Yihen.Liu
 * Create: 2021-09-30
 */
package tables

type Block struct {
	Hash          string  `json:"hash"`
	Confirmations int64   `json:"confirmations"`
	StrippedSize  int32   `json:"strippedsize"`
	Size          int32   `json:"size"`
	Weight        int32   `json:"weight"`
	Height        int64   `json:"height" gorm:"index:idx_height, unique"`
	Version       int32   `json:"version"`
	VersionHex    string  `json:"versionHex"`
	MerkleRoot    string  `json:"merkleroot"`
	Time          int64   `json:"time"`
	Nonce         uint32  `json:"nonce"`
	Bits          string  `json:"bits"`
	Difficulty    float64 `json:"difficulty"`
	PreviousHash  string  `json:"previousblockhash"`
	NTx           uint32  `json:"nTx"`
	Fee           float64 `json:"fee"`
	Value         float64 `json:"value"`
	NextHash      string  `json:"nextblockhash,omitempty"`
}

func CreateBlockTable(suffix string) {

}
