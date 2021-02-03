package communication

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/schoeppi5/libts"
)

// UnmarshalResponse attempts to parse body to value
// If value is a single struct, only the first element of body is unmarsheld
// If value is a slice of structs, UnmarshalResponse will append to that slice
// If value is an array of structs, UnmarshalResponse will try to set the indices of the array
func UnmarshalResponse(body []map[string]interface{}, value interface{}) error {
	kind := reflect.Indirect(reflect.ValueOf(value)).Kind()
	if kind == reflect.Struct { // one item expected
		err := Decode(body[0], value)
		if err != nil {
			return err
		}
	} else if kind == reflect.Slice { // slice of items expected
		inter := getTypeOfSlice(value)
		v := reflect.ValueOf(value).Elem()
		for i := range body {
			err := Decode(body[i], &inter)
			if err != nil {
				return err
			}
			v.Set(reflect.Append(v, reflect.ValueOf(inter)))
		}
	} else if reflect.Indirect(reflect.ValueOf(value)).Kind() == reflect.Array { // array of items expected
		inter := getTypeOfSlice(value)
		v := reflect.ValueOf(value).Elem()
		for i := range body {
			if i > reflect.Indirect(reflect.ValueOf(value)).Len()-1 { // reached end of expected output
				break
			}
			err := Decode(body[i], &inter)
			if err != nil {
				return err
			}
			v.Index(i).Set(reflect.ValueOf(inter))
		}
	} else {
		return fmt.Errorf("unsupported type %s. Expected type struct, slice or array", reflect.Indirect(reflect.ValueOf(value)).Kind())
	}
	return nil
}

// Decode creates a new decoder and decodes m to v
func Decode(m map[string]interface{}, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return errors.New("expected pointer to value, not value")
	}
	decodeConfig := &mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.TextUnmarshallerHookFunc(),
		Metadata:         nil,
		Result:           v,
		WeaklyTypedInput: true,
	}
	decoder, _ := mapstructure.NewDecoder(decodeConfig)
	err := decoder.Decode(m)
	if err != nil {
		return err
	}
	return nil
}

// getTypeOfSlice returns the zero value for the type of a slice or arry
// e.g.: []struct{} -> zero value of struct{}
func getTypeOfSlice(s interface{}) interface{} {
	return reflect.Zero(reflect.ValueOf(s).Elem().Type().Elem()).Interface()
}

// IsError returns the error if any
func IsError(r []byte) error {
	parts := bytes.SplitN(r, []byte(" "), 2)
	if bytes.Compare(parts[0], []byte("error")) == 0 {
		e := &QueryError{}
		err := UnmarshalResponse(ConvertResponse(parts[1]), e)
		if err != nil {
			return err
		}
		return *e
	}
	return nil
}

// ConvertResponse parses the serverquery response to a map
func ConvertResponse(r []byte) []map[string]interface{} {
	var list []map[string]interface{}
	r = bytes.Trim(bytes.TrimSpace(r), "|")
	items := bytes.Split(r, []byte("|"))
	for j := range items {
		list = append(list, responseToMap(items[j]))
	}
	return list
}

func responseToMap(r []byte) map[string]interface{} {
	m := make(map[string]interface{})
	pairs := bytes.Split(r, []byte(" "))
	for i := range pairs {
		kV := bytes.SplitN(pairs[i], []byte("="), 2)
		key := string(kV[0])
		if len(kV) != 2 {
			m[key] = ""
			continue
		}
		if strings.Contains(key, "client_default_channel") { // I do not know why. The client_default_channel is wrongly encoded. I don't know why
			kV[1] = []byte(libts.QueryDecoder.Replace(string(kV[1])))
		}
		m[key] = libts.QueryDecoder.Replace(string(kV[1]))
	}
	return m
}
