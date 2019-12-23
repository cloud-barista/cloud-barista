package util

import (
	"bytes"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
)

func StructToMap(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		values.Set(typ.Field(i).Name, v)
	}
	return
}

func ToMap(val interface{}) (map[string]interface{}, error) {

	// Convert struct to bytes
	bytes := new(bytes.Buffer)
	if err := json.NewEncoder(bytes).Encode(val); err != nil {
		return nil, err
	}

	// Convert bytes to map
	byteData := bytes.Bytes()
	resultMap := map[string]interface{}{}
	if err := json.Unmarshal(byteData, &resultMap); err != nil {
		return nil, err
	}

	return resultMap, nil
}

func GetFields(val reflect.Value) []string {
	var fieldArr []string
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldArr = append(fieldArr, t.Field(i).Tag.Get("json"))
	}
	return fieldArr
}
