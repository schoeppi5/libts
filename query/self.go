package query

import (
	"github.com/schoeppi5/libts"
)

// Me is the answer to WhoAmI ;)
type Me struct {
	ChannelID      int    `mapstructure:"client_channel_id"`
	DBID           int    `mapstructure:"client_database_id"`
	ID             int    `mapstructure:"client_id"`
	Name           string `mapstructure:"client_login_name"`
	Nickname       string `mapstructure:"client_nickname"`
	OriginServerID int    `mapstructure:"client_origin_server_id"`
	UID            string `mapstructure:"client_unique_identifier"`
	ServerID       int    `mapstructure:"virtualserver_id"`
	ServerPort     int    `mapstructure:"virtualserver_port"`
	ServerStatus   string `mapstructure:"virtualserver_status"`
	ServerUID      string `mapstructure:"virtualserver_unique_identifier"`
}

// WhoAmI returns Me
// sid - optional - It depends on the used query what happens. WebQuery will execute WhoAmI without a server. Telnet and SSH query will use the currently selected virtual server
// The output (filled fields) obviously vary based on the circumstances of the call (with or whithout sid)
func (a Agent) WhoAmI(sid int) (*Me, error) {
	whoami := libts.Request{
		Command: "whoami",
	}
	if sid != 0 {
		whoami.ServerID = sid
	}
	me := &Me{}
	err := a.Query.Do(whoami, me)
	if err != nil {
		return nil, err
	}
	return me, nil
}
