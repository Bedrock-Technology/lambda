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

func makeInjections(ctx *gin.Context) map[string]map[string]any {
	return map[string]map[string]any{
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
			"select": func(db, query string, values ...any) ([]map[string]any, error) {
				return core.TableSelect(dbs[db], query, values...)
			},
			"insert": func(db, table string, obj map[string]any) error {
				return core.TableInsert(dbs[db], table, obj)
			},
		},
		"clickhouse": {
			"description": map[string]any{
				"select": "Selects records from the clickhouse database.",
			},
			"select": func(db, query string, values ...any) ([]map[string]any, error) {
				return core.TableSelect(clickhouseDB[db], query, values...)
			},
		},
		"redis": {
			"hget": func(db, key, field string) (string, error) {
				return core.RedisHGet(redisDB[db], key, field)
			},
			"hset": func(db, key string, values ...any) (int64, error) {
				return core.RedisHSet(redisDB[db], key, values...)
			},
			"hexpire": func(db, key string, duration string, fields ...string) ([]int64, error) {
				return core.RedisHExpire(redisDB[db], key, duration, fields...)
			},
			"hkeys": func(db, key string) ([]string, error) {
				return core.RedisHKeys(redisDB[db], key)
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
				"first_or":           "Returns the first non-empty value from a list of values.",
			},
			"hex_to_address":     core.HexToAddress,
			"hex_to_hash":        core.HexToHash,
			"strings_equal_fold": strings.EqualFold,
			"hash_typed_data":    core.HashTypedData,
			"bech32_address":     core.Bech32Address,
			"decimal_add":        core.DecimalAdd,
			"decimal_sub":        core.DecimalSub,
			"decimal_mul":        core.DecimalMul,
			"decimal_divround":   core.DecimalDivRound,
			"csv_read":           core.CSVRead,
			"first_or":           lo.FirstOr[any],
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
		"shared": {
			"description": map[string]any{
				"dict_get":         "Gets a value from the shared dictionary for the current service.",
				"dict_set":         "Sets a value in the shared dictionary for the current service.",
				"dict_keys":        "Gets all keys from the shared dictionary for the current service.",
				"dict_del":         "Deletes a key from the shared dictionary for the current service.",
				"dict_get_global":  "Gets a value from the global shared dictionary.",
				"dict_set_global":  "Sets a value in the global shared dictionary.",
				"dict_keys_global": "Gets all keys from the global shared dictionary.",
				"dict_del_global":  "Deletes a key from the global shared dictionary.",
			},
			"dict_get": func(key string) (any, bool) {
				serviceName, _ := ctx.Get(serviceNameKey)
				return core.DictGet(serviceName.(string), key)
			},
			"dict_set": func(key string, val any) {
				serviceName, _ := ctx.Get(serviceNameKey)
				core.DictSet(serviceName.(string), key, val)
			},
			"dict_keys": func() []string {
				serviceName, _ := ctx.Get(serviceNameKey)
				return core.DictKeys(serviceName.(string))
			},
			"dict_del": func(key string) {
				serviceName, _ := ctx.Get(serviceNameKey)
				core.DictDel(serviceName.(string), key)
			},
			"dict_get_global":  core.DictGetGlobal,
			"dict_set_global":  core.DictSetGlobal,
			"dict_keys_global": core.DictKeysGlobal,
			"dict_del_global":  core.DictDelGlobal,
		},
	}
}

func injectorFor(vm *goja.Runtime, ctx *gin.Context) *goja.Object {
	mp := make(map[string]any)

	injections := makeInjections(ctx)
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
	Form    map[string]string   `json:"form"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func makeVarsObj(vm *goja.Runtime, ctx *gin.Context) *goja.Object {
	cfgLock.Lock()
	vars, varsDesc := cfg.Vars, cfg.VarsDesc
	varsDesc["req"] = "The request object."
	cfgLock.Unlock()

	varsObj := mapToObject(vm, lo.MapEntries(vars, func(k1, v1 string) (string, any) {
		return k1, v1
	}))

	formMap := make(map[string]string)
	if ctxForm, err := ctx.MultipartForm(); err == nil {
		for k, v := range ctxForm.File {
			if len(v) > 0 {
				f := lo.Must(v[0].Open())
				fileContent := lo.Must(io.ReadAll(f))
				formMap[k] = string(fileContent)
			}
		}

		for k, v := range ctxForm.Value {
			if len(v) > 0 {
				formMap[k] = v[0]
			}
		}
	}

	body := ""
	if ctx.Request.Body != nil {
		body = string(lo.Must(io.ReadAll(ctx.Request.Body)))
	}
	r := rawRequest{
		Method:  ctx.Request.Method,
		Path:    ctx.Request.URL.Path,
		Query:   ctx.Request.URL.Query(),
		Headers: ctx.Request.Header,
		Form:    formMap,
		Body:    body,
	}
	varsObj.Set("req", r)
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
