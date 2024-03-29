package subscriber

import "github.com/schoeppi5/libts"

// ChannelCreated subscribes to the ChannelCreated events for all channels
// events recieved on the channel will always be of type ChannelCreatedEvent
func (a Agent) ChannelCreated(c chan interface{}) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: 0,
		Events: map[string]libts.Event{
			ChannelCreated: {
				Template: &ChannelCreatedEvent{},
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

// ChannelDeleted subscribes to the ChannelDeleted events for channel cid (0 for all)
// events recieved on the channel will always be of type ChannelDeletedEvent
func (a Agent) ChannelDeleted(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ChannelDeleted: {
				Template: &ChannelDeletedEvent{},
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

// ChannelMoved subscribes to the ChannelMoved events for channel cid (0 for all)
// events recieved on the channel will always be of type ChannelMovedEvent
func (a Agent) ChannelMoved(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ChannelMoved: {
				Template: &ChannelMovedEvent{},
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

// ChannelEdited subscribes to the ChannelEdited events for channel cid (0 for all)
// events recieved on the channel will always be of type ChannelEditedEvent
func (a Agent) ChannelEdited(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ChannelEdited: {
				Template: &ChannelEditedEvent{},
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

// ChannelDescriptionChanged subscribes to the ChannelDescriptionChanged events for channel cid (0 for all)
// events recieved on the channel will always be of type ChannelDescriptionChangedEvent
func (a Agent) ChannelDescriptionChanged(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ChannelDescriptionChanged: {
				Template: &ChannelDescriptionChangedEvent{},
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

// ChannelPasswordChanged subscribes to the ChannelPasswordChanged events for channel cid (0 for all)
// events recieved on the channel will always be of type ChannelPasswordChangedEvent
// Is only omited when password is addded/removed not when actually changed
func (a Agent) ChannelPasswordChanged(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ChannelPasswordChanged: {
				Template: &ChannelPasswordChangedEvent{},
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

// ClientMoved subscribes to the ClientMoved events for channel cid (0 for all)
// events recieved on the channel will always be of type ClientMovedEvent
func (a Agent) ClientMoved(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ClientMoved: {
				Template: &ClientMovedEvent{},
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

// ClientJoinedChannel subscribes to the ClientEnterView events for channel cid (0 for all)
// events recieved on the channel will always be of type ClientEnterViewEvent
func (a Agent) ClientJoinedChannel(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ClientEnterView: {
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

// ClientLeftChannel subscribes to the ClientLeftView events for channel cid (0 for all)
// events recieved on the channel will always be of type ClientLeftViewEvent
func (a Agent) ClientLeftChannel(c chan interface{}, cid int) error {
	s := libts.Subscription{
		Name:      Channel,
		ChannelID: cid,
		Events: map[string]libts.Event{
			ClientLeftView: {
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
