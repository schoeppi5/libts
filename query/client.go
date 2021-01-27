package query

import (
	"github.com/schoeppi5/libts"
)

// ClientList returns an array clients
func (a Agent) ClientList(sid int) ([]Client, error) {
	list, err := a.ClientIDList(sid)
	if err != nil {
		return nil, err
	}
	return a.Clients(sid, list...)
}

// ClientIDList returns a list of only the client ids
// This has the same problem as all the list commands
// It provides different information than the corresponding info command
func (a Agent) ClientIDList(sid int) ([]int, error) {
	cList := []struct {
		ID int `mapstructure:"clid"`
	}{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "clientlist",
		}, &cList)
	if err != nil {
		return nil, err
	}
	list := make([]int, len(cList))
	for i := range cList {
		list[i] = cList[i].ID
	}
	return list, nil
}

// Client returns a single client
func (a Agent) Client(sid int, clid int) (*Client, error) {
	client, err := a.Clients(sid, clid)
	if err != nil {
		return nil, err
	}
	return &client[0], nil
}

// Clients returns all clients specified by []clids
// Only here does the query provide a possibility to query multiple objects as you would expect them with one command
// So we might as well use it
func (a Agent) Clients(sid int, clids ...int) ([]Client, error) {
	if len(clids) == 0 { // Don't panic if there isn't anything to do
		return nil, nil
	}
	clients := []Client{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "clientinfo",
			Args:     map[string]interface{}{"clid": clids},
		}, &clients)
	if err != nil {
		return nil, err
	}
	// Add clientid to clients
	// This assumes that you get the clients as you requested them, which is the case atm
	// Watch out, they might screw this up some day
	for i := range clients {
		clients[i].ID = clids[i]
	}
	return clients, nil
}

// MoveClient cid form channel from to channel to
func (a Agent) MoveClient(sid, clid, from, to int, cpw string) error {
	if from == to { // Don't have to move the client
		return nil
	}
	_, err := a.Query.DoRaw(libts.Request{
		ServerID: sid,
		Command:  "clientmove",
		Args: map[string]interface{}{
			"clid": clid,
			"cid":  to,
			"cpw":  cpw,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
