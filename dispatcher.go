package main

import (
	"fmt"
	"net/http"
)

var WorkQueue = make(chan *http.Request)

func StartDispatcher(nworkers int) {
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkQueue)
		worker.Start()
	}
}
