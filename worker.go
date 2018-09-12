package main

import (
	"fmt"
	"net/http"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan *http.Request) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:   id,
		Work: workerQueue,
	}
	return worker
}

type Worker struct {
	ID   int
	Work chan *http.Request
}

func (w Worker) Start() {
	go func() {
		for {
			select {
			case signTx := <-w.Work:
				for i := 0; i < 3; i++ {
					err := SendRequest(signTx, httpClient)
					if (err != nil) {
						fmt.Println(err, signTx)
					} else {
						break
					}

				}
			}
		}
	}()
}
