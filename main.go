package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"os"
	"sync"
	"time"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"net"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	NReq     = flag.Int("req", 1, "The number of transactions")
	Url      = flag.String("url", "http://localhost:8545", "That you want to connect")
	KeyFile  = flag.String("key-file", "key.json", "Key file name")
	Password = flag.String("password", "", "Keyfile password")
)
var wg sync.WaitGroup
var nonce uint64
var unlockedKey *keystore.Key
var request []*http.Request
var httpClient *http.Client

func main() {
	flag.Parse()
	key_file := *KeyFile
	fmt.Println(key_file, *NWorkers, *Url, )
	if _, err := os.Stat(key_file); err != nil {
		fmt.Println(err)
	}
	// get key file
	data, err := ioutil.ReadFile(key_file)
	if err != nil {
		fmt.Println(err)
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient = &http.Client{Transport: transport}
	unlockedKey, _ = keystore.DecryptKey(data, *Password)
	nonce, err = GetTransactionCount(unlockedKey.Address, httpClient)
	fmt.Println(unlockedKey.Address.Hex(), "nonce", nonce, time.Now().Format(time.RFC3339), err)
	balance, err := GetBalance(unlockedKey.Address, httpClient)
	fmt.Println(unlockedKey.Address.Hex(), "balance", balance, time.Now().Format(time.RFC3339), err)
	request = make([]*http.Request, *NReq * *NWorkers)
	prepareData(*NReq, *NWorkers)
	attack(*NReq, *NWorkers)
}

func attack(nReq int, nWorkers int) {
	// Start the dispatcher.
	StartDispatcher(nWorkers)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	stop := false
	go func() {
		select {
		case <-c:
			fmt.Println("Waiting Stop")
			stop = true
		}
	}()
	for (!stop) {

		start := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Println("Start send ", len(request), "request ")
		for i := 0; i < len(request); i++ {
			WorkQueue <- request[i]
		}
		balance, _ := GetBalance(unlockedKey.Address, httpClient)
		prepare := time.Now().UnixNano() / int64(time.Millisecond)
		prepareData(nReq, nWorkers)
		end := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Println("Done a round with time = ", prepare-start, "prepare data ", end-prepare, unlockedKey.Address.Hex(), "balance", balance, time.Now().Format(time.RFC3339))
		if (end-start < 1000) {
			sleep := int(1000 + start - end)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
}

func prepareData(nReq int, nWorkers int) {
	fmt.Println("Prepare data")
	wg.Add(nWorkers)
	for workerIndex := 0; workerIndex < nWorkers; workerIndex++ {
		go func(workerIndex int) {
			start := nReq * workerIndex
			end := start + nReq
			var err error
			for i := start; i < end; i++ {
				request[i], err = CreateRequestSendRawTransaction(uint64(i)+nonce, unlockedKey.Address)
				if err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}(workerIndex)
	}
	wg.Wait()
	nonce = nonce + uint64(nReq*nWorkers)
}
