package query

import (
	"github.com/schoeppi5/libts"
)

// ChannelList returns a array of all channels
// Teamspeak returns different information dependening if you ask for one channel or multiple
// To compensate for that, ChannelList first asks for all channels and then for each channel individually resulting in a larger amount of requests
func (a Agent) ChannelList(sid int) ([]Channel, error) {
	list, err := a.ChannelIDList(sid)
	if err != nil {
		return nil, err
	}
	channels := make([]Channel, len(list))
	for i := range list {
		channel, err := a.Channel(sid, list[i])
		if err != nil {
			return nil, err
		}
		channels[i] = *channel
	}
	return channels, nil
}

// ChannelIDList only returns the ids of all channels on a virtual server
// This func executes the 'channellist' command, but discards all information except the cids for a uniform behaviour
func (a Agent) ChannelIDList(sid int) ([]int, error) {
	cList := []struct {
		ID int `mapstructure:"cid"`
	}{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "channellist",
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

// Channel returns detailed info about a channel specified by cid on a server specified by sid
func (a Agent) Channel(sid int, cid int) (*Channel, error) {
	channel := Channel{}
	err := a.Query.Do(
		libts.Request{
			ServerID: sid,
			Command:  "channelinfo",
			Args: map[string]interface{}{
				"cid": cid,
			},
		}, &channel)
	if err != nil {
		return nil, err
	}
	// manually set the cid because logic
	channel.ID = cid
	return &channel, nil
}
