package main

import (
	"github.com/ethereum/go-ethereum/core/types"
	"context"
	"time"
)

func Sender(signTx *types.Transaction) error {
	ctx, _ := context.WithTimeout(context.Background(), 100000*time.Millisecond)
	return client.SendTransaction(ctx, signTx)
}
