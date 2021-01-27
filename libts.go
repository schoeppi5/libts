package libts

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	// QueryEncoder escapes special characters as needed
	QueryEncoder = strings.NewReplacer(
		`\`, `\\`,
		`/`, `\/`,
		` `, `\s`,
		`|`, `\p`,
		"\a", `\a`,
		"\b", `\b`,
		"\f", `\f`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
		"\v", `\v`,
	)

	// QueryDecoder replaces escaped characters with thier string representation
	QueryDecoder = strings.NewReplacer(
		`\\`, "\\",
		`\/`, "/",
		`\s`, " ",
		`\p`, "|",
		`\a`, "\a",
		`\b`, "\b",
		`\f`, "\f",
		`\n`, "\n",
		`\r`, "\r",
		`\t`, "\t",
		`\v`, "\v",
	)
)

// Request contains all necessary infos for a query request towards teamspeak
// You need to know, what to expect as an answer
type Request struct {
	ServerID int
	Command  string
	Args     map[string]interface{}
}

// String returns the correct string representation for a serverquery
func (r Request) String() string {
	c := r.Command
	for i, v := range r.Args {
		vt := reflect.ValueOf(v)
		if vt.Kind() == reflect.Array || vt.Kind() == reflect.Slice {
			for j := 0; j < vt.Len(); j++ {
				if j > 0 {
					c += "|"
				} else {
					c += " "
				}
				c += printArg(i, vt.Index(j))
			}
		} else {
			c += fmt.Sprintf(" %s", printArg(i, v))
		}
	}
	c = c + "\n"
	return c
}

func printArg(key string, value interface{}) string {
	return fmt.Sprintf("%s=%s", key, QueryEncoder.Replace(fmt.Sprint(value)))
}

// Subscription describes an event a subscriber can subscribe to
// Name is the Name of the "Event" you are subscribing to (server, channel, etc.)
// Events are the events with the type that represents them you want to recieve. All other are discarded
// The subscriber will attempt to parse the event to the type of the channel
// ChannelID is ignored if Subscription is not libts.ChannelEvents
type Subscription struct {
	ChannelID int
	Name      string
	Events    map[string]Event
}

// Event describes the structure and the channel for a single event (cliententerview, clientleftview, etc.)
type Event struct {
	Template interface{}
	C        chan<- interface{}
}

// Query is the interface for all implementation (webquery, serverquery, more?)
// Be aware, that different queries can do different things (e.g. serverquery can recieve notifications, webquery can't atm)
type Query interface {
	// Do executes a given request against teamspeak and tries to parse the answer in the given interface
	// You can give it either one single PTR to a struct or a PTR to a slice if you expect to recieve more than one answer
	Do(req Request, res interface{}) error
	// DoRaw just returns the answer of teamspeak
	DoRaw(req Request) ([]byte, error)
	// Notifications provides a io.Reader which includes only the notifications
	Notification() <-chan []byte
}

// Subscriber can subscribe to events
type Subscriber interface {
	// Subscribe to s and attempt to parse the response
	Subscribe(s Subscription) error
	// Unsubscribe from notification n
	Unsubscribe(n string)
	// UnsubscribeAll form all subscriptions
	UnsubscribeAll() error
}
