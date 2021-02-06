package communication_test

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
	"github.com/schoeppi5/libts/communication"
)

func TestQueryErrorIsError(t *testing.T) {
	t.Parallel()
	// given
	qe := &communication.QueryError{}

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
	qe := &communication.QueryError{
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
	qe := &communication.QueryError{
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

func TestKeepAlive(t *testing.T) {
	// given
	b := &bytes.Buffer{}
	interval := 2
	timer := time.NewTimer(time.Duration(interval)*time.Second + 10*time.Millisecond)
	want := strings.Repeat(" \n", interval)

	// when
	go communication.KeepAlive(b, time.Second)
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
	err := communication.ReadHeader(c)

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
	err := communication.ReadHeader(c)

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
	err := communication.ReadHeader(c)

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
	data, err := communication.Run(in, writer, cmd)

	// then
	if err != nil {
		LogTestError(err.Error(), "", t, "error while running cmd")
	}
	if data != nil {
		LogTestError(data, "", t, "didn't expect return value")
	}
	have := writer.String()
	if have != fmt.Sprintf("%s\n", cmd) {
		LogTestError(have, fmt.Sprintf("%s\n", cmd), t)
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
	data, err := communication.Run(in, writer, cmd)

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
	_, err = communication.Run(in, writer, cmd)

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
	_, err := communication.Run(in, writer, cmd)

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
	want := communication.QueryError{
		ID:      1234,
		Message: "it worked",
	}
	in <- []byte(fmt.Sprintf("error id=%d msg=%s", want.ID, libts.QueryEncoder.Replace(want.Message)))

	// when
	_, err := communication.Run(in, writer, []byte(""))

	// then
	if err == nil {
		LogTestError("", "", t, "expected error from run")
	}
	if _, ok := err.(communication.QueryError); !ok {
		LogTestError(reflect.TypeOf(err), "communication.QueryError", t)
	}
	if e := err.(communication.QueryError); !(e.ID == want.ID) || !(e.Message == want.Message) {
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
	_, err := communication.Run(in, writer, []byte(""))

	// then
	if err == nil {
		LogTestError("", "", t, "expected error parsing error")
	}
	if err.Error() != want {
		LogTestError(err.Error(), want, t)
	}
}

func TestSplit(t *testing.T) {
	// given
	out := make(chan []byte)
	notify := make(chan []byte)
	in := bytes.NewBufferString("\rtest\n\rerror id=0 msg=\n\rnotify test\n\r")

	// when
	go communication.Split(in, out, notify)

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
