package main

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func Sender(signTx *types.Transaction) error {
	return client.SendTransaction(ctx, signTx)
}
