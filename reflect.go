package mp

import (
	"fmt"
	"reflect"
	"strconv"
)

type Balance struct {
	InitCash   float64 `json:"initCash"`
	BoughtCash float64 `json:"boughtCash"`
	SoldCash   float64 `json:"soldCash"`

	BidWorth float64 `json:"bidWorth"`
	Worth    float64 `json:"worth"`
	AskWorth float64 `json:"askWorth"`
}

/*func IsBasicValue(v reflect.Value) bool {
	if v.IsValid() {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return true
		}
	}
	return false
}*/

func IsBasic(val interface{}) bool {
	switch val.(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16,
		uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func BasicString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

func StructToAttrs(val interface{}, keys ...string) map[string]string {
	v := reflect.Indirect(reflect.ValueOf(val))

	if v.Kind() != reflect.Struct { // including !v.IsValid() where v.Kind() == reflect.Invalid
		return nil
	}

	ret := map[string]string{}
	if len(keys) == 0 { // NOTE: not including inline structs
		for i, n := 0, v.NumField(); i < n; i++ {
			name := v.Type().Field(i).Name
			if f := v.Field(i); f.IsValid() {
				if fv := f.Interface(); IsBasic(fv) {
					ret[name] = BasicString(fv)
				}
			}
		}
	} else {
		for _, key := range keys {
			if f := v.FieldByName(key); f.IsValid() {
				if fv := f.Interface(); IsBasic(fv) {
					ret[key] = BasicString(fv)
				}
			}
		}
	}
	return ret
}
