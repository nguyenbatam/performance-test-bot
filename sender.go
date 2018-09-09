package main

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func Sender(signTx *types.Transaction) error {
	return SendRawTransaction(signTx, httpClient)
}
