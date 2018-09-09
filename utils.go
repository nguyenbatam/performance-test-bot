package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/core/types"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"errors"
	"math/big"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	ErrHTTP = errors.New("HTTP_ERROR")
)

func GetBalance(address common.Address, client *http.Client) (*big.Int, error) {
	data := "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"" +
		address.Hex() + "\", \"latest\"],\"id\":1}"
	req, err := http.NewRequest("POST", *Url, bytes.NewReader([]byte(data)))
	if err != nil {
		return nil, ErrHTTP
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	var mapObject = make(map[string]interface{})
	err = json.Unmarshal(respByte, &mapObject)
	if (err != nil) {
		return nil, errors.New(string(respByte))
	}
	error, check := mapObject["error"].(string)
	if (check) {
		return nil, errors.New(error)
	}
	result, check := mapObject["result"].(string)
	if (! check) {
		return nil, errors.New(string(respByte))
	}
	return hexutil.DecodeBig(result)
}

func GetTransactionCount(address common.Address, client *http.Client) (uint64, error) {
	data := "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionCount\",\"params\":[\"" +
		address.Hex() + "\", \"latest\"],\"id\":1}"
	req, err := http.NewRequest("POST", *Url, bytes.NewReader([]byte(data)))
	if err != nil {
		return 0, ErrHTTP
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	var mapObject = make(map[string]interface{})
	err = json.Unmarshal(respByte, &mapObject)
	if (err != nil) {
		return 0, errors.New(string(respByte))
	}
	error, check := mapObject["error"].(string)
	if (check) {
		return 0, errors.New(error)
	}
	result, check := mapObject["result"].(string)
	if (! check) {
		return 0, errors.New(string(respByte))
	}
	count, err := hexutil.DecodeBig(result)
	if err != nil {
		return 0, err
	}
	return count.Uint64(), nil
}

func SendRawTransaction(tx *types.Transaction, client *http.Client) (error) {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	body := "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendRawTransaction\",\"params\":[\"" +
		common.ToHex(data) + "\"],\"id\":1}"
	req, err := http.NewRequest("POST", *Url, bytes.NewReader([]byte(body)))
	if err != nil {
		return ErrHTTP
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respByte, err := ioutil.ReadAll(resp.Body)
	var mapObject = make(map[string]interface{})
	err = json.Unmarshal(respByte, &mapObject)
	if (err != nil) {
		return errors.New(string(respByte))
	}
	error, check := mapObject["error"].(string)
	if (check) {
		return errors.New(error)
	}
	return nil
}
