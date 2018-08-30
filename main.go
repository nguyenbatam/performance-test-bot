package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"os"
	"sync"
	"time"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	NReq     = flag.Int("req", 1, "The number of transactions")
	Url      = flag.String("url", "http://localhost:8545", "That you want to connect")
	KeyFile  = flag.String("key-file", "key.json", "Key file name")
	Password = flag.String("password", "", "Keyfile password")
)
var wg sync.WaitGroup
var client *ethclient.Client
var nonce uint64
var unlockedKey *keystore.Key
var request []*types.Transaction

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
	client, err = ethclient.Dial(*Url)
	if err != nil {
		log.Fatal(err)
	}
	unlockedKey, _ = keystore.DecryptKey(data, *Password)
	ctx, _ := context.WithTimeout(context.Background(), 100000*time.Millisecond)
	nonce, _ = client.NonceAt(ctx, unlockedKey.Address, nil)
	request = make([]*types.Transaction, *NReq * *NWorkers)
	prepareData(*NReq, *NWorkers)
	attack(*NReq, *NWorkers)
}

func attack(nReq int, nWorkers int) {
	StartDispatcher(nWorkers)
	// Start the dispatcher.
	for {
		start := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Println("Start send ", len(request), "request ")
		for i := 0; i < len(request); i++ {
			if (request)[i] == nil {
				fmt.Println(i, (request)[i])
			}
			WorkQueue <- (request)[i]
		}
		prepareData(nReq, nWorkers)
		end := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Println("Done a round with time = ", end-start)
		if (end-start < 1000) {
			sleep := int(1000 + start - end)
			time.Sleep(time.Duration(sleep) * time.Microsecond)
		}
	}
}

func prepareData(nReq int, nWorkers int) {
	// Now, create all of our workers.
	fmt.Println("Prepare data")
	wg.Add(nWorkers)
	for workerIndex := 0; workerIndex < nWorkers; workerIndex++ {
		go func(workerIndex int) {
			start := nReq * workerIndex
			end := start + nReq
			var err error
			for i := start; i < end; i++ {
				tx := types.NewTransaction(uint64(i)+nonce, unlockedKey.Address, big.NewInt(int64(i+1)+int64(nonce)), 21000, big.NewInt(int64(100000+i)), nil)
				request[i], err = types.SignTx(tx, types.NewEIP155Signer(big.NewInt(89)), unlockedKey.PrivateKey)
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
