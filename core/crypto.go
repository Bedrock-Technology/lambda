package core

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func Keccak256(data string) string {
	return hexutil.Encode(crypto.Keccak256([]byte(data)))
}

func Ecrecover(hash, sig string) (string, error) {
	hashBytes, err := hexutil.Decode(hash)
	if err != nil {
		return "", err
	}

	sigBytes, err := hexutil.Decode(sig)
	if err != nil {
		return "", err
	}

	if sigBytes[crypto.RecoveryIDOffset] == 27 || sigBytes[crypto.RecoveryIDOffset] == 28 {
		sigBytes[crypto.RecoveryIDOffset] -= 27
	}

	pub, err := crypto.SigToPub(hashBytes, sigBytes)
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*pub).Hex(), nil
}
