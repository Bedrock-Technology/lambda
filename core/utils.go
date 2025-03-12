package core

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func ValidateAddress(address string) bool {
	return strings.EqualFold(address, common.HexToAddress(address).Hex())
}
