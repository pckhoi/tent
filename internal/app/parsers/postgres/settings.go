package postgres

import (
	"fmt"
	"sync"
)

var settingsMutex = sync.Mutex{}
var settings = make(map[string][]interface{})

func updateSettings(key string, vals []interface{}) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	if val, ok := settings[key]; ok {
		settings[key] = append(val, vals...)
	} else {
		settings[key] = vals
	}
}

func setSettings(key string, vals []interface{}) {
	settingsMutex.Lock()
	defer settingsMutex.Unlock()
	settings[key] = vals
}

func typeExists(typename string) error {
	types, ok := settings["custom_types"]
	if ok {
		for _, t := range types {
			if t == typename {
				return nil
			}
		}
	}
	return fmt.Errorf("Type %s is not defined", typename)
}
