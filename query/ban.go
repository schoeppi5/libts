package query

import (
	"time"

	"github.com/schoeppi5/libts"
)

// Ban is a single ban on a virtualserver obtained by BanList
type Ban struct {
	ID           int    `mapstructure:"banid"`
	IP           string `mapstructure:"ip"`
	Created      int64  `mapstructure:"created"`
	InvokerName  string `mapstructure:"ivokername"`
	InvokerDBID  int    `mapstructure:"invokercldbid"`
	InvokerUID   string `mapstructure:"invokeruid"`
	Reason       string `mapstructure:"reason"`
	Enforcements int    `mapstructure:"enforcements"`
}

// BanAdd a ban to server sid. Client to ban can be specified using ip, name, uid or myTeamSpeakId (I don't know what the lastnickname does, since the docs didn't bother to mention it)
// sid - required
// ip - optional - regexp for last ip of client
// name - optional - regexp for nickname
// uid - optional - uid of client
// mytsid - optional - myTeamSpeak ID
// time - optional - Default 0 (infinite)
// reason - optional
// lastnickname - optional - I dunno what ever
// Possibly returns multiple banids if more than one matche is found
func (a Agent) BanAdd(sid int, ip string, name string, uid string, mytsid string, time time.Duration, reason string, lastnickname string) (int, error) {
	req := libts.Request{
		Command:  "banadd",
		ServerID: sid,
		Args: map[string]interface{}{
			"time":         int(time.Seconds()),
			"banreason":    reason,
			"lastnickname": lastnickname,
		},
	}
	if ip != "" {
		req.Args["ip"] = ip
	}
	if name != "" {
		req.Args["name"] = name
	}
	if uid != "" {
		req.Args["uid"] = uid
	}
	id := struct {
		ID int `mapstructure:"banid"`
	}{}
	err := a.Query.Do(req, &id)
	if err != nil {
		return 0, err
	}
	return id.ID, nil
}

// BanClients clid on server sid for time time for reason reason
// sid - required
// time - optional
// reason - optional
// clid - required
func (a Agent) BanClients(sid int, time time.Duration, reason string, clid ...int) ([]int, error) {
	req := libts.Request{
		Command:  "banclient",
		ServerID: sid,
		Args: map[string]interface{}{
			"time":             int(time.Seconds()),
			"reason":           reason,
			"clid":             clid,
			"-continueonerror": "",
		},
	}
	ids := []struct {
		ID int `mapstructure:"banid"`
	}{}
	err := a.Query.Do(req, &ids)
	if err != nil {
		return nil, err
	}
	bans := []int{}
	for i := range ids {
		bans = append(bans, ids[i].ID)
	}
	return bans, nil
}

// BanDel delete ban banid on server sid
// sid - required
// banid - required
func (a Agent) BanDel(sid int, banid int) error {
	req := libts.Request{
		Command:  "bandel",
		ServerID: sid,
		Args: map[string]interface{}{
			"banid": banid,
		},
	}
	_, err := a.Query.DoRaw(req)
	return err
}

// BanDelAll bans on server sid
// sid - required
func (a Agent) BanDelAll(sid int) error {
	req := libts.Request{
		Command:  "bandelall",
		ServerID: sid,
	}
	_, err := a.Query.DoRaw(req)
	return err
}

// BanList all bans on server sid with offset and a maximum of count
// sid - required
// offset - optional - Default 0 - skip the first 'offset' bans
// count - optional - Default 0
func (a Agent) BanList(sid int, offset int, count int) error {
	req := libts.Request{
		Command:  "banlist",
		ServerID: sid,
		Args: map[string]interface{}{
			"start": offset,
		},
	}
	if count != 0 {
		req.Args["duration"] = count
	}
	_, err := a.Query.DoRaw(req)
	return err
}
