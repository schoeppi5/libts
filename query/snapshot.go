package query

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/DataDog/zstd"
	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/communication"
)

const (
	endVirtualserver              = "end_virtualserver"
	startChannels                 = "begin_channels"
	endChannels                   = "end_channels"
	startClients                  = "begin_clients"
	endClients                    = "end_clients"
	startPermissions              = "begin_permissions"
	endPermissions                = "end_permissions"
	startServerGroups             = "server_groups"
	startChannelGroups            = "channel_groups"
	endGroup                      = "end_group"
	endGroups                     = "end_groups"
	endRelations                  = "end_relations"
	startClientPermissions        = "client_flat"
	endClientPermissions          = "end_flat"
	startChannelPermissions       = "channel_flat"
	endChannelPermissions         = "end_flat"
	startChannelClientPermissions = "channel_client_flat"
	endChannelClientPermissions   = "end_flat"
	startAPIKeys                  = "begin_apikeys"
	endAPIKeys                    = "end_apikeys"
)

// RawSnapshot contains a created snapshot from a virtualserver
type RawSnapshot struct {
	Version int    `mapstructure:"version"`
	Salt    string `mapstructure:"salt"`
	Data    string `mapstructure:"data"`
}

// Snapshot contains a decoded and decompressed snapshot
type Snapshot struct {
	VirtualServer struct {
		MaxClients                       int     `mapstructure:"virtualserver_maxclients"`
		SecurityLevel                    int     `mapstructure:"virtualserver_needed_identity_security_level"`
		Weblist                          bool    `mapstructure:"virtualserver_weblist_enabled"`
		HostbannerGFXURL                 string  `mapstructure:"virtualserver_hostbanner_gfx_url"`
		ComplainAutobanCount             int     `mapstructure:"virtualserver_complain_autoban_count"`
		IconID                           int     `mapstructure:"virtualserver_icon_id"`
		Password                         string  `mapstructure:"virtualserver_password"`
		DefaultChannelGroup              int     `mapstructure:"virtualserver_default_channel_group"`
		HostbannerURL                    string  `mapstructure:"virtualserver_hostbanner_url"`
		PrioritySpeakerMod               float32 `mapstructure:"virtualserver_priority_speaker_dimm_modificator"`
		FileStorageClass                 string  `mapstructure:"virtualserver_file_storage_class"`
		LogChannel                       bool    `mapstructure:"virtualserver_log_channel"`
		LogServer                        bool    `mapstructure:"virtualserver_log_server"`
		CodecEncryption                  int     `mapstructure:"virtualserver_codec_encryption_mode"`
		HostbuttonToolTip                string  `mapstructure:"virtualserver_hostbutton_tooltip"`
		DefaultServerGroup               int     `mapstructure:"virtualserver_default_server_group"`
		LogPermissions                   bool    `mapstructure:"virtualserver_log_permissions"`
		HostMessage                      string  `mapstructure:"virtualserver_hostmessage"`
		HostMessageMode                  int     `mapstructure:"virtualserver_hostmessage_mode"`
		HasPassword                      bool    `mapstructure:"virtualserver_flag_password"`
		LogQuery                         bool    `mapstructure:"virtualserver_log_query"`
		MaxDownloadTotal                 big.Int `mapstructure:"virtualserver_max_download_total_bandwidth"`
		UploadQuota                      big.Int `mapstructure:"virtualserver_upload_quota"`
		ReservedSlots                    int     `mapstructure:"virtualserver_reserved_slots"`
		AntifloodPluginBlock             int     `mapstructure:"virtualserver_antiflood_points_needed_plugin_block"`
		LogClient                        bool    `mapstructure:"virtualserver_log_client"`
		Keypair                          string  `mapstructure:"virtualserver_keypair"`
		DefaultChannelAdminGroup         int     `mapstructure:"virtualserver_default_channel_admin_group"`
		AntifloodPointsTickReduce        int     `mapstructure:"virtualserver_antiflood_points_tick_reduce"`
		AntifloodIPBlock                 int     `mapstructure:"virtualserver_antiflood_points_needed_ip_block"`
		HostbannerMode                   int     `mapstructure:"virtualserver_hostbanner_mode"`
		MinIOSVersion                    string  `mapstructure:"virtualserver_min_ios_version"`
		ProtocolVerifyKeypair            string  `mapstructure:"virtualserver_protocol_verify_keypair"`
		Name                             string  `mapstructure:"virtualserver_name"`
		LogFileTransfer                  bool    `mapstructure:"virtualserver_log_filetransfer"`
		Nickname                         string  `mapstructure:"virtualserver_nickname"`
		Filebase                         string  `mapstructure:"virtualserver_filebase"`
		HostbannerGFXInterval            string  `mapstructure:"virtualserver_hostbanner_gfx_interval"`
		AccountingToken                  string  `mapstructure:"virtualserver_accounting_token"`
		Created                          int64   `mapstructure:"virtualserver_created"`
		MinClientsInChannelForcedSilence int     `mapstructure:"virtualserver_min_clients_in_channel_before_forced_silence"`
		HostbuttonURL                    string  `mapstructure:"virtualserver_hostbutton_url"`
		HostbuttonGFXURL                 string  `mapstructure:"virtualserver_hostbutton_gfx_url"`
		TempChannelDeleteDelayDefault    int     `mapstructure:"virtualserver_channel_temp_delete_delay_default"`
		WelcomeMessage                   string  `mapstructure:"virtualserver_welcomemessage"`
		MinAndroidVersion                string  `mapstructure:"virtualserver_min_android_version"`
		AntifloodCommandBlock            int     `mapstructure:"virtualserver_antiflood_points_needed_command_block"`
		MinClientVersion                 string  `mapstructure:"virtualserver_min_client_version"`
		UID                              string  `mapstructure:"virtualserver_unique_identifier"`
		MaxUploadTotalBandwidth          big.Int `mapstructure:"virtualserver_max_upload_total_bandwidth"`
		DownloadQuota                    big.Int `mapstructure:"virtualserver_download_quota"`
		ComplainAutobanTime              int64   `mapstructure:"virtualserver_complain_autoban_time"`
		ComplainRemoveTime               int64   `mapstructure:"virtualserver_complain_remove_time"`
		NamePhonetic                     string  `mapstructure:"virtualserver_name_phonetic"`
	}
	Channels []struct {
		PID                         int    `mapstructure:"channel_pid"`
		Description                 string `mapstructure:"channel_description"`
		MaxClients                  int    `mapstructure:"channel_maxclients"`
		IsPermanent                 bool   `mapstructure:"channel_flag_permanent"`
		CodecIsUnencrypted          bool   `mapstructure:"channel_codec_is_unencrypted"`
		BannerMode                  int    `mapstructure:"channel_banner_mode"`
		Password                    string `mapstructure:"channel_password"`
		CodecQuality                int    `mapstructure:"channel_codec_quality"`
		Order                       int    `mapstructure:"channel_order"`
		HasPassword                 bool   `mapstructure:"channel_flag_password"`
		MaxFamilyClientsIsUnlimited bool   `mapstructure:"channel_flag_maxfamilyclients_unlimited"`
		MaxFamilyClientsIsInherited bool   `mapstructure:"channel_flag_maxfamilyclients_inherited"`
		ID                          int    `mapstructure:"channel_id"`
		Name                        string `mapstructure:"channel_name"`
		Codec                       Codec  `mapstructure:"channel_codec"`
		CodecLatencyFactor          int    `mapstructure:"channel_codec_latency_factor"`
		Filepath                    string `mapstructure:"channel_filepath"`
		BannerGFXURL                string `mapstructure:"channel_banner_gfx_url"`
		NamePhonetic                string `mapstructure:"channel_name_phonetic"`
		Topic                       string `mapstructure:"channel_topic"`
		MaxFamilyClients            int    `mapstructure:"channel_maxfamilyclients"`
		IsSemiPermanent             bool   `mapstructure:"channel_flag_semi_permanent"`
		IsDefault                   bool   `mapstructure:"channel_flag_default"`
		SecuritySalt                string `mapstructure:"channel_security_salt"`
		UID                         string `mapstructure:"channel_unique_identifier"`
		MaxClientsIsUnlimited       bool   `mapstructure:"channel_flag_maxclients_unlimited"`
	}
	Clients []struct {
		UID              string `mapstructure:"client_unique_id"`
		Nickname         string `mapstructure:"client_nickname"`
		LastConnected    int64  `mapstructure:"client_lastconnected"`
		TotalConnections int    `mapstructure:"client_totalconnections"`
		Created          int64  `mapstructure:"client_created"`
		Description      string `mapstructure:"client_description"`
		ID               int    `mapstructure:"client_id"`
	}
	ServerGroups         []snapshotGroup
	ChannelGroups        []snapshotGroup
	ServerGroupRelations []struct {
		ClientDBID int `mapstructure:"cldbid"`
		GroupID    int `mapstructure:"gid"`
	}
	ChannelGroupRelations []struct {
		ChannelID  int `mapstructure:"iid"`
		ClientDBID int `mapstructure:"cldbid"`
		GroupID    int `mapstructure:"gid"`
	}
	ClientPermissions []struct {
		ClientDBID        int    `mapstructure:"id1"`
		PermissionID      string `mapstructure:"permid"`
		PermissionValue   int    `mapstructure:"permvalue"`
		PermissionSkip    bool   `mapstructure:"permskip"`
		PermissionNegated bool   `mapstructure:"permnegated"`
	}
	ChannelPermissions []struct {
		ChannelID         int    `mapstructure:"id1"`
		PermissionID      string `mapstructure:"permid"`
		PermissionValue   int    `mapstructure:"permvalue"`
		PermissionSkip    bool   `mapstructure:"permskip"`
		PermissionNegated bool   `mapstructure:"permnegated"`
	}
	ChannelClientPermissions []struct {
		ChannelID         int    `mapstructure:"id1"`
		ClientDBID        int    `mapstructure:"id2"`
		PermissionID      string `mapstructure:"permid"`
		PermissionValue   int    `mapstructure:"permvalue"`
		PermissionSkip    bool   `mapstructure:"permskip"`
		PermissionNegated bool   `mapstructure:"permnegated"`
	}
	APIKeys []struct {
		Hash      string `mapstructure:"hash"`
		ClientUID string `mapstructure:"cluid"`
		Scope     Scope  `mapstructure:"scope"`
		Created   int64  `mapstructure:"created_at"`
		Expires   int64  `mapstructure:"expires_at"`
	}
}

