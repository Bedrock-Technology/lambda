package main

import (
	"io"
	"log/slog"
	"strings"

	"github.com/Bedrock-Technology/lambda/core"

	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var (
	injections = map[string]map[string]any{
		"crypto": {
			"description": map[string]any{
				"keccak256": "Calculates the Keccak-256 hash of the input.",
				"ecrecover": "Recovers the address associated with the public key from a message and a signature.",
				"merkle":    "Generates a Merkle tree and provides its root hash.",
			},
			"keccak256": core.Keccak256,
			"ecrecover": core.Ecrecover,
			"merkle":    core.Merkle,
		},
		"net": {
			"description": map[string]any{
				"fetch": "Performs an HTTP request.",
			},
			"fetch": core.Fetch,
		},
		"db": {
			"description": map[string]any{
				"insert": "Inserts a record into the postgres database.",
				"select": "Selects records from the postgres database.",
			},
			"select": func(db, query string) ([]map[string]any, error) {
				return core.TableSelect(dbs[db], query)
			},
			"insert": func(db, table string, obj map[string]any) error {
				return core.TableInsert(dbs[db], table, obj)
			},
		},
		"clickhouse": {
			"description": map[string]any{
				"select": "Selects records from the clickhouse database.",
			},
			"select": func(db, query string) ([]map[string]any, error) {
				return core.TableSelect(clickhouseDB[db], query)
			},
		},
		"utils": {
			"description": map[string]any{
				"hex_to_address":     "Converts a hexadecimal string to an Ethereum address.",
				"hex_to_hash":        "Converts a hexadecimal string to an Ethereum hash.",
				"strings_equal_fold": "Compares two strings case-insensitively.",
				"hash_typed_data":    "Hashes a typed data object.",
				"bech32_address":     "Converts a Bech32 address to its prefix and encoded form.",
				"decimal_add":        "Adds two decimal strings.",
				"decimal_sub":        "Subtracts two decimal strings.",
				"decimal_mul":        "Multiplies two decimal strings.",
				"decimal_divround":   "Divides two decimal strings and rounds the result.",
				"csv_read":           "Reads a CSV file and returns its contents.",
			},
			"hex_to_address":     core.HexToAddress,
			"strings_equal_fold": strings.EqualFold,
			"hash_typed_data":    core.HashTypedData,
			"bech32_address":     core.Bech32Address,
			"decimal_add":        core.DecimalAdd,
			"decimal_sub":        core.DecimalSub,
			"decimal_mul":        core.DecimalMul,
			"decimal_divround":   core.DecimalDivRound,
			"csv_read":           core.CSVRead,
		},
		"slog": {
			"description": map[string]any{
				"debug": "Logs a debug message.",
				"info":  "Logs an info message.",
				"warn":  "Logs a warning message.",
				"error": "Logs an error message.",
			},
			"debug": slog.Debug,
			"info":  slog.Info,
			"warn":  slog.Warn,
			"error": slog.Error,
		},
	}
)

func injectorFor(vm *goja.Runtime, ctx *gin.Context) *goja.Object {
	mp := make(map[string]any)
	for k, v := range injections {
		mp[k] = mapToObject(vm, v)
	}
	mp["vars"] = makeVarsObj(vm, ctx)

	return mapToObject(vm, mp)
}

type rawRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func makeVarsObj(vm *goja.Runtime, ctx *gin.Context) *goja.Object {
	cfgLock.RLock()
	vars, varsDesc := cfg.Vars, cfg.VarsDesc
	cfgLock.RUnlock()

	varsObj := mapToObject(vm, lo.MapEntries(vars, func(k1, v1 string) (string, any) {
		return k1, v1
	}))

	body := ""
	if ctx.Request.Body != nil {
		body = string(lo.Must(io.ReadAll(ctx.Request.Body)))
	}
	r := rawRequest{
		Method:  ctx.Request.Method,
		Path:    ctx.Request.URL.Path,
		Query:   ctx.Request.URL.Query(),
		Headers: ctx.Request.Header,
		Body:    body,
	}
	varsObj.Set("req", r)

	varsDesc["req"] = "The request object."
	varsObj.Set("description", varsDesc)

	return varsObj
}

func mapToObject(vm *goja.Runtime, mp map[string]any) *goja.Object {
	obj := vm.NewObject()
	for k, v := range mp {
		obj.Set(k, v)
	}
	return obj
}
