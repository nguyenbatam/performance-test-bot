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
	"github.com/ethereum/go-ethereum"
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
	balance, _ := client.BalanceAt(ctx, unlockedKey.Address, nil)
	gasprice, _ := client.EstimateGas(ctx, ethereum.CallMsg{})
	fmt.Println(unlockedKey.Address.Hex(), "balance", balance, "gasprice", gasprice, time.Now().Format(time.RFC3339))
	attack(*NReq, *NWorkers)
}

func attack(nReq int, nWorkers int) {
	// Start the dispatcher.
	StartDispatcher(nWorkers)
	c := make(chan os.Signal)
	stop := false
	for (!stop) {
		request := make([]*types.Transaction, nReq*nWorkers)
		prepareData(request)
		start := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Println("Start send ", len(request), "request ")
		for i := 0; i < len(request); i++ {
			WorkQueue <- request[i]
		}
		select {
		case <-c:
			stop = true
			fmt.Println("Waiting Stop")
		}
		ctx, _ := context.WithTimeout(context.Background(), 100000*time.Millisecond)
		balance, _ := client.BalanceAt(ctx, unlockedKey.Address, nil)
		end := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Println("Done a round with time = ", end-start, unlockedKey.Address.Hex(), "balance", balance, time.Now().Format(time.RFC3339))
		if (end-start < 1000) {
			sleep := int(1000 + start - end)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
}

func prepareData(request []*types.Transaction) {
	// Now, create all of our workers.
	fmt.Println("Prepare data")
	var err error
	for i := 0; i < len(request); i++ {
		tx := types.NewTransaction(nonce, unlockedKey.Address, big.NewInt(1), 21000, big.NewInt(1), nil)
		request[i], err = types.SignTx(tx, types.NewEIP155Signer(big.NewInt(89)), unlockedKey.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}
	}
	nonce = nonce + uint64(len(request))
}