type snapshotGroup struct {
	ID          int    `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Permissions []struct {
		ID      string `mapstructure:"permid"`
		Value   int    `mapstructure:"permvalue"`
		Skip    bool   `mapstructure:"permskip"`
		Negated bool   `mapstructure:"permnegated"`
	}
}

// SnapshotCreate returns a newly created snapshot form server sid with password password
// sid - required
// password - optional
func (a Agent) SnapshotCreate(sid int, password string) (*RawSnapshot, error) {
	req := libts.Request{
		Command:  "serversnapshotcreate",
		ServerID: sid,
		Args: map[string]interface{}{
			"password": password,
		},
	}
	snapshot := &RawSnapshot{}
	err := a.Query.Do(req, snapshot)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

// SnapshotDeploy a snapshot snapshot with password password to server sid and keepfiles
// sid - required - 0 for new virtualserver
// keepFiles - optional
// password - optional
// snapshot - required
func (a Agent) SnapshotDeploy(sid int, keepFiles bool, password string, snapshot RawSnapshot) error {
	req := libts.Request{
		Command:  "serversnapshotdeploy",
		ServerID: sid,
		Args: map[string]interface{}{
			"password": password,
		},
	}
	if keepFiles {
		req.Args["-keepfiles"] = ""
	}
	req.Args["version"] = snapshot.Version
	req.Args["salt"] = snapshot.Salt
	req.Args["data"] = snapshot.Data
	return a.Query.Do(req, nil)
}

// DecodeToString decodes to snapshot to a string
// At the moment only supports non encrypted snapshots
func (s RawSnapshot) DecodeToString() (string, error) {
	// Base64 decode
	b64, err := base64.StdEncoding.DecodeString(s.Data)
	if err != nil {
		return "", err
	}
	// ZStandard decompression
	decomp, err := ioutil.ReadAll(zstd.NewReader(bytes.NewReader(b64)))
	if err != nil {
		return "", err
	}
	return string(decomp), err
}

// Decode a binary snapshot to a Snapshot struct
// At the moment only supports non encrypted snapshots
func (s RawSnapshot) Decode() (*Snapshot, error) {
	snapshot := &Snapshot{}
	raw, err := s.DecodeToString()
	if err != nil {
		return nil, err
	}
	// Virtualserver
	err = parseSection(selectSectionByIndeces(raw, 0, strings.Index(raw, endVirtualserver)), &snapshot.VirtualServer)
	if err != nil {
		return nil, err
	}
	// Channels
	err = parseSection(selectSection(raw, startChannels, endChannels, 0), &snapshot.Channels)
	if err != nil {
		return nil, err
	}
	// Clients
	err = parseSection(selectSection(raw, startClients, endClients, 0), &snapshot.Clients)
	if err != nil {
		return nil, err
	}
	// Groups
	// ServerGroups
	snapshot.ServerGroups, err = parseGroups(
		strings.Index(raw, startServerGroups),
		raw,
	)
	if err != nil {
		return nil, err
	}
	// ChannelGroups
	snapshot.ChannelGroups, err = parseGroups(
		strings.Index(raw, startChannelGroups),
		raw,
	)
	if err != nil {
		return nil, err
	}
	// Relations
	// ServerGroupRelations
	err = parseSection(selectSection(raw, "iid=0", endRelations, 0), &snapshot.ServerGroupRelations)
	if err != nil {
		return nil, err
	}
	// ChannelGroupRealations
	err = parseSection(
		selectSectionByIndeces(raw,
			strings.LastIndex(raw, endGroups),
			strings.LastIndex(raw, endRelations),
		), &snapshot.ChannelGroupRelations)
	if err != nil {
		return nil, err
	}
	// other permissions
	// ClientPermissions
	err = parseSection(selectSection(raw, startClientPermissions, endClientPermissions, 0), &snapshot.ClientPermissions)
	if err != nil {
		return nil, err
	}
	// ChannelPermissions
	err = parseSection(selectSection(raw, startChannelPermissions, endChannelPermissions, 0), &snapshot.ChannelPermissions)
	if err != nil {
		return nil, err
	}
	// ChannelClientPermissions
	err = parseSection(selectSection(raw, startChannelClientPermissions, endChannelClientPermissions, 0), &snapshot.ChannelClientPermissions)
	if err != nil {
		return nil, err
	}
	// APIKeys
	err = parseSection(selectSection(raw, startAPIKeys, endAPIKeys, 0), &snapshot.APIKeys)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func parseSection(section []byte, i interface{}) error {
	section = bytes.TrimSpace(section)
	if len(section) == 0 {
		return nil
	}
	m := communication.ConvertResponse(section)
	if len(m) != 0 {
		return communication.UnmarshalResponse(m, i)
	}
	return nil
}

func parseGroups(start int, raw string) ([]snapshotGroup, error) {
	groups := []snapshotGroup{}
	// the end of a group block is end_group|end_groups where the end_groups is the end for the last group
	end := strings.Index(raw[start:], fmt.Sprintf("%s|%s", endGroup, endGroups)) + start
	// we jump behind the next closing tag
	for groupStart := start; groupStart < end; groupStart = strings.Index(raw[groupStart:], endGroup) + len(endGroup) + groupStart {
		sg := &snapshotGroup{}
		// we can't parse the whole group at once so we split between the "meta" (id and name) and the permissions
		permissionStart := strings.Index(raw[groupStart:], "permid=")
		// first parse the meta
		err := parseSection(selectSectionByIndeces(raw, groupStart, groupStart+permissionStart), sg)
		if err != nil {
			return nil, err
		}
		// then parse the permissions
		err = parseSection(selectSectionByIndeces(raw,
			permissionStart,
			strings.Index(raw[groupStart:], endGroup)+groupStart,
		), &sg.Permissions)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *sg)
	}
	return groups, nil
}

func selectSection(text, start, end string, offset int) []byte {
	text = text[offset:]
	// start index: index of starttag -> place cursor behind tag
	startIndex := strings.Index(text, start) + len(start)
	endIndex := strings.Index(text[startIndex:], end) + startIndex
	return selectSectionByIndeces(text, startIndex, endIndex)
}

func selectSectionByIndeces(text string, start, end int) []byte {
	return []byte(text[start:end])
}
