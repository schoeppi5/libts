package subscriber

const (
	// Subscriptions

	// Server registers the following events:
	// `cliententerview`, `clientleftview`, `serveredited`.
	Server string = "server"

	// Channel registers the following events:
	// `cliententerview`, `clientleftview`, `channeldescriptionchanged`, `channelpasswordchanged`
	// `channelmoved`, `channeledited`, `channelcreated`, `channeldeleted`, `clientmoved`.
	Channel string = "channel"

	// TextServer registers the `textmessage` event with `targetmode = 3`.
	TextServer string = "textserver"

	// TextChannel registers the `textmessage` event with `targetmode = 2`.
	//
	// Notifications are only received for messages that are sent in the channel that the client is in.
	TextChannel string = "textchannel"

	// TextPrivate registers the `textmessage` event with `targetmode = 1`.
	TextPrivate string = "textprivate"

	// Token registers the `tokenused` event.
	Token string = "tokenused"

	// Events

	// ServerEdited -> server properties were changed
	ServerEdited string = "notifyserveredited"
	// ClientLeftView -> client left the server/channel
	ClientLeftView string = "notifyclientleftview"
	// ClientEnterView -> client entered server/channel
	ClientEnterView string = "notifycliententerview"
	// ChannelDescriptionChanged -> description of the channel was changed/removed
	ChannelDescriptionChanged string = "notifychanneldescriptionchanged"
	// ChannelPasswordChanged -> password of the channel was set/removed (not emited when changed)
	ChannelPasswordChanged string = "notifychannelpasswordchanged"
	// ChannelMoved -> channel was moved
	ChannelMoved string = "notifychannelmoved"
	// ChannelEdited -> channel properties were changed
	ChannelEdited string = "notifychanneledited"
	// ChannelCreated -> channel was created (only emited when all channel are subscribed to (cid 0))
	ChannelCreated string = "notifychannelcreated"
	// ChannelDeleted -> channel was deleted
	ChannelDeleted string = "notifychanneldeleted"
	// ClientMoved -> client was moved
	ClientMoved string = "notifyclientmoved"
	// TextMessage -> the Query recieved a textmessage
	TextMessage string = "notifytextmessage"
	// TokenUsed -> someone used a token
	TokenUsed string = "notifytokenused"
)
