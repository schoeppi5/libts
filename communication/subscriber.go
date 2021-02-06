package communication

import (
	"bytes"

	"github.com/schoeppi5/libts"
)

// shared codebase for subscriber

// SortEvents to the appropriate channels
// exits when in is closed
// either sends the event as event.Template or an error when failed to parse
// SortEvents returns if es is nil
func SortEvents(in <-chan []byte, es *EventStore) {
	if es == nil {
		return
	}
	for {
		notify, open := <-in
		if !open {
			for _, v := range es.Keys() { // close all listening channels
				e, _ := es.Get(v)
				close(e.C)
			}
			return
		}
		n := bytes.SplitN(notify, []byte(" "), 2) // 0 -> type | 1 -> notification data
		if e, ok := es.Get(string(n[0])); ok {
			err := UnmarshalResponse(ConvertResponse(n[1]), e.Template)
			if err != nil { // If an error occures, discard event
				e.C <- err
				continue
			}
			e.C <- e.Template
		}
	}
}

// BasicSubscriber implements the libts.Subscriber interface
type BasicSubscriber struct {
	query         libts.Query
	subscriptions *EventStore
	serverID      int
}

// NewBasicSubscriber takes a Query and adds Subscriber capabilities
func NewBasicSubscriber(q libts.Query, serverID int) *BasicSubscriber {
	bs := &BasicSubscriber{
		query:         q,
		subscriptions: NewEventStore(),
		serverID:      serverID,
	}
	go SortEvents(bs.query.Notification(), bs.subscriptions)
	return bs
}

// Subscribe to s and attempt to parse the responses
func (bs *BasicSubscriber) Subscribe(s libts.Subscription) error {
	return bs.register(s)
}

// register sets up the subscription
func (bs *BasicSubscriber) register(s libts.Subscription) error {
	// create command for subscription
	r := libts.Request{
		ServerID: bs.serverID,
		Command:  "servernotifyregister",
		Args: map[string]interface{}{
			"event": s.Name,
		},
	}
	if s.Name == "channel" {
		r.Args["id"] = s.ChannelID
	}
	// run command
	_, err := bs.query.DoRaw(r)
	if err != nil {
		return err
	}
	for i, v := range s.Events {
		bs.subscriptions.Add(i, &v) // override existing subscriptions and add new ones
	}
	return nil
}

// Unsubscribe from e
func (bs *BasicSubscriber) Unsubscribe(e string) {
	bs.subscriptions.Del(e) // delete e from map
}

// UnsubscribeAll removes all current subscriptions
func (bs *BasicSubscriber) UnsubscribeAll() error {
	_, err := bs.query.DoRaw(
		libts.Request{
			ServerID: bs.serverID,
			Command:  "servernotifyunregister",
		},
	)
	bs.serverID = 0
	bs.subscriptions.Clean()
	return err
}
