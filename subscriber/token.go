package subscriber

import "github.com/schoeppi5/libts"

// TokenUsed subscribes to the TokenUsed events
// events recieved on the channel will always be of type TokenUsedEvent
func (a Agent) TokenUsed(c chan interface{}) error {
	s := libts.Subscription{
		Name: Token,
		Events: map[string]libts.Event{
			TokenUsed: libts.Event{
				Template: &TokenUsedEvent{},
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
