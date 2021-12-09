// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"log"

	rpcclient "github.com/0xb10c/memo/core-rpc"
)

func main() {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host: "api.anyblock.tools/bitcoin/bitcoin/mainnet/rpc/f14d7962-573e-49f8-a80d-ecb53f8dc135/",
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	txs, err := client.GetRawMempoolVerbose()
	if err != nil {
		log.Fatal(err)
	}
	for hash, tx := range txs {
		log.Println(hash, tx)
	}
	log.Printf("Block count: %d", blockCount)
}
