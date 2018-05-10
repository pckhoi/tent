package postgres

import (
	"fmt"
	"reflect"
	"strconv"
)

func interfaceToString(val interface{}) string {
	return fmt.Sprintf("%s", reflect.ValueOf(val))
}

func getValueAndTypeAsStrings(val interface{}) (string, string) {
	var value interface{}
	switch v := val.(type) {
	case int64:
		value = strconv.FormatInt(v, 10)
	case bool:
		value = strconv.FormatBool(v)
	default:
		value = v
	}
	return interfaceToString(value), fmt.Sprintf("%s", reflect.TypeOf(val))
}
