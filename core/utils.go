package core

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func HexToAddress(address string) string {
	return common.HexToAddress(address).Hex()
}

func HexToHash(input string) string {
	return common.HexToHash(input).Hex()
}

func Bech32Address(address string) (string, string, error) {
	prefix, data, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", "", err
	}

	addr, err := bech32.ConvertAndEncode(prefix, data)
	if err != nil {
		return "", "", err
	}

	return prefix, addr, err
}

func DecimalAdd(a, b string) (string, error) {
	decA, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}

	decB, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}

	decA = decA.Add(decB)
	return decA.String(), nil
}

func DecimalSub(a, b string) (string, error) {
	decA, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}

	decB, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}

	decA = decA.Sub(decB)
	return decA.String(), nil
}

func DecimalMul(a, b string) (string, error) {
	decA, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}

	decB, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}

	decA = decA.Mul(decB)
	return decA.String(), nil
}

func DecimalDivRound(a, b string, precision int32) (string, error) {
	decA, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}

	decB, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}

	decA = decA.DivRound(decB, precision)
	return decA.String(), nil
}

func CSVRead(filename string) ([][]string, error) {
	key := "csvRead:" + filename

	f, _ := sharedDict.LoadOrStore(key, sync.OnceValues(func() ([][]string, error) {
		fmt.Println("Reading CSV file:", filename)
		return csvRead(filename)
	}))

	return f.(func() ([][]string, error))()
}

func csvRead(filename string) ([][]string, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}
