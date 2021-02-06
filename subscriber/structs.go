package subscriber

import (
	"fmt"
	"strconv"

	"github.com/schoeppi5/libts/query"
)

// This file contains structs for all event types

// TextMessageEvent event for TextMessage
type TextMessageEvent struct {
	// 1 = privat | 2 = channel | 3 = server
	TargetMode int    `mapstructure:"targetmode"`
	Message    string `mapstructure:"msg"`
	// Only set for private messages
	Target int `mapstructure:"target"`
	// 0 for gms
	InvokerID int `mapstructure:"invokerid"`
	// Server for gms
	InvokerName string `mapstructure:"invokername"`
	InvokerUID  string `mapstructure:"invokeruid"`
}

// TokenUsedEvent event for TokenUsed
type TokenUsedEvent struct {
	ClientID       int    `mapstructure:"clid"`
	ClientDBID     int    `mapstructure:"cldbid"`
	ClientUID      string `mapstructure:"cluid"`
	Token          string `mapstructure:"token"`
	TokenCustomSet string `mapstructure:"tokencustomset"`
	Token1         string `mapstructure:"token1"`
	Token2         int    `mapstructure:"token2"`
}

// ClientMovedEvent event for ClientMoved
type ClientMovedEvent struct {
	To          int    `mapstructure:"ctid"`
	Reason      Reason `mapstructure:"reasonid"`
	InvokerID   int    `mapstructure:"invokerid"`
	InvokerName string `mapstructure:"invokername"`
	InvokerUID  string `mapstructure:"invokeruid"`
	ID          int    `mapstructure:"clid"`
}

// ChannelEditedEvent event for ChannelEdited
type ChannelEditedEvent struct {
	ID                        int         `mapstructure:"cid"`
	Reason                    Reason      `mapstructure:"reasonid"`
	InvokerID                 int         `mapstructure:"invokerid"`
	InvokerName               string      `mapstructure:"invokername"`
	InvokerUID                string      `mapstructure:"invokeruid"`
	Name                      string      `mapstructure:"channel_name"`
	Topic                     string      `mapstructure:"channel_topic"`
	Codec                     query.Codec `mapstructure:"channel_codec"`
	CodecQuality              int         `mapstructure:"channel_codec_quality"`
	MaxClients                int         `mapstructure:"channel_maxclients"`
	MaxFamilyClients          int         `mapstructure:"channel_maxfamilyclients"`
	Order                     int         `mapstructure:"channel_order"`
	Permanent                 bool        `mapstructure:"channel_flag_permanent"`
	SemiPermanent             bool        `mapstructure:"channel_flag_semi_permanent"`
	Default                   bool        `mapstructure:"channel_flag_default"`
	Password                  bool        `mapstructure:"channel_flag_password"`
	CodecLatencyFactor        int         `mapstructure:"channel_codec_latency_factor"`
	CodecIsUnencrypted        int         `mapstructure:"channel_codec_is_unencrypted"`
	DeleteDelay               int         `mapstructure:"channel_delete_delay"`
	ClientsUnlimited          bool        `mapstructure:"channel_flag_maxclients_unlimited"`
	FamilyClientsUnlimited    bool        `mapstructure:"channel_flag_maxfamilyclients_unlimited"`
	MaxFamilyClientsInherited bool        `mapstructure:"channel_flag_maxfamilyclients_inherited"`
	NeededTalkPower           int         `mapstructure:"channel_needed_talk_power"`
	NamePhonetic              string      `mapstructure:"channel_name_phonetic"`
	IconID                    int         `mapstructure:"channel_icon_id"`
}

// ChannelDeletedEvent event for ChannelDeleted
type ChannelDeletedEvent struct {
	InvokerID   int    `mapstructure:"invokerid"`
	InvokerName string `mapstructure:"invokername"`
	InvokerUID  string `mapstructure:"invokeruid"`
	ID          int    `mapstructure:"cid"`
}

