package postgres

import (
	"fmt"
	"reflect"
	"strconv"
)

type Update struct {
	TableName string
	Row       map[string]interface{}
}

type String string

type Identifier string

func updateSettings(key string, val interface{}) Update {
	var setting interface{}
	switch v := val.(type) {
	case int64:
		setting = strconv.FormatInt(v, 10)
	case bool:
		setting = strconv.FormatBool(v)
	default:
		setting = v
	}
	return Update{
		TableName: "postgres_settings",
		Row: map[string]interface{}{
			"name":    key,
			"setting": fmt.Sprintf("%s", reflect.ValueOf(setting)),
			"type":    fmt.Sprintf("%s", reflect.TypeOf(val)),
		},
	}
}
