package query

import (
	"fmt"
	"strconv"
	"strings"
)

// Server represents the actual server instance
// Can contain multiple virtual servers
// This is the combined result of "hostinfo" and "instanceinfo"
type Server struct {
	Uptime                         int64 `mapstructure:"instance_uptime"`
	ServerTime                     int64 `mapstructure:"host_timestamp_utc"`
	VirtualServerCount             int   `mapstructure:"virtualservers_running_total"`
	MaxClients                     int   `mapstructure:"virtualservers_total_maxclients"`
	TotalClients                   int   `mapstructure:"virtualservers_total_clients_online"`
	TotalChannels                  int   `mapstructure:"virtualservers_total_channels_online"`
	FTPPort                        int   `mapstructure:"serverinstance_filetransfer_port"`
	DefaultServerQueryGroup        int   `mapstructure:"serverinstance_guest_serverquery_group"`
	MaxServerQueryCommands         int   `mapstructure:"serverinstance_serverquery_flood_commands"`
	ServerQueryFloodTime           int   `mapstructure:"serverinstance_serverquery_flood_time"` //seconds
	ServerQueryFloodBanTime        int   `mapstructure:"serverinstance_serverquery_ban_time"`   //seconds
	DatabaseVersion                int   `mapstructure:"serverinstance_database_version"`
	PendingConnectionsPerIP        int   `mapstructure:"serverinstance_pending_connections_per_ip"`
	PermissionsVersion             int   `mapstructure:"serverinstance_permissions_version"`
	ServerQueryMaxConnectionsPerIP int   `mapstructure:"serverinstance_serverquery_max_connections_per_ip"`
	TemplateChannelAdminGroup      int   `mapstructure:"serverinstance_template_channeladmin_group"`
	TemplateChannelDefaultGroup    int   `mapstructure:"serverinstance_template_channeldefault_group"`
	TemplateServerAdminGroup       int   `mapstructure:"serverinstance_template_serveradmin_group"`
	TemplateServerDefaultGroup     int   `mapstructure:"serverinstance_template_serverdefault_group"`
}

// VirtualServer represents a single virtual server instance
type VirtualServer struct {
	UID                              string  `mapstructure:"virtualserver_unique_identifier"`
	Name                             string  `mapstructure:"virtualserver_name"`
	WelcomeMessage                   string  `mapstructure:"virtualserver_welcomemessage"`
	Platform                         string  `mapstructure:"virtualserver_platform"`
	Version                          string  `mapstructure:"virtualserver_version"`
	MaxClients                       int     `mapstructure:"virtualserver_maxclients"`
	TotalClients                     int     `mapstructure:"virtualserver_clientsonline"`
	ChannelCount                     int     `mapstructure:"virtualserver_channelsonline"`
	Created                          int64   `mapstructure:"virtualserver_created"`
	Uptime                           int64   `mapstructure:"virtualserver_uptime"`
	CodecEncryption                  int     `mapstructure:"virtualserver_codec_encryption_mod"`
	HostMessage                      string  `mapstructure:"virtualserver_hostmessage"`
	HostMessageMode                  int     `mapstructure:"virtualserver_hostmessage_mode"`
	FileBase                         string  `mapstructure:"virtualserver_filebase"`
	DefaultServerGroup               int     `mapstructure:"virtualserver_default_server_group"`
	DefaultChannelGroup              int     `mapstructure:"virtualserver_default_channel_group"`
	Password                         bool    `mapstructure:"virtualserver_flag_password"`
	DefaultChannelAdminGroup         int     `mapstructure:"virtualserver_default_channel_admin_group"`
	HostbannerURL                    string  `mapstructure:"virtualserver_hostbanner_url"`
	HostbannerGFXURL                 string  `mapstructure:"virtualserver_hostbanner_gfx_url"`
	HostbannerGFXInterval            string  `mapstructure:"virtualserver_hostbanner_gfx_interval"`
	ComplainAutobanCount             int     `mapstructure:"virtualserver_complain_autoban_count"`
	ComplainAutobanTime              int     `mapstructure:"virtualserver_complain_autoban_time"`
	ComplainRemoveTime               int     `mapstructure:"virtualserver_complain_remove_time"`
	MinClientsInChannelForcedSilence int     `mapstructure:"virtualserver_min_clients_in_channel_before_forced_silence"`
	PrioritySpeakerMod               float32 `mapstructure:"virtualserver_priority_speaker_dimm_modificator"`
	ID                               int     `mapstructure:"virtualserver_id"`
	AntifloodPointsTickReduce        int     `mapstructure:"virtualserver_antiflood_points_tick_reduce"`
	AntifloodCommandBlock            int     `mapstructure:"virtualserver_antiflood_points_needed_command_block"`
	AntifloodIPBlock                 int     `mapstructure:"virtualserver_antiflood_points_needed_ip_block"`
	TotalClientConnections           int     `mapstructure:"virtualserver_client_connections"`       //I think thats all connections since startup (not creation)
	TotalQueryConnections            int     `mapstructure:"virtualserver_query_client_connections"` // Same here
	HostbuttonToolTip                string  `mapstructure:"virtualserver_hostbutton_tooltip"`
	HostbuttonURL                    string  `mapstructure:"virtualserver_hostbutton_url"`
	HostbuttonGFXURL                 string  `mapstructure:"virtualserver_hostbutton_gfx_url"`
	QueryClientCount                 int     `mapstructure:"virtualserver_queryclientsonline"`
	Port                             int     `mapstructure:"virtualserver_port"`
	Autostart                        bool    `mapstructure:"virtualserver_autostart"`
	SecurityLevel                    int     `mapstructure:"virtualserver_needed_identity_security_level"`
	NamePhonetic                     string  `mapstructure:"virtualserver_name_phonetic"`
	IconID                           int     `mapstructure:"virtualserver_icon_id"`
	ReservedSlots                    int     `mapstructure:"virtualserver_reserved_slots"`
	Ping                             float32 `mapstructure:"virtualserver_total_ping"`
	Weblist                          bool    `mapstructure:"virtualserver_weblist_enabled"`
	HostbannerMode                   int     `mapstructure:"virtualserver_hostbanner_mode"`
	TempChannelDeleteDelayDefault    int     `mapstructure:"virtualserver_channel_temp_delete_delay_default"`
	Nickname                         string  `mapstructure:"virtualserver_nickname"`
	AntifloodPluginBlock             int     `mapstructure:"virtualserver_antiflood_points_needed_plugin_block"`
	Status                           string  `mapstructure:"virtualserver_status"`
}

