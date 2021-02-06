package query

import "github.com/schoeppi5/libts"

// There are no info commands for groups

// ServerGroupList returns a list of all servergroups
// Type Servergroups -> 2 = Server Query Group | 0 = Template | 1 = Normal (I think)
func (a Agent) ServerGroupList(sid int) ([]Group, error) {
	groups := []Group{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "servergrouplist",
		}, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// ChannelGroupList returns a list of all servergroups
// Type Channelgroups -> 0 = Template | 1 = Normal
func (a Agent) ChannelGroupList(sid int) ([]Group, error) {
	groups := []Group{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "channelgrouplist",
		}, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
