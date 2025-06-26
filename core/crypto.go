package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/samber/lo"
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

type MerkleArgs struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type MerkleTree struct {
	Root    string              `json:"root"`
	Proof   map[string][]string `json:"proof"`
	Amounts map[string]string   `json:"amounts"`
}

func TryMerkle(name string) (*MerkleTree, bool) {
	key := "merkle:" + name

	f, ok := sharedDict.Load(key)
	if !ok {
		return nil, false
	}

	tree, err := f.(func() (*MerkleTree, error))()
	if err != nil {
		return nil, false
	}

	return tree, true
}

func Merkle(name string, args []MerkleArgs) (*MerkleTree, error) {
	key := "merkle:" + name

	f, _ := sharedDict.LoadOrStore(key, sync.OnceValues(func() (*MerkleTree, error) {
		return merkle(args)
	}))

	return f.(func() (*MerkleTree, error))()
}

func merkle(args []MerkleArgs) (*MerkleTree, error) {
	leaves := lo.Map(args, func(a MerkleArgs, _ int) []any {
		return []any{
			any(smt.SolAddress(a.Address)),
			any(smt.SolNumber(a.Amount)),
		}
	})

	leafEncodings := []string{
		smt.SOL_ADDRESS,
		smt.SOL_UINT256,
	}

	tree, err := smt.Of(leaves, leafEncodings)
	if err != nil {
		return nil, err
	}

	proofMap := make(map[string][]string)
	for _, leaf := range leaves {
		proof, err := tree.GetProof(leaf)
		if err != nil {
			return nil, err
		}

		ok, err := tree.Verify(proof, leaf)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, fmt.Errorf("invalid proof for leaf %v", leaf)
		}

		proofMap[leaf[0].(common.Address).Hex()] = lo.Map(proof, func(p []byte, _ int) string {
			return hexutil.Encode(p)
		})
	}

	return &MerkleTree{
		Root:  hexutil.Encode(tree.GetRoot()),
		Proof: proofMap,
		Amounts: lo.SliceToMap(args, func(a MerkleArgs) (string, string) {
			return a.Address, a.Amount
		}),
	}, nil
}
