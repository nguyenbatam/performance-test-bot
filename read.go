package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var notified int32

func main() {
	file, err := os.Open("vps.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	host := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Copy") {
			host = host + "http://" + strings.TrimSpace(strings.Replace(line, " Copy", ":8545,", 1))
			fmt.Println(line)
		}
	}
	fmt.Println(host)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	second := time.Now().Unix()
	fmt.Println(second)
	time.Sleep(2 * time.Second)
	second = time.Now().Unix()
	fmt.Println(second)

}
