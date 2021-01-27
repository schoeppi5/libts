package core

import (
	"fmt"
	"testing"

	"github.com/schoeppi5/libts"
)

func TestAdd(t *testing.T) {
	// given
	es := NewEventStore()
	key := "test"
	value := &libts.Event{
		C:        nil,
		Template: "string",
	}

	// when
	es.Add(key, value)

	// then
	if es.store[key] != value {
		t.Errorf("Test %s failed!\n\tHave: %v\n\tWant: %v", t.Name(), es.store[key], value)
	}
}

func TestDel(t *testing.T) {
	// given
	es := NewEventStore()
	key := "test"
	value := &libts.Event{
		C:        nil,
		Template: nil,
	}
	es.Add(key, value)

	// when
	es.Del(key)

	// then
	if _, ok := es.store[key]; ok {
		t.Errorf("Test %s failed!\n\tKey %s still present in store", t.Name(), key)
	}
}

func TestGetPresent(t *testing.T) {
	// given
	es := NewEventStore()
	key := "test"
	value := &libts.Event{
		C:        nil,
		Template: nil,
	}
	es.Add(key, value)

	// when
	v, ok := es.Get(key)

	// then
	if !ok {
		t.Errorf("Test %s failed!\n\tKey %s not present in store", t.Name(), key)
	}
	if v != value {
		t.Errorf("Test %s failed!\n\tHave: %v\n\tWant: %v", t.Name(), v, value)
	}
}

func TestGetNotPresent(t *testing.T) {
	// given
	es := NewEventStore()
	key := "test"

	// when
	v, ok := es.Get(key)

	// then
	if ok {
		t.Errorf("Test %s failed!\n\tFound value for key %s: %v", t.Name(), key, v)
	}
}

func TestLen(t *testing.T) {
	// given
	es := NewEventStore()
	keyCount := 10
	for i := 0; i < keyCount; i++ {
		es.Add(fmt.Sprintf("%d", i), nil)
	}

	// when
	l := es.Len()

	if l != keyCount {
		t.Errorf("Test %s failed!\n\tHave: %d\n\tWant: %d", t.Name(), l, keyCount)
	}
}

func TestClean(t *testing.T) {
	// given
	es := NewEventStore()
	es.Add("test", nil)

	// when
	es.Clean()

	// then
	if len(es.store) != 0 {
		t.Errorf("Test %s failed!\n\tHave: %d\n\tWant: %d", t.Name(), len(es.store), 0)
	}
}

func TestKeys(t *testing.T) {
	// given
	es := NewEventStore()
	es.Add("test", nil)

	// when
	keys := es.Keys()

	// then
	if len(keys) != 1 {
		t.Errorf("Test %s failed!\n\tHave: %d\n\tWant: %d", t.Name(), len(keys), 1)
	}
	if keys[0] != "test" {
		t.Errorf("Test %s failed!\n\tHave: %s\n\tWant: %s", t.Name(), keys[0], "test")
	}
}