// ChannelCreatedEvent event for ChannelCreated
type ChannelCreatedEvent struct {
	ID                        int         `mapstructure:"cid"`
	ParentID                  int         `mapstructure:"cpid"`
	Name                      string      `mapstructure:"channel_name"`
	Topic                     string      `mapstructure:"channel_topic"`
	Codec                     query.Codec `mapstructure:"channel_codec"`
	CodecQuality              int         `mapstructure:"channel_codec_quality"`
	MaxClients                int         `mapstructure:"channel_maxclients"`
	MaxFamilyClients          int         `mapstructure:"channel_maxfamilyclients"`
	Order                     int         `mapstructure:"channel_order"`
	Permanent                 bool        `mapstructure:"channel_flag_permanent"`
	SemiPermanent             bool        `mapstructure:"channel_flag_semi_permanent"`
	Default                   bool        `mapstructure:"channel_flag_default"`
	Password                  bool        `mapstructure:"channel_flag_password"`
	CodecLatencyFactor        int         `mapstructure:"channel_codec_latency_factor"`
	CodecIsUnencrypted        bool        `mapstructure:"channel_codec_is_unencrypted"`
	DeleteDelay               int         `mapstructure:"channel_delete_delay"`
	ClientsUnlimited          bool        `mapstructure:"channel_flag_maxclients_unlimited"`
	FamilyClientsUnlimited    bool        `mapstructure:"channel_flag_maxfamilyclients_unlimited"`
	MaxFamilyClientsInherited bool        `mapstructure:"channel_flag_maxfamilyclients_inherited"`
	NeededTalkPower           int         `mapstructure:"channel_needed_talk_power"`
	NamePhoenetic             string      `mapstructure:"channel_name_phonetic"`
	IconID                    int         `mapstructure:"channel_icon_id"`
	InvokerID                 int         `mapstructure:"invokerid"`
	InvokerName               string      `mapstructure:"invokername"`
	InvokerUID                string      `mapstructure:"invokeruid"`
}

// ChannelMovedEvent event for ChannelMoved
type ChannelMovedEvent struct {
	ID          int    `mapstructure:"cid"`
	ParentID    int    `mapstructure:"cpid"`
	Order       int    `mapstructure:"order"`
	Reason      Reason `mapstructure:"reasonid"`
	InvokerID   int    `mapstructure:"invokerid"`
	InvokerName string `mapstructure:"invokername"`
	InvokerUID  string `mapstructure:"invokeruid"`
}

// ChannelDescriptionChangedEvent event for ChannelDescriptionChanged
type ChannelDescriptionChangedEvent struct {
	ID int `mapstructure:"cid"`
}

// ChannelPasswordChangedEvent event for ChannelPasswordChanged
type ChannelPasswordChangedEvent struct {
	ID int `mapstructure:"cid"`
}

// ServerEditedEvent event for ServerEdited
type ServerEditedEvent struct {
	Reason                          Reason `mapstructure:"reasonid"`
	InvokerID                       int    `mapstructure:"invokerid"`
	InvokerName                     string `mapstructure:"invokername"`
	InvokerUID                      string `mapstructure:"invokeruid"`
	Name                            string `mapstructure:"virtualserver_name"`
	CodecEncryptionMode             string `mapstructure:"virtualserver_codec_encryption_mode"`
	DefaultServerGroup              int    `mapstructure:"virtualserver_default_server_group"`
	DefaultChannelGroup             int    `mapstructure:"virtualserver_default_channel_group"`
	HostbannerURL                   string `mapstructure:"virtualserver_hostbanner_url"`
	HostbannerGFXURL                string `mapstructure:"virtualserver_hostbanner_gfx_url"`
	HostbannerGFXInterval           int    `mapstructure:"virtualserver_hostbanner_gfx_interval"`
	PrioritySpeakerDimmModification int    `mapstructure:"virtualserver_priority_speaker_dimm_modification"`
	HostbuttonTooltip               string `mapstructure:"virtualserver_hostbutton_tooltip"`
	HostbuttonURL                   string `mapstructure:"virtualserver_hostbutton_url"`
	HostbuttonGFXURL                string `mapstructure:"virtualserver_hostbutton_gfx_url"`
	NamePhoenetic                   string `mapstructure:"virtualserver_name_phoenetic"`
	IconID                          int    `mapstructure:"virtualserver_icon_id"`
	HostbannerMode                  string `mapstructure:"virtualserver_hostbanner_mode"`
	TempChannelDefaultDeleteDelay   int    `mapstructure:"virtualserver_channel_temp_delete_delay_default"`
}

