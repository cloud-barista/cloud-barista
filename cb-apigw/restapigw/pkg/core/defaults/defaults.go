// Package defaults - Golang Structure에 Defaults 값 설정 기능 제공 패키지
package defaults

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"
)

// ===== [ Constants and Variables ] =====

const (
	fieldTagName = "default"
)

var (
	errInvalidType = errors.New("not a structure pointer")
)

// ===== [ Types ] =====
// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// shouldInitializeField - 지정한 필드에 태그 정보로 기본 값 설정이 가능한지 검증
func shouldInitializeField(fv reflect.Value, tag string) bool {
	switch fv.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		if !fv.IsNil() && fv.Elem().Kind() == reflect.Struct {
			return true
		}
	case reflect.Slice:
		return fv.Len() > 0 || tag != ""
	}

	return tag != ""
}

// setField - 지정한 필드에 지정한 기본 값 설정
func setField(fv reflect.Value, defaultValue string) error {
	if !fv.CanSet() {
		return nil
	}

	if !shouldInitializeField(fv, defaultValue) {
		return nil
	}

	if isInitialValue(fv) {
		switch fv.Kind() {
		case reflect.Bool:
			if val, err := strconv.ParseBool(defaultValue); nil == err {
				fv.Set(reflect.ValueOf(val).Convert(fv.Type()))
			}
		case reflect.Int:
			if val, err := strconv.ParseInt(defaultValue, 0, strconv.IntSize); nil == err {
				fv.Set(reflect.ValueOf(int(val)).Convert(fv.Type()))
			}
		case reflect.Int8:
			if val, err := strconv.ParseInt(defaultValue, 0, 8); nil == err {
				fv.Set(reflect.ValueOf(int8(val)).Convert(fv.Type()))
			}
		case reflect.Int16:
			if val, err := strconv.ParseInt(defaultValue, 0, 16); nil == err {
				fv.Set(reflect.ValueOf(int16(val)).Convert(fv.Type()))
			}
		case reflect.Int32:
			if val, err := strconv.ParseInt(defaultValue, 0, 32); nil == err {
				fv.Set(reflect.ValueOf(int32(val)).Convert(fv.Type()))
			}
		case reflect.Int64:
			if val, err := time.ParseDuration(defaultValue); nil == err {
				fv.Set(reflect.ValueOf(val).Convert(fv.Type()))
			} else if val, err := strconv.ParseInt(defaultValue, 0, 64); nil == err {
				fv.Set(reflect.ValueOf(val).Convert(fv.Type()))
			}
		case reflect.Uint:
			if val, err := strconv.ParseUint(defaultValue, 0, strconv.IntSize); nil == err {
				fv.Set(reflect.ValueOf(uint(val)).Convert(fv.Type()))
			}
		case reflect.Uint8:
			if val, err := strconv.ParseUint(defaultValue, 0, 8); nil == err {
				fv.Set(reflect.ValueOf(uint8(val)).Convert(fv.Type()))
			}
		case reflect.Uint16:
			if val, err := strconv.ParseUint(defaultValue, 0, 16); nil == err {
				fv.Set(reflect.ValueOf(uint16(val)).Convert(fv.Type()))
			}
		case reflect.Uint32:
			if val, err := strconv.ParseUint(defaultValue, 0, 32); nil == err {
				fv.Set(reflect.ValueOf(uint32(val)).Convert(fv.Type()))
			}
		case reflect.Uint64:
			if val, err := strconv.ParseUint(defaultValue, 0, 64); nil == err {
				fv.Set(reflect.ValueOf(val).Convert(fv.Type()))
			}
		case reflect.Uintptr:
			if val, err := strconv.ParseUint(defaultValue, 0, strconv.IntSize); nil == err {
				fv.Set(reflect.ValueOf(uintptr(val)).Convert(fv.Type()))
			}
		case reflect.Float32:
			if val, err := strconv.ParseFloat(defaultValue, 32); nil == err {
				fv.Set(reflect.ValueOf(float32(val)).Convert(fv.Type()))
			}
		case reflect.Float64:
			if val, err := strconv.ParseFloat(defaultValue, 64); nil == err {
				fv.Set(reflect.ValueOf(val).Convert(fv.Type()))
			}
		case reflect.String:
			fv.Set(reflect.ValueOf(defaultValue).Convert(fv.Type()))

		case reflect.Slice:
			ref := reflect.New(fv.Type())
			ref.Elem().Set(reflect.MakeSlice(fv.Type(), 0, 0))
			if defaultValue != "" && defaultValue != "[]" {
				if err := json.Unmarshal([]byte(defaultValue), ref.Interface()); nil != err {
					return err
				}
			}
			fv.Set(ref.Elem().Convert(fv.Type()))
		case reflect.Map:
			ref := reflect.New(fv.Type())
			ref.Elem().Set(reflect.MakeMap(fv.Type()))
			if defaultValue != "" && defaultValue != "{}" {
				if err := json.Unmarshal([]byte(defaultValue), ref.Interface()); nil != err {
					return err
				}
			}
			fv.Set(ref.Elem().Convert(fv.Type()))
		case reflect.Struct:
			if defaultValue != "" && defaultValue != "{}" {
				if err := json.Unmarshal([]byte(defaultValue), fv.Addr().Interface()); nil != err {
					return err
				}
			}
		case reflect.Ptr:
			fv.Set(reflect.New(fv.Type().Elem()))
		}
	}

	switch fv.Kind() {
	case reflect.Ptr:
		setField(fv.Elem(), defaultValue)
		callSetter(fv.Interface())
	case reflect.Struct:
		if err := Set(fv.Addr().Interface()); nil != err {
			return err
		}
	case reflect.Slice:
		for j := 0; j < fv.Len(); j++ {
			if err := setField(fv.Index(j), defaultValue); nil != err {
				return err
			}
		}
	}

	return nil
}

// isInitialValue - 지정한 필드 값이 초기 값(형식에 따른 Zero Value)인지 검증
func isInitialValue(fv reflect.Value) bool {
	// 필드 형식의 Zero Value 검증
	return reflect.DeepEqual(reflect.Zero(fv.Type()).Interface(), fv.Interface())
}

// ===== [ Public Functions ] =====

// CanUpdate - 지정한 값을 초기값으로 설정할 수있는지 여부 검증
func CanUpdate(v interface{}) bool {
	return isInitialValue(reflect.ValueOf(v))
}

// Set - Pointer로 참조되는 Structure의 Field들에 대한 초기값을 설정한다.
func Set(ptr interface{}) error {
	// Checking structure pointer
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return errInvalidType
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return errInvalidType
	}

	for i := 0; i < t.NumField(); i++ {
		if defaultValue := t.Field(i).Tag.Get(fieldTagName); defaultValue != "-" {
			if err := setField(v.Field(i), defaultValue); nil != err {
				return err
			}
		}
	}

	callSetter(ptr)
	return nil
}
