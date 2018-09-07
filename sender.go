package main

import (
	"github.com/ethereum/go-ethereum/core/types"
	"context"
	"time"
	"math/rand"
)

func Sender(signTx *types.Transaction) error {
	rand := rand.Intn(len(clients))
	ctx, _ := context.WithTimeout(context.Background(), 100000*time.Millisecond)
	return clients[rand].SendTransaction(ctx, signTx)
}
