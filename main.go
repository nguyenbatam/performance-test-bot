package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	NReq     = flag.Int("req", 1, "The number of transactions")
	Urls     = flag.String("urls", "http://localhost:8545,http://localhost:8546", "That you want to connect")
	KeyFile  = flag.String("key-file", "key.json", "Key file name")
	Password = flag.String("password", "", "Keyfile password")
)
var wg sync.WaitGroup
var mainClient *ethclient.Client
var mainNonce uint64
var unlockedKey *keystore.Key
var nAccount int
var ks *keystore.KeyStore
var _10TOMO = big.NewInt(1).Mul(big.NewInt(int64(math.Pow10(9))), big.NewInt(int64(math.Pow10(10))))
var gasLimit = uint64(21000)
var LockSendMoneyToBot sync.RWMutex
var urls []string
var unlockedBotKeys []*keystore.Key

func main() {
	flag.Parse()
	key_file := *KeyFile
	if _, err := os.Stat(key_file); err != nil {
		fmt.Println(err)
	}
	// get key file
	data, err := ioutil.ReadFile(key_file)
	if err != nil {
		fmt.Println(err)
	}
	unlockedKey, err = keystore.DecryptKey(data, *Password)
	if err != nil {
		fmt.Println(err)
	}

	urls = strings.Split(*Urls, ",")
	mainClient, err = ethclient.Dial(urls[0])
	if err != nil {
		log.Fatal(err)
	}
	mainNonce, _ = mainClient.PendingNonceAt(context.Background(), unlockedKey.Address)

	fmt.Println("read account bot mainNonce ", mainNonce)
	nAccount = len(urls)
	// Create a new account with the specified encryption passphrase
	fmt.Println("create ", nAccount, "new account and create transaction send money for bots", _10TOMO)
	ks = keystore.NewKeyStore("", keystore.StandardScryptN, keystore.StandardScryptP)
	for i := 0; i < nAccount; i++ {
		newAcc, err1 := ks.NewAccount("")
		if err1 != nil {
			log.Fatal(err1)
		}
		ks.Unlock(newAcc, "")
		jsonByte, _ := ks.Export(newAcc, "", "")
		unlockedBot, _ := keystore.DecryptKey(jsonByte, "")
		unlockedBotKeys = append(unlockedBotKeys, unlockedBot)
		sendMoneyToBot(newAcc.Address)

	}

	fmt.Println("wait done send money for bot")
	done := false
	for !done {
		time.Sleep(5 * time.Second)
		done = true
		for i := 0; i < nAccount; i++ {
			result, err1 := mainClient.BalanceAt(context.Background(), unlockedBotKeys[i].Address, nil)
			if err1 != nil || result.Uint64() < gasLimit {
				fmt.Println("stil waiting send money to ", unlockedBotKeys[i].Address.Hex(), err1, result)
				done = false
				break
			}
		}
	}

	fmt.Println("Start run ", nAccount, "bot")
	wg.Add(nAccount)
	for i := 0; i < nAccount; i++ {
		go attack(*NReq, urls[i], unlockedBotKeys[i])
	}
	wg.Wait()
}

func attack(request int, url string, account *keystore.Key) {
	fmt.Println("Start run ", request, " request  with account ", account.Address.Hex(), "bot with", url)
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	nonce, _ := client.PendingNonceAt(context.Background(), account.Address)
	// Start the dispatcher.
	for {
		start := time.Now().UnixNano() / int64(time.Millisecond)
		balance, _ := client.BalanceAt(context.Background(), account.Address, nil)
		if balance.Uint64() < gasLimit*10 {
			fmt.Println("Request  money from main account to ", account.Address.Hex(), "balance", balance)
			sendMoneyToBot(account.Address)
			time.Sleep(5 * time.Second)
		}
		for i := 0; i < request; i++ {
			tx := types.NewTransaction(nonce, account.Address, big.NewInt(1), gasLimit, big.NewInt(1), nil)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(89)), account.PrivateKey)
			if err != nil {
				log.Fatal(err)
			}
			err = client.SendTransaction(context.Background(), signedTx)
			if (err != nil) {
				fmt.Println(err, signedTx)
			}
			nonce++
		}
		end := time.Now().UnixNano() / int64(time.Millisecond)
		if (end-start < 1000) {
			sleep := int(1000 + start - end)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		} else {
			fmt.Println("Done a round with time = ", end-start, "url", url)
		}
	}
	wg.Done()
}

func sendMoneyToBot(address common.Address) {
	fmt.Println("start create transaction for bot ", address.Hex())
	LockSendMoneyToBot.Lock()
	defer LockSendMoneyToBot.Unlock()
	tx := types.NewTransaction(mainNonce, address, _10TOMO, gasLimit, big.NewInt(1), nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(89)), unlockedKey.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = mainClient.SendTransaction(context.Background(), signedTx)
	if (err != nil) {
		log.Fatal(err, signedTx)
	}
	fmt.Println("done send transaction for bot ", address.Hex(), "nonce", mainNonce)
	mainNonce++
}
