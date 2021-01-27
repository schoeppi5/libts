package sshquery

import (
	"github.com/schoeppi5/libts/core"

	"github.com/schoeppi5/libts"
)

// Do tries to unmarshal the response from teamspeak to value
func (sq *SSHQuery) Do(request libts.Request, value interface{}) error {
	raw, err := sq.DoRaw(request)
	if err != nil {
		return err
	}
	return core.UnmarshalResponse(core.ConvertResponse(raw), value)
}

// DoRaw returns the raw response from teamspeak
func (sq *SSHQuery) DoRaw(request libts.Request) ([]byte, error) {
	if request.ServerID != 0 {
		err := sq.Use(request.ServerID)
		if err != nil {
			return nil, err
		}
	}
	return core.Run(sq.in, sq.out, []byte(request.String()))
}

// Notification returns an io.Reader for arriving events
func (sq *SSHQuery) Notification() <-chan []byte {
	return sq.notify
}

// Use selects a virtual server for the client
func (sq *SSHQuery) Use(sid int) error {
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
