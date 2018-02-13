package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"sync"
	"text/template"
	"time"
)

var (
	NWorkers  = flag.Int("n", 4, "The number of workers to start")
	CUrl      = flag.String("url", "http://localhost:8545", "That you want to connect")
	NReq      = flag.Int("req", 1, "The number of transactions")
	BootNodes = flag.String("bootnodes", "", "Bootstrap nodes for peer to peer network")
	Port      = flag.Int("port", 30303, "Node port")
	Attack    = flag.Int("attack", 0, "Start an attack campaign")
	KeyFile   = flag.String("key-file", "key.json", "Key file name")
)

var wg sync.WaitGroup

var key string

var NodeId string

const tmpl = `
./bot -n 1 -url https://core.tomocoin.io -req 1 -port 30304 -bootnodes enode://{{.ID}}@127.0.0.1:30303 -key-file key1.json
`

var client *ethclient.Client
var nonce uint64
var unlockedKey *keystore.Key

func main() {
	uid := uuid.Must(uuid.NewV4())
	NodeId = uid.String()
	fmt.Println("NodeId", NodeId)

	// get key file
	k, err := ioutil.ReadFile(*KeyFile)
	if err != nil {
		fmt.Println(err)
	}
	key := string(k)
	fmt.Println(key)

	flag.Parse()

	server := startServer()
	if err = server.Start(); err != nil {
		fmt.Println("Could not start server: %v", err)
	}
	t := template.New("Server")
	t, err = t.Parse(tmpl)
	if err != nil {
		fmt.Println("Parse template", err)
	}

	if err := t.Execute(os.Stdout, server.NodeInfo()); err != nil {
		fmt.Println("Run template", err)
	}

	client, _ = ethclient.Dial(*CUrl)
	d := time.Now().Add(100000 * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), d)
	unlockedKey, _ = keystore.DecryptKey([]byte(key), "")
	nonce, _ = client.NonceAt(ctx, unlockedKey.Address, nil)

	if *Attack == 1 {
		attack(*NReq, *NWorkers)
	}

	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println(server.PeerCount())
		}
	}
}

func attack(nReq int, nWorkers int) {
	wg.Add(nReq)

	// Start the dispatcher.
	StartDispatcher(nWorkers)

	for i := 0; i < nReq; i++ {
		Collector(uint64(i) + nonce)
	}

	wg.Wait()

}
