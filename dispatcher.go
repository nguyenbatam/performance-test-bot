package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
)

var WorkQueue = make(chan *types.Transaction)

func StartDispatcher(nworkers int) {
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkQueue)
		worker.Start()
	}
}
