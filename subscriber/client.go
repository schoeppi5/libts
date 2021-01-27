package subscriber

import "github.com/schoeppi5/libts"

// PrivateTextMessages subscribes to the TextMessage events send to the client directly
// events recieved on the channel will always be of type TextMessageEvent
func (a Agent) PrivateTextMessages(c chan interface{}) error {
	s := libts.Subscription{
		Name: TextPrivate,
		Events: map[string]libts.Event{
			TextMessage: libts.Event{
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
