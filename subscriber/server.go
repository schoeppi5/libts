package subscriber

import "github.com/schoeppi5/libts"

// ServerEdited subscribes to the ServerEdited event
// events recieved will be of type ServerEditedEvent
func (a Agent) ServerEdited(c chan interface{}) error {
	s := libts.Subscription{
		Name: Server,
		Events: map[string]libts.Event{
			ServerEdited: libts.Event{
				Template: &ServerEditedEvent{},
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

// ClientJoinedServer subscribes to the ClientEnterView events on server level
// events recieved on the channel will always be of type ClientEnterViewEvent
func (a Agent) ClientJoinedServer(c chan interface{}) error {
	s := libts.Subscription{
		Name: Server,
		Events: map[string]libts.Event{
			ClientEnterView: libts.Event{
				Template: &ClientEnterViewEvent{},
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

// ClientLeftServer subscribes to the ClientLeftView events on server level
// events recieved on the channel will always be of type ClientLeftViewEvent
func (a Agent) ClientLeftServer(c chan interface{}) error {
	s := libts.Subscription{
		Name: Server,
		Events: map[string]libts.Event{
			ClientLeftView: libts.Event{
				Template: &ClientLeftViewEvent{},
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
