package core

import (
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
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
