package sshquery

import (
	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/communication"
)

// Do tries to unmarshal the response from teamspeak to value
func (sq *SSHQuery) Do(request libts.Request, value interface{}) error {
	raw, err := sq.DoRaw(request)
	if err != nil {
		return err
	}
	return communication.UnmarshalResponse(communication.ConvertResponse(raw), value)
}

// DoRaw returns the raw response from teamspeak
func (sq *SSHQuery) DoRaw(request libts.Request) ([]byte, error) {
	if request.ServerID != 0 {
		err := sq.Use(request.ServerID)
		if err != nil {
			return nil, err
		}
	}
	return communication.Run(sq.in, sq.out, []byte(request.String()))
}

// Notification returns an io.Reader for arriving events
func (sq *SSHQuery) Notification() <-chan []byte {
	sq.notify = make(chan []byte, 5)
	return sq.notify
}

// Connected sends the version command and returns the recieved error, if any
func (sq *SSHQuery) Connected() (bool, error) {
	version := libts.Request{
		Command: "version",
	}
	_, err := sq.DoRaw(version)
	return err == nil, err
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
