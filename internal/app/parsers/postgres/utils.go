package postgres

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func interfaceToString(val interface{}) string {
	return fmt.Sprintf("%s", reflect.ValueOf(val))
}

func interfaceToMap(val interface{}) map[string]string {
	ref := reflect.ValueOf(val)
	result := map[string]string{}
	for _, key := range ref.MapKeys() {
		result[key.String()] = ref.MapIndex(key).String()
	}
	return result
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

func extractBytesSlice(src []interface{}) []byte {
	bytes := []byte{}
	for _, val := range src {
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case []byte:
			bytes = append(bytes, v...)
		case []interface{}:
			bytes = append(bytes, extractBytesSlice(v)...)
		case string:
			bytes = append(bytes, []byte(v)...)
		}
	}
	return bytes
}

func complexString(src interface{}) string {
	valsSlice := toIfaceSlice(src)
	return string(extractBytesSlice(valsSlice))
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func extractList(list interface{}, firstElIndex int, restIndex int, restElIndex int) []interface{} {
	vals := []interface{}{}
	valsSlice := toIfaceSlice(list)
	if firstElIndex != -1 {
		vals = append(vals, valsSlice[firstElIndex])
	}
	var restSlice []interface{}
	if firstElIndex == -1 {
		restSlice = valsSlice
	} else {
		restSlice = toIfaceSlice(valsSlice[restIndex])
	}
	for _, v := range restSlice {
		vSlice := toIfaceSlice(v)
		vals = append(vals, vSlice[restElIndex])
	}
	return vals
}

func toByteSlice(v interface{}) []byte {
	valsSl := toIfaceSlice(v)
	var result []byte
	for _, val := range valsSl {
		result = append(result, val.([]byte)[0])
	}
	return result
}

func serializeStringSlice(slice []String) string {
	formattedStrings := []string{}
	for _, s := range slice {
		escapedString := strings.Replace(interfaceToString(s), "'", "''", -1)
		formattedStrings = append(formattedStrings, fmt.Sprintf("'%s'", escapedString))
	}
	return strings.Join(formattedStrings, ",")
}
