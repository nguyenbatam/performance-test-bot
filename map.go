package main

import (
	"fmt"
	"github.com/hashicorp/golang-lru"
	"time"
)
func main() {
	knownTxs, _ := lru.New(100000)
	for i:=0;i<100*1000;i++{
		fmt.Println(i)
		fmt.Println(knownTxs.ContainsOrAdd(i,nil))
		fmt.Println(knownTxs.Contains(i))
	}
	time.Sleep(1 *time.Hour)
}
