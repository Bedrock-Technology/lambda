package core

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func HexToAddress(address string) string {
	return common.HexToAddress(address).Hex()
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
