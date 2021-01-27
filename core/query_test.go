package core_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/schoeppi5/libts"

	"github.com/schoeppi5/libts/core"
)

func TestQueryErrorIsError(t *testing.T) {
	t.Parallel()
	// given
	qe := &core.QueryError{}

	// when
	_, isError := interface{}(qe).(error)

	// then
	if !isError {
		LogTestError(qe, "error{}", t, "QueryError does not implement the error interface")
	}
}

func TestQueryErrorString(t *testing.T) {
	t.Parallel()
	// given
	qe := &core.QueryError{
		ID:      1234,
		Message: "This is a test",
	}
	want := "Query error(1234): This is a test"

	// when
	have := qe.Error()

	// then
	if have != want {
		LogTestError(have, want, t)
	}
}

func TestQueryErrorStringExtraMessage(t *testing.T) {
	t.Parallel()
	// given
	qe := &core.QueryError{
		ID:           1234,
		Message:      "This is a test",
		ExtraMessage: "A very serious test",
	}
	want := "Query error(1234): This is a test (Extra message: A very serious test)"

	// when
	have := qe.Error()

	// then
	if have != want {
		LogTestError(have, want, t)
	}
}

func TestUnmarshalResponseMapToStruct(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1",
		"Tag2": "test2",
	}}
	have := struct {
		Tag1 string
		Tag2 string
	}{}

	want := struct {
		Tag1 string
		Tag2 string
	}{
		Tag1: "test1",
		Tag2: "test2",
	}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("Test failed with error: %+v", err))
	}
	if have != want {
		LogTestError(have, want, t)
	}
}

func TestUnmarshalResponseMapToStructError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1",
		"Tag2": "test2",
	}}
	have := struct {
		Tag1 []bool
		Tag2 string
	}{}
	want := "1 error(s) decoding:\n\n* cannot parse 'Tag1[0]' as bool: strconv.ParseBool: parsing \"test1\": invalid syntax"

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected unmarshal error for type []bool")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestUnmarshalResponseMapToArray(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := [2]struct {
		Tag1 string
		Tag2 string
	}{}
	want := [2]struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "test1_1",
		Tag2: "test1_2",
	}, {
		Tag1: "test2_1",
		Tag2: "test2_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToShorterArray(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := [1]struct {
		Tag1 string
		Tag2 string
	}{}
	want := [1]struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "test1_1",
		Tag2: "test1_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToArrayError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := [2]struct {
		Tag1 []bool
		Tag2 string
	}{}
	want := "1 error(s) decoding:\n\n* cannot parse 'Tag1[0]' as bool: strconv.ParseBool: parsing \"test1_1\": invalid syntax"

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected error unmarshaling")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestUnmarshalResponseMapToSlice(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := []struct {
		Tag1 string
		Tag2 string
	}{}
	want := []struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "test1_1",
		Tag2: "test1_2",
	}, {
		Tag1: "test2_1",
		Tag2: "test2_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToSliceWithContent(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := []struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "do not",
		Tag2: "override",
	}}
	want := []struct {
		Tag1 string
		Tag2 string
	}{{
		Tag1: "do not",
		Tag2: "override",
	}, {
		Tag1: "test1_1",
		Tag2: "test1_2",
	}, {
		Tag1: "test2_1",
		Tag2: "test2_2",
	}}

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err != nil {
		LogTestError(have, want, t, fmt.Sprintf("test failed with error %s", err))
	}
	for i, v := range have {
		if v.Tag1 != want[i].Tag1 {
			LogTestError(v.Tag1, want[i].Tag1, t)
		}
		if v.Tag2 != want[i].Tag2 {
			LogTestError(v.Tag2, want[i].Tag2, t)
		}
	}
}

func TestUnmarshalResponseMapToSliceError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test1_1",
		"Tag2": "test1_2",
	}, {
		"Tag1": "test2_1",
		"Tag2": "test2_2",
	}}
	have := []struct {
		Tag1 []bool
		Tag2 string
	}{}
	want := "1 error(s) decoding:\n\n* cannot parse 'Tag1[0]' as bool: strconv.ParseBool: parsing \"test1_1\": invalid syntax"

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected error unmarshaling")
	}
	if err.Error() != want {
		LogTestError(have, want, t)
	}
}

func TestUnmarshalResponseMapToUnsupportedTypeError(t *testing.T) {
	// given
	m := []map[string]interface{}{{
		"Tag1": "test",
	}}
	have := ""

	// when
	err := core.UnmarshalResponse(m, &have)

	// then
	if err == nil {
		LogTestError(have, "", t, "expected unmarshal error")
	}
	if err.Error() != "unsupported type string. Expected type struct, slice or array" {
		LogTestError(err.Error(), "unsupported type string. Expected type struct, slice or array", t)
	}
}

