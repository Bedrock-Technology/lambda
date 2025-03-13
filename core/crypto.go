package core

import (
	"encoding/json"
	"log/slog"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
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

func HashTypedData(typedDataJson string) (string, error) {
	data := apitypes.TypedData{}
	if err := json.Unmarshal([]byte(typedDataJson), &data); err != nil {
		return "", err
	}

	slog.Debug("HashTypedData", "data", data)

	hash, _, err := apitypes.TypedDataAndHash(data)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(hash), nil
}
