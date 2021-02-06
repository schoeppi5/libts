package query

import (
	"math/big"
	"net"

	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/communication"
)

// HostInfo contains the info about the TeamSpeak 3 server host
type HostInfo struct {
	MaxClients                     int     `mapstructure:"virtualservers_total_maxclients"`
	FileTransferBandwithSent       big.Int `mapstructure:"connection_filetransfer_bandwidth_sent"`
	FileTransferBytesReceivedTotal big.Int `mapstructure:"connection_filetransfer_bytes_received_total"`
	SentLastSecondTotal            big.Int `mapstructure:"connection_bandwidth_sent_last_second_total"`
	ReceivedLastSecondTotal        big.Int `mapstructure:"connection_bandwidth_received_last_second_total"`
	ClientsOnlineTotal             int     `mapstructure:"virtualservers_total_clients_online"`
	ChannelsOnlineTotal            int     `mapstructure:"virtualservers_total_channels_online"`
	Uptime                         int     `mapstructure:"instance_uptime"`
	FileTransferBandwidthReceived  big.Int `mapstructure:"connection_filetransfer_bandwidth_received"`
	FileTransferBytesSentTotal     big.Int `mapstructure:"connection_filetransfer_bytes_sent_total"`
	BytesSendTotal                 big.Int `mapstructure:"connection_bytes_sent_total"`
	ReceivedLastMinuteTotal        big.Int `mapstructure:"connection_bandwidth_received_last_minute_total"`
	SentLastMinuteTotal            big.Int `mapstructure:"connection_bandwidth_sent_last_minute_total"`
	TimestampUTC                   int64   `mapstructure:"host_timestamp_utc"`
	VirtualserverRunningTotal      int     `mapstructure:"virtualservers_running_total"`
	PacketsSentTotal               big.Int `mapstructure:"connection_packets_sent_total"`
	PacketsReceivedTotal           big.Int `mapstructure:"connection_packets_received_total"`
	BytesReceivedTotal             big.Int `mapstructure:"connection_bytes_received_total"`
}

// InstanceInfo contains info about the TeamSpeak 3 server instance
type InstanceInfo struct {
	FileTransterPort               int     `mapstructure:"serverinstance_filetransfer_port"`
	MaxDownloadBandwidthTotal      big.Int `mapstructure:"serverinstance_max_download_total_bandwidth"`
	ServerQueryFloodCommands       int     `mapstructure:"serverinstance_serverquery_flood_commands"`
	ChannelDefaultGroupTemplate    int     `mapstructure:"serverinstance_template_channeldefault_group"`
	PendingConnectionsPerIP        int     `mapstructure:"serverinstance_pending_connections_per_ip"`
	ServerQueryMaxConnectionsPerIP int     `mapstructure:"serverinstance_serverquery_max_connections_per_ip"`
	ServerQueryGuestGroup          int     `mapstructure:"serverinstance_guest_serverquery_group"`
	ServerQueryFloodTime           int     `mapstructure:"serverinstance_serverquery_flood_time"`
	ServerDefaultGroupTemplate     int     `mapstructure:"serverinstance_template_serverdefault_group"`
	ChannelAdminGroupTemplate      int     `mapstructure:"serverinstance_template_channeladmin_group"`
	PermissionsVersion             string  `mapstructure:"serverinstance_permissions_version"`
	DatabaseVersion                string  `mapstructure:"serverinstance_database_version"`
	MaxUploadTotalBandwidth        big.Int `mapstructure:"serverinstance_max_upload_total_bandwidth"`
	ServerQueryBanTime             int     `mapstructure:"serverinstance_serverquery_ban_time"`
	ServerAdminGroupTemplate       int     `mapstructure:"serverinstance_template_serveradmin_group"`
}

// Host retrieves info about the host of the TeamSpeak server
func (a Agent) Host() (*HostInfo, error) {
	req := libts.Request{
		Command: "hostinfo",
	}
	host := &HostInfo{}
	err := a.Query.Do(req, host)
	if err != nil {
		return nil, err
	}
	return host, nil
}

// Instance returns the combined result of hostinfo and instanceinfo since I didn't really see point in dividing them
func (a Agent) Instance() (*InstanceInfo, error) {
	instance := &InstanceInfo{}
	err := a.Query.Do(
		libts.Request{
			Command: "instanceinfo",
		}, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// InstanceEdit changes the instance according to instance
// instance - required - be aware that also null values are used. e.g.: if FileTransferPort is left undeclared in the struct it will have 0 as its value and the function will attempt to change the port to 0
// So it is advice to first get the InstanceInfo using Instance(), then change what you want to change and send it back
func (a Agent) InstanceEdit(instance InstanceInfo) error {
	args, err := communication.MarshalRequest(instance)
	if err != nil {
		return err
	}
	req := libts.Request{
		Command: "instanceedit",
		Args:    args,
	}
	return a.Query.Do(req, nil)
}

// BindingList returns the bindings for the specified subsystem
// sid - required
// subsystem - optional - Default 'voice' - Possible 'voice', 'query', 'filetransfer'
func (a Agent) BindingList(sid int, subsystem string) ([]net.IP, error) {
	req := libts.Request{
		Command:  "bindinglist",
		ServerID: sid,
		Args:     map[string]interface{}{},
	}
	if subsystem != "" {
		req.Args["subsystem"] = subsystem
	}
	ips := []struct {
		IP string `mapstructure:"ip"`
	}{}
	err := a.Query.Do(req, &ips)
	if err != nil {
		return nil, err
	}
	ipList := []net.IP{}
	for i := range ips {
		ipList = append(ipList, net.ParseIP(ips[i].IP))
	}
	return ipList, nil
}

// GlobalMessage sends message message to virtualservers in the server chat (using the gm command - I just don't know, why this command has such a non descriptive name)
// message - required
func (a Agent) GlobalMessage(message string) error {
	req := libts.Request{
		Command: "gm",
		Args: map[string]interface{}{
			"msg": message,
		},
	}
	return a.Query.Do(req, nil)
}
