package query

import (
	"github.com/schoeppi5/libts"
)

// Client is a single logged in client on a virtual server
type Client struct {
	ID                         int
	Away                       int       `mapstructure:"client_away"`
	AwayMessage                string    `mapstructure:"client_away_message"`
	Base64Hash                 string    `mapstructure:"client_base64HashClientUID"`
	ChannelGroups              GroupList `mapstructure:"client_channel_group_id"`
	ChannelID                  int       `mapstructure:"cid"`
	ConnectedTime              int       `mapstructure:"connection_connected_time"`
	Country                    string    `mapstructure:"client_country"`
	DBID                       int       `mapstructure:"client_database_id"`
	DefaultChannel             string    `mapstructure:"client_default_channel"`
	DefaultToken               string    `mapstructure:"client_default_token"`
	Description                string    `mapstructure:"client_description"`
	EstimatedLocation          string    `mapstructure:"client_estimated_location"`
	FirstConnect               int64     `mapstructure:"client_created"`
	FlagAvatar                 string    `mapstructure:"client_flag_avatar"`
	IP                         string    `mapstructure:"connection_client_ip"`
	IconID                     uint32    `mapstructure:"client_icon_id"`
	IdleTime                   int       `mapstructure:"client_idle_time"`
	InputHardware              bool      `mapstructure:"client_input_hardware"`
	InputMuted                 bool      `mapstructure:"client_input_muted"`
	Integrations               string    `mapstructure:"client_integrations"`
	IsChannelCommander         bool      `mapstructure:"client_is_channel_commander"`
	IsPrioritySpeaker          bool      `mapstructure:"client_is_priority_speaker"`
	IsRecording                bool      `mapstructure:"client_is_recording"`
	IsServerQuery              bool      `mapstructure:"client_type"`
	IsTalker                   bool      `mapstructure:"client_is_talker"`
	LastConnect                int64     `mapstructure:"client_lastconnected"`
	LoginName                  string    `mapstructure:"client_login_name"`
	Metadata                   string    `mapstructure:"client_meta_data"`
	MyTeamspeakAvatar          string    `mapstructure:"client_myteamspeak_avatar"`
	MyTeamspeakID              string    `mapstructure:"client_myteamspeak_id"`
	NeededServerQueryViewPower int       `mapstructure:"client_needed_serverquery_view_power"`
	Nickname                   string    `mapstructure:"client_nickname"`
	NicknamePhonetic           string    `mapstructure:"client_nickname_phonetic"`
	OutputHardware             bool      `mapstructure:"client_output_hardware"`
	OutputMuted                bool      `mapstructure:"client_output_muted"`
	OutputOnlyMuted            bool      `mapstructure:"client_outputonly_muted"`
	Platform                   string    `mapstructure:"client_platform"`
	SecurityHash               string    `mapstructure:"client_security_hash"`
	ServerGroups               GroupList `mapstructure:"client_servergroups"`
	TalkRequest                bool      `mapstructure:"client_talk_request"`
	TalkRequestMessage         string    `mapstructure:"client_talk_request_msg"`
	Talkpower                  int       `mapstructure:"client_talk_power"`
	TotalConnections           int       `mapstructure:"client_totalconnections"`
	UID                        string    `mapstructure:"client_unique_identifier"`
	Version                    string    `mapstructure:"client_version"`
	VersionSign                string    `mapstructure:"client_version_sign"`
}

// DBClient is a single client in the database of a virtualserver
type DBClient struct {
	FlagAvatar       string `mapstructure:"client_flag_avatar"`
	Base64Hash       string `mapstructure:"client_base64HashClientUID"`
	DBID             int    `mapstructure:"client_database_id"`
	LastIP           string `mapstructure:"client_lastip"`
	Created          int64  `mapstructure:"client_created"`
	TotalConnections int    `mapstructure:"client_totalconnections"`
	UID              string `mapstructure:"client_unique_identifier"`
	Nickname         string `mapstructure:"client_nickname"`
	LastConnected    int64  `mapstructure:"client_lastconnected"`
	Description      string `mapstructure:"client_description"`
}

// ClientList returns all currently logged in clients
// sid - required
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
// sid - required
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

// Client returns client clid on server sid
// sid - required
// clid - required
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
// sid - required
// clids - required (returns when len(clids) == 0)
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

// MoveClient cid form channel "from" to channel "to"
// sid - required
// clid - required
// from - required
// to - required
// cpw - optional
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

// DBClients returns clients from database of server sid
// sid - required
// cldbids - required (returns when len(cldbids) == 0)
func (a Agent) DBClients(sid int, cldbids ...int) (*[]DBClient, error) {
	if len(cldbids) == 0 {
		return nil, nil
	}
	dbclient := &[]DBClient{}
	err := a.Query.Do(libts.Request{
		Command:  "clientdbinfo",
		ServerID: sid,
		Args: map[string]interface{}{
			"cldbid": cldbids,
		},
	}, dbclient)
	if err != nil {
		return nil, err
	}
	return dbclient, nil
}