// Channel is a single Channel on a virtual server
type Channel struct {
	ID               int    `mapstructure:"id"`
	ParentID         int    `mapstructure:"pid"`
	IconID           int    `mapstructure:"channel_icon_id"`
	Name             string `mapstructure:"channel_name"`
	Topic            string `mapstructure:"channel_topic"`
	Description      string `mapstructure:"channel_description"`
	Password         bool   `mapstructure:"channel_flag_password"`
	Codec            Codec  `mapstructure:"channel_codec"`
	CodecQuality     int    `mapstructure:"channel_codec_quality"`
	MaxClients       int    `mapstructure:"channel_maxclients"`
	MaxFamilyClients int    `mapstructure:"channel_maxfamilyclients"`
	Order            int    `mapstructure:"channel_order"`
	Permanent        bool   `mapstructure:"channel_flag_permanent"`
	SemiPermanent    bool   `mapstructure:"channel_flag_semi_permanent"`
	Temporary        bool   `mapstructure:"channel_flag_temporary"`
	DefaultChannel   bool   `mapstructure:"channel_flag_default"`
	TalkPower        int    `mapstructure:"channel_needed_talk_power"`
	NamePhonetic     string `mapstructure:"channel_name_phonetic"`
	FilePath         string `mapstructure:"channel_filepath"`
	Silenced         bool   `mapstructure:"channel_forced_silence"`
	SecondsEmpty     int64  `mapstructure:"seconds_empty"`
}

// Group represents a group on a virtual server (server|channel)
type Group struct {
	ID     int    `mapstructure:"sgid" mapstructure:"cgid"`
	Name   string `mapstructure:"name"`
	Type   int    `mapstructure:"type"`   // 1 -> Server | 2 -> Channel
	IconID int32  `mapstructure:"iconid"` // ok so here me out: There is a bug, that is not a bug. Read all about it here (https://community.teamspeak.com/t/bug-query-sends-wrong-icon-id-in-response/15054)
	SaveDB bool   `mapstructure:"savedb"`
}

// Codec represents the possible codecs of a channel
type Codec string

// UnmarshalText turns number to correct Codec
func (c *Codec) UnmarshalText(text []byte) error {
	switch string(text) {
	case "0":
		*c = "Speex Narrowband"
	case "1":
		*c = "Speex Wideband"
	case "2":
		*c = "Speex Ultrawideband"
	case "3":
		*c = "Celt Mono"
	case "4":
		*c = "Opus Voice"
	case "5":
		*c = "Opus Music"
	default:
		return fmt.Errorf("Failed to unmarshal %s to codec", text)
	}
	return nil
}

// GroupList represents a list of group ids
type GroupList []int

// UnmarshalText unmarshals Teamspeaks "arrays" into real int arrays
func (gl *GroupList) UnmarshalText(text []byte) error {
	groups := strings.Split(string(text), ",")
	list := make([]int, len(groups))
	for i := range groups {
		j, err := strconv.Atoi(groups[i])
		if err != nil {
			return err
		}
		list[i] = j
	}
	*gl = GroupList(list)
	return nil
}
