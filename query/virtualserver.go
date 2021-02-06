package query

import "github.com/schoeppi5/libts"

// VirtualServerList returns an array of all virtual servers
func (a Agent) VirtualServerList() ([]VirtualServer, error) {
	list, err := a.VirtualServerIDList()
	if err != nil {
		return nil, err
	}
	vServers := make([]VirtualServer, len(list))
	for i := range list {
		vServer, err := a.VirtualServer(list[i])
		if err != nil {
			return nil, err
		}
		vServers[i] = *vServer
	}
	return vServers, nil
}

// VirtualServerIDList only returns the ids of all virtual servers
func (a Agent) VirtualServerIDList() ([]int, error) {
	vsList := []struct {
		ID int `mapstructure:"virtualserver_id"`
	}{}
	err := a.Query.Do(
		libts.Request{
			Command: "serverlist",
		}, &vsList)
	if err != nil {
		return nil, err
	}
	list := make([]int, len(vsList))
	for i := range vsList {
		list[i] = vsList[i].ID
	}
	return list, nil
}

// VirtualServer returns the virtual server represented by id
func (a Agent) VirtualServer(sid int) (*VirtualServer, error) {
	virtualServer := VirtualServer{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "serverinfo",
		}, &virtualServer)
	if err != nil {
		return nil, err
	}
	return &virtualServer, nil
}
