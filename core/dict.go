package core

import (
	"fmt"
	"sync"
)

var (
	sharedDict = new(sync.Map)
)

func sharedDictKey(service, key string) string {
	return fmt.Sprintf("%s:%s", service, key)
}

func DictGet(service, key string) (any, bool) {
	return sharedDict.Load(sharedDictKey(service, key))
}

func DictGetGlobal(key string) (any, bool) {
	return sharedDict.Load(key)
}

func DictSet(service, key string, val any) {
	sharedDict.Store(sharedDictKey(service, key), val)
	sharedDict.Range(func(key, value any) bool {
		fmt.Printf("[DEBUG] dump sharedDict, key: %s, value: %T\n", key, value)
		return true
	})
}

func DictSetGlobal(key string, val any) {
	sharedDict.Store(key, val)
	sharedDict.Range(func(key, value any) bool {
		fmt.Printf("[DEBUG] dump sharedDict, key: %s, value: %T\n", key, value)
		return true
	})
}

func DictKeys(service string) []string {
	keys := make([]string, 0)
	sharedDict.Range(func(key, value any) bool {
		if k, ok := key.(string); ok && k[:len(service)+1] == service+":" {
			keys = append(keys, k[len(service)+1:])
		}
		return true
	})
	return keys
}

func DictKeysGlobal() []string {
	keys := make([]string, 0)
	sharedDict.Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			keys = append(keys, k)
		}
		return true
	})
	return keys
}