// ClientLeftViewEvent event for ClientLeftView
type ClientLeftViewEvent struct {
	From          int    `mapstructure:"cfid"`
	To            int    `mapstructure:"ctid"`
	Reason        Reason `mapstructure:"reasonid"`
	InvokerID     int    `mapstructure:"invokerid"`
	InvokerName   string `mapstructure:"invokername"`
	InvokerUID    string `mapstructure:"invokeruid"`
	ReasonMessage string `mapstructure:"reasonmsg"`
	Bantime       int    `mapstructure:"bantime"`
	ID            int    `mapstructure:"clid"`
}

// ClientEnterViewEvent event for ClientEnterView
type ClientEnterViewEvent struct {
	From                           int             `mapstructure:"cfid"`
	To                             int             `mapstructure:"ctid"`
	Reason                         Reason          `mapstructure:"reasonid"`
	ID                             int             `mapstructure:"clid"`
	UID                            string          `mapstructure:"client_unique_identifier"`
	Nickname                       string          `mapstructure:"client_nickname"`
	InputMuted                     bool            `mapstructure:"client_input_muted"`
	OutputMuted                    bool            `mapstructure:"client_output_muted"`
	OutputOnlyMuted                bool            `mapstructure:"client_output_muted"`
	InputHardware                  bool            `mapstructure:"client_input_hardware"`
	OutputHardware                 bool            `mapstructure:"client_output_hardware"`
	IsRecording                    bool            `mapstructure:"client_is_recording"`
	DBID                           int             `mapstructure:"client_database_id"`
	ChannelGroupID                 int             `mapstructure:"client_channel_group_id"`
	ServerGroups                   query.GroupList `mapstructure:"client_servergroups"`
	Away                           bool            `mapstructure:"client_away"`
	AwayMessage                    string          `mapstructure:"client_away_message"`
	IsServerQuery                  bool            `mapstructure:"client_type"`
	AvatarFlag                     string          `mapstructure:"client_flag_avatar"`
	TalkPower                      int             `mapstructure:"client_talk_power"`
	TalkRequest                    bool            `mapstructure:"client_talk_request"`
	TalkRequestMessage             string          `mapstructure:"client_talk_request_msg"`
	Description                    string          `mapstructure:"client_description"`
	IsTalker                       bool            `mapstructure:"client_is_talker"`
	IsPrioritySpeaker              bool            `mapstructure:"client_is_priority_speaker"`
	NicknamePhonetic               string          `mapstructure:"client_nickname_phonetic"`
	ServerQueryNeededViewPower     int             `mapstructure:"client_needed_serverquery_view_power"`
	IconID                         int             `mapstructure:"client_icon_id"`
	IsChannelCommander             string          `mapstructure:"client_is_channel_commander"`
	Country                        string          `mapstructure:"client_country"`
	ChannelGroupInheritedChannelID int             `mapstructure:"client_channel_group_inherited_channel_id"`
	Badges                         string          `mapstructure:"client_badges"`
}

// Reason for different events
type Reason string

// UnmarshalText to Reason
func (r *Reason) UnmarshalText(text []byte) error {
	i, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	switch i {
	case 0:
		*r = "Moved itself"
	case 1:
		*r = "Moved"
	case 3:
		*r = "Timeout"
	case 4:
		*r = "Kicked from channel"
	case 5:
		*r = "Kicked from server"
	case 6:
		*r = "Banned"
	case 8:
		*r = "Left"
	case 10:
		*r = "Edited"
	case 11:
		*r = "Server shutdown"
	default:
		return fmt.Errorf("failed to unmarshal %s to reason", text)
	}
	return nil
}
