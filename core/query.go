package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
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

// Run writes r to in and reads until it encounters an error. If error has id 0, the read data is returned
func Run(in <-chan []byte, out io.Writer, r []byte) ([]byte, error) {
	if !bytes.HasSuffix(r, []byte("\n")) { // a command must be suffixed by \n
		r = append(r, byte('\n'))
	}
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
		if err = IsError(d); err != nil {
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

// Split the input from c
// notifications (prefixed with notify.*) are send to notify, everything else is send to out
// Stops when c is closed or it encounters an error while reading from c
// If notify is nil, notifications are discarded
// If out is nil, Split returns
func Split(c io.Reader, out chan<- []byte, notify chan<- []byte) {
	if out == nil {
		return
	}
	reader := bufio.NewReader(c)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			close(out)
			close(notify)
			return
		}
		data = bytes.TrimRight(bytes.TrimLeft(data, "\r"), "\n") // normalize data
		if bytes.HasPrefix(data, []byte("notify")) {
			if notify == nil {
				continue
			}
			select { // non blocking write to notify
			case notify <- data:
			default:
			}
			continue
		}
		out <- data // block on write to out
	}
}
