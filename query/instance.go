package query

import "github.com/schoeppi5/libts"

// Instance returns the combined result of hostinfo and instanceinfo since I didn't really see point in dividing them
func (a Agent) Instance() (*Server, error) {
	server := &Server{}
	err := a.Query.Do(
		libts.Request{
			Command: "hostinfo",
		}, server)
	if err != nil {
		return nil, err
	}
	err = a.Query.Do(
		libts.Request{
			Command: "instanceinfo",
		}, server)
	if err != nil {
		return nil, err
	}
	return server, nil
}
