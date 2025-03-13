package main

import (
	"io"
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
			},
			"keccak256": core.Keccak256,
			"ecrecover": core.Ecrecover,
		},
		"net": {
			"description": map[string]any{
				"fetch": "Performs an HTTP request.",
			},
			"fetch": core.Fetch,
		},
		"db": {
			"description": map[string]any{
				"insert": "Inserts a record into the database.",
				"select": "Selects records from the database.",
			},
			"select": func(query string) ([]map[string]any, error) {
				return core.TableSelect(db, query)
			},
			"insert": func(table string, obj map[string]any) error {
				return core.TableInsert(db, table, obj)
			},
		},
		"utils": {
			"description": map[string]any{
				"hex_to_address":     "Converts a hexadecimal string to an Ethereum address.",
				"strings_equal_fold": "Compares two strings case-insensitively.",
			},
			"hex_to_address":     core.HexToAddress,
			"strings_equal_fold": strings.EqualFold,
		},
	}
)

func injectorFor(vm *goja.Runtime, ctx *gin.Context) *goja.Object {
	mp := make(map[string]any)
	for k, v := range injections {
		if k == "db" && db == nil {
			continue
		}
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
