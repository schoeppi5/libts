package serverquery

import (
	"time"

	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/core"
)

// This file fullfills the libts.Query interface

// Do a command and attempt to parse the response in value
func (sq *ServerQuery) Do(request libts.Request, value interface{}) error {
	raw, err := sq.DoRaw(request)
	if err != nil {
		return err
	}
	return core.UnmarshalResponse(core.ConvertResponse(raw), value)
}

// DoRaw a command and return the raw response
func (sq *ServerQuery) DoRaw(request libts.Request) ([]byte, error) {
	if request.ServerID != 0 {
		err := sq.use(request.ServerID)
		if err != nil {
			return nil, err
		}
	}
	sq.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	sq.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer func() {
		sq.conn.SetReadDeadline(time.Time{})
		sq.conn.SetWriteDeadline(time.Time{})
	}()
	return core.Run(sq.in, sq.conn, []byte(request.String()))
}

// Notification returns an io.Reader for arriving events
func (sq *ServerQuery) Notification() <-chan []byte {
	sq.notify = make(chan []byte, 5)
	return sq.notify
}

// use selects the virtual server for queries. Only changes the virtual server when needed
func (sq *ServerQuery) use(sid int) error {
	if sid == sq.serverID {
		return nil
	}
	use := libts.Request{
		Command: "use",
		Args: map[string]interface{}{
			"sid": sid,
		},
	}
	_, err := sq.DoRaw(use)
	if err != nil {
		return err
	}
	sq.serverID = sid
	return nil
}