func TestDecodeMapToNonPtrError(t *testing.T) {
	// given
	m := map[string]interface{}{
		"Tag1": "test",
	}
	have := struct {
		Tag1 string
	}{}
	want := "expected pointer to value, not value"

	// when
	err := core.Decode(m, have)

	// then
	if err == nil {
		LogTestError(have, want, t, "expected error because of non pointer value")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestKeepAlive(t *testing.T) {
	// given
	b := &bytes.Buffer{}
	interval := 2
	timer := time.NewTimer(time.Duration(interval)*time.Second + 10*time.Millisecond)
	want := strings.Repeat(" \n", interval)

	// when
	go core.KeepAlive(b, time.Second)
	<-timer.C

	// then
	have := b.String()
	b = nil // "close" b
	if have != want {
		LogTestError(fmt.Sprintf("%q", have), fmt.Sprintf("%q", want), t)
	}
}

func TestReadHeader(t *testing.T) {
	// given
	c := make(chan []byte, 2)
	header := []byte("TS3")
	banner := []byte(" ")

	// when
	c <- header
	c <- banner
	err := core.ReadHeader(c)

	// then
	if err != nil {
		LogTestError("", "", t, "failed reading header", err.Error())
	}
}

func TestReadHeaderWrongHeaderError(t *testing.T) {
	// given
	c := make(chan []byte, 2)
	header := []byte("TS4")
	banner := []byte(" ")
	want := "wrong header"

	// when
	c <- header
	c <- banner
	err := core.ReadHeader(c)

	// then
	if err == nil {
		LogTestError("", want, t, "expected error reading header")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestReadHeaderChannelClosedError(t *testing.T) {
	// given
	c := make(chan []byte)
	close(c)
	want := "unable to read header: connection closed"

	// when
	err := core.ReadHeader(c)

	// then
	if err == nil {
		LogTestError("", want, t, "expected error reading header")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestRunSimpleResponse(t *testing.T) {
	// given
	writer := &bytes.Buffer{}
	in := make(chan []byte, 1)
	defer close(in)
	in <- []byte("error id=0 msg=\"\"")
	cmd := []byte("test")

	// when
	data, err := core.Run(in, writer, cmd)

	// then
	if err != nil {
		LogTestError(err.Error(), "", t, "error while running cmd")
	}
	if data != nil {
		LogTestError(data, "", t, "didn't expect return value")
	}
	have := writer.String()
	if have != string(cmd) {
		LogTestError(have, cmd, t)
	}
}

func TestRunComplexResponse(t *testing.T) {
	// given
	writer := &bytes.Buffer{}
	in := make(chan []byte, 2)
	defer close(in)
	in <- []byte("test")
	in <- []byte("error id=0 msg=")
	cmd := []byte("test")

	// when
	data, err := core.Run(in, writer, cmd)

	// then
	if err != nil {
		LogTestError(err.Error(), "", t, "error while running cmd")
	}
	if string(data) != "test" {
		LogTestError(data, "test", t)
	}
}

func TestRunClosedWriter(t *testing.T) {
	// given
	writer, err := ioutil.TempFile("./", "TestRunClosedWriter")
	if err != nil {
		LogTestError("", "", t, "error occured opening temp file")
	}
	writer.Close()
	in := make(chan []byte)
	defer func() {
		close(in)
		os.Remove(writer.Name())
	}()
	cmd := []byte("test")

	// when
	_, err = core.Run(in, writer, cmd)

	// then
	if err == nil {
		LogTestError("", "", t, "expected error using run on nil writer")
	}
	if err.Error() != fmt.Sprintf("write %s: file already closed", writer.Name()) {
		LogTestError(err.Error(), fmt.Sprintf("write %s: file already closed", writer.Name()), t)
	}
}

func TestRunClosedChannel(t *testing.T) {
	// given
	writer := &bytes.Buffer{}
	in := make(chan []byte)
	close(in)
	cmd := []byte("test")

	// when
	_, err := core.Run(in, writer, cmd)

	// then
	if err == nil {
		LogTestError("", "", t, "expected error using run with closed channel")
	}
	if err.Error() != "unable to read response: connection closed" {
		LogTestError(err.Error(), "unable to read response: connection closed", t)
	}
}

func TestRunErrorResponse(t *testing.T) {
	// given
	writer := ioutil.Discard
	in := make(chan []byte, 1)
	defer close(in)
	want := core.QueryError{
		ID:      1234,
		Message: "it worked",
	}
	in <- []byte(fmt.Sprintf("error id=%d msg=%s", want.ID, libts.QueryEncoder.Replace(want.Message)))

	// when
	_, err := core.Run(in, writer, []byte(""))

	// then
	if err == nil {
		LogTestError("", "", t, "expected error from run")
	}
	if _, ok := err.(core.QueryError); !ok {
		LogTestError(reflect.TypeOf(err), "core.QueryError", t)
	}
	if e := err.(core.QueryError); !(e.ID == want.ID) || !(e.Message == want.Message) {
		LogTestError(e, want, t)
	}
}

func TestRunUnexpectedError(t *testing.T) {
	// given
	writer := ioutil.Discard
	in := make(chan []byte, 1)
	in <- []byte("error id=true msg=fail")
	want := "1 error(s) decoding:\n\n* cannot parse 'id' as int: strconv.ParseInt: parsing \"true\": invalid syntax"

	// when
	_, err := core.Run(in, writer, []byte(""))

	// then
	if err == nil {
		LogTestError("", "", t, "expected error parsing error")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestDemultiplexer(t *testing.T) {
	// given
	out := make(chan []byte)
	notify := make(chan []byte)
	in := bytes.NewBufferString("\rtest\n\rerror id=0 msg=\n\rnotify test\n\r")

	// when
	go core.Demultiplexer(in, out, notify)

	// then
	go func() {
		e := <-out
		if string(e) != "test" {
			LogTestError(e, "test", t)
		}
		e = <-out
		if string(e) != "error id=0 msg=" {
			LogTestError(e, "error id=0 msg=", t)
		}
	}()
	n := <-notify
	if string(n) != "notify test" {
		LogTestError(n, "notify test", t)
	}
}
