package communication

import (
	"sync"

	"github.com/schoeppi5/libts"
)

// EventStore is used to manage events across a subscriber and the event loop (communication.SortEvents)
// it is basically a thread-safe map
type EventStore struct {
	store map[string]*libts.Event
	lock  sync.Locker
}

// NewEventStore retuns a thread safe EventStore
func NewEventStore() *EventStore {
	return &EventStore{
		store: make(map[string]*libts.Event),
		lock:  &sync.Mutex{},
	}
}

// Add adds/updates name event pairs to/in the store
func (es *EventStore) Add(name string, event *libts.Event) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.store[name] = event
}

// Del deletes name event pairs from the store
func (es *EventStore) Del(name string) {
	es.lock.Lock()
	defer es.lock.Unlock()
	delete(es.store, name)
}

// Get returns the event to a name
func (es *EventStore) Get(name string) (*libts.Event, bool) {
	es.lock.Lock()
	defer es.lock.Unlock()
	e, ok := es.store[name]
	return e, ok
}

// Clean the whole store
func (es *EventStore) Clean() {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.store = make(map[string]*libts.Event)
}

// Len returns the number of events in the store
func (es *EventStore) Len() int {
	es.lock.Lock()
	defer es.lock.Unlock()
	return len(es.store)
}

// Keys all keys in store
func (es *EventStore) Keys() (keys []string) {
	es.lock.Lock()
	defer es.lock.Unlock()
	for key := range es.store {
		keys = append(keys, key)
	}
	return
}
