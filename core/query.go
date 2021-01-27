package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/schoeppi5/libts"
)

// Shared codebase for queries

// QueryError is returned when the query answered with an error
type QueryError struct {
	ID           int    `mapstructure:"id" json:"code"`
	Message      string `mapstructure:"msg" json:"message"`
	ExtraMessage string `json:"extra_message"` // only set on webquery
}

func (qe QueryError) Error() string {
	s := fmt.Sprintf("Query error(%d): %s", qe.ID, qe.Message)
	if qe.ExtraMessage != "" {
		s += fmt.Sprintf(" (Extra message: %s)", qe.ExtraMessage)
	}
	return s
}

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

// Handling of persistent connections

// KeepAlive sends " \n" every t to conn
func KeepAlive(conn io.Writer, t time.Duration) {
	ticker := time.NewTicker(t)
	for {
		<-ticker.C
		conn.Write([]byte(" \n"))
	}
}

// ReadHeader slurps the header from the io.Reader
func ReadHeader(r <-chan []byte) error {
	// header
	header, open := <-r
	if !open {
		return errors.New("unable to read header: connection closed")
	}
	if string(header) != "TS3" {
		return errors.New("wrong header")
	}
	// banner
	_, open = <-r
	if !open {
		return errors.New("unable to read banner: connection closed")
	}
	return nil
}

// Run writes r to in and reads the first line from out and checks, if it is an error
// r must be terminated by \n
func Run(in <-chan []byte, out io.Writer, r []byte) ([]byte, error) {
	_, err := out.Write(r)
	if err != nil {
		return nil, err
	}
	var data []byte
	for {
		d, open := <-in
		if !open {
			return nil, errors.New("unable to read response: connection closed")
		}
		if err = isError(d); err != nil {
			if e, ok := err.(QueryError); ok {
				if e.ID == 0 {
					return data, nil
				}
				return nil, err
			}
			return nil, err
		}
		data = d
	}
}

// Demultiplexer seperates the notifications and everything else
// Stops when m is closed or it encounters an error while reading from m
func Demultiplexer(m io.Reader, out chan<- []byte, notify chan<- []byte) {
	reader := bufio.NewReader(m)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			close(out)
			close(notify)
			return
		}
		data = bytes.TrimRight(bytes.TrimLeft(data, "\r"), "\n") // normalize data
		if bytes.Index(data, []byte("notify")) == 0 {
			notify <- data
			continue
		}
		out <- data
	}
}

// Helpers for converting responses

// IsError returns the error if any
func isError(r []byte) error {
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
	items := bytes.Split(r, []byte("|"))
	for i := range items {
		list = append(list, responseToMap(items[i]))
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
