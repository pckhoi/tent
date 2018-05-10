package postgres

import (
	"sync"
)

var settingsMutex = sync.Mutex{}
var settings = make(map[string][]interface{})

func updateSettings(key string, vals []interface{}) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	settings[key] = vals
}
