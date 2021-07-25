package subscriber

import "github.com/schoeppi5/libts"

// TextMessageTarget represents one of three possible Targets for text messages
type TextMessageTarget string

const (
	// ChannelMessages -> recieve messages from the channel the query client is currently in
	ChannelMessages TextMessageTarget = TextMessageTarget(TextChannel)
	// ServerMessages -> recieve messages the query client is connected to
	// When the query client switches the virtualserver, the subscription has to be renewed
	ServerMessages TextMessageTarget = TextMessageTarget(TextServer)
	// PrivateMessages -> recieve messages send directly to the query client
	PrivateMessages TextMessageTarget = TextMessageTarget(TextPrivate)
)

// TextMessage subscribes to the textmessage event with the given target
func (a Agent) TextMessage(c chan interface{}, target TextMessageTarget) error {
	s := libts.Subscription{
		Name: string(target),
		Events: map[string]libts.Event{
			TextMessage: {
				Template: &TextMessageEvent{},
				C:        c,
			},
		},
	}
	err := a.Subscriber.Subscribe(s)
	if err != nil {
		return err
	}
	return nil
}
