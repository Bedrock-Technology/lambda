package core

import (
	"github.com/ethereum/go-ethereum/common"
)

func HexToAddress(address string) string {
	return common.HexToAddress(address).Hex()
}
