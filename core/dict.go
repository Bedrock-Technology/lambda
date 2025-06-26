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
