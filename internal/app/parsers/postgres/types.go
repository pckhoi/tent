package postgres

import (
	"fmt"
	"github.com/pckhoi/tent/internal/app/storage"
	"reflect"
	"strconv"
)

type String string

type Identifier string

func updateSettings(key string, val interface{}) storage.DataRow {
	var setting interface{}
	switch v := val.(type) {
	case int64:
		setting = strconv.FormatInt(v, 10)
	case bool:
		setting = strconv.FormatBool(v)
	default:
		setting = v
	}
	return storage.DataRow{
		TableName: "postgres_settings",
		ID:        key,
		Content: map[string]string{
			"name":    key,
			"setting": fmt.Sprintf("%s", reflect.ValueOf(setting)),
			"type":    fmt.Sprintf("%s", reflect.TypeOf(val)),
		},
	}
}
