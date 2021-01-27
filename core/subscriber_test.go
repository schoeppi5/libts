package core_test

import (
	"reflect"
	"testing"

	"github.com/schoeppi5/libts/core"

	"github.com/schoeppi5/libts"
)

func TestSortEvents(t *testing.T) {
	// given
	c := make(chan interface{})
	in := make(chan []byte, 1)
	defer close(in)
	in <- []byte("notifytest test1=1 test2=2")
	type tmp struct {
		Test1 int `mapstructure:"test1"`
		Test2 int `mapstructure:"test2"`
	}
	want := &tmp{
		Test1: 1,
		Test2: 2,
	}
	event := &libts.Event{
		C:        c,
		Template: &tmp{},
	}
	es := core.NewEventStore()
	es.Add("notifytest", event)

	// when
	go core.SortEvents(in, es)

	// then
	notify := <-c
	if _, ok := notify.(*tmp); !ok {
		LogTestError(reflect.TypeOf(notify), "*core_test.tmp", t)
	}
	if e := notify.(*tmp); *e != *want {
		LogTestError(e, want, t)
	}
}

func TestSortEventsNonPointerTemplate(t *testing.T) {
	// given
	c := make(chan interface{})
	in := make(chan []byte, 1)
	defer close(in)
	in <- []byte("notifytest test1=1 test2=2")
	type tmp struct {
		Test1 int `mapstructure:"test1"`
		Test2 int `mapstructure:"test2"`
	}
	event := &libts.Event{
		C:        c,
		Template: tmp{}, // non pointer
	}
	es := core.NewEventStore()
	es.Add("notifytest", event)

	// when
	go core.SortEvents(in, es)

	// then
	notify := <-c
	if _, ok := notify.(error); !ok {
		LogTestError(reflect.TypeOf(c), "error", t)
	}
	if e := notify.(error); e.Error() != "expected pointer to value, not value" {
		LogTestError(e, "expected pointer to value, not value", t)
	}
}

func TestSortEventsClosedChannel(t *testing.T) {
	t.Parallel()
	// given
	in := make(chan []byte)
	close(in)

	// when
	core.SortEvents(in, core.NewEventStore())

	// then
	// passes when not killed by deadline
}

// TODO: BasicSubscriber test (needs query Mock)
