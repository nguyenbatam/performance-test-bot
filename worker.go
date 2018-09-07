package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"time"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan *types.Transaction) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:   id,
		Work: workerQueue,
	}
	return worker
}

type Worker struct {
	ID   int
	Work chan *types.Transaction
}

func (w Worker) Start() {
	go func() {
		for {
			select {
			case signTx := <-w.Work:
				start := time.Now().UnixNano()
				err := Sender(signTx)
				end := time.Now().UnixNano()
				fmt.Println("send",(end-start)/int64(time.Millisecond),signTx.Hash().Hex())
				if (err != nil) {
					fmt.Println(err,signTx)
				}
			}
		}
	}()
}
