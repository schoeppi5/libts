package query

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/communication"
)

const (
	// LevelError equals 1
	LevelError LogLevel = 1
	// LevelWarning equals 2
	LevelWarning = LevelError + 1
	// LevelDebug equals 3
	LevelDebug = LevelWarning + 1
	// LevelInfo equals 4
	LevelInfo = LevelDebug + 1
)

// LogLevel represents one of the four log levels teamspeak uses
type LogLevel int

// UnmarshalText used for decoding
func (ll *LogLevel) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "error":
		*ll = LevelError
	case "warning":
		*ll = LevelWarning
	case "debug":
		*ll = LevelDebug
	case "info":
		*ll = LevelInfo
	default:
		return fmt.Errorf("invalid loglevel %s", text)
	}
	return nil
}

// MarshalText used for encoding
func (ll LogLevel) MarshalText() ([]byte, error) {
	switch ll {
	case 1:
		return []byte("ERROR"), nil
	case 2:
		return []byte("WARNING"), nil
	case 3:
		return []byte("DEBUG"), nil
	case 4:
		return []byte("INFO"), nil
	default:
		return nil, fmt.Errorf("invalid loglevel %d", ll)
	}
}

// LogEntry represents one line of log
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Logger    string
	Filed3    string // I don't know if this field is just spacing or whatever
	Message   string
}

// LogView returns either the instances logs (sid 0) or the server logs
// It always returns the whole log, what could take a while
// TODO: Pageination
// There is currently a bug where the first three bytes of the first line of a log file are not returned, thus the timestamp of the first logentry is empty
// sid - optional
func (a Agent) LogView(sid int) ([]LogEntry, error) {
	log := []LogEntry{}
	req := libts.Request{
		Command: "logview",
		Args: map[string]interface{}{
			"lines":   99,
			"reverse": 1,
		},
	}
	if sid == 0 {
		req.Args["instance"] = 1
	} else {
		req.ServerID = sid
	}
	for beginPos := int64(1); beginPos != 0; { // keep querying the log until there is no more to query
		if beginPos > 1 { // beginPos needed to start somewhere
			req.Args["begin_pos"] = beginPos
		}
		logs, err := a.Query.DoRaw(req)
		if err != nil {
			return nil, err
		}
		meta := struct {
			LastPos  int64 `mapstructure:"last_pos"`
			FileSize int64 `mapstructure:"file_size"`
		}{}
		firstLog := bytes.Index(logs, []byte("l="))
		err = communication.UnmarshalResponse(communication.ConvertResponse(logs[:firstLog]), &meta)
		if err != nil {
			return nil, err
		}
		beginPos = meta.LastPos
		entries := []struct {
			L string `mapstructure:"l"`
		}{}
		err = communication.UnmarshalResponse(communication.ConvertResponse(logs[firstLog:]), &entries)
		if err != nil {
			return nil, err
		}
		for i := range entries { // parse the log entries
			entry, err := parseLogEntry(entries[i].L)
			if err != nil {
				return nil, err
			}
			log = append(log, *entry)
		}
	}
	for left, right := 0, len(log)-1; left < right; left, right = left+1, right-1 { // reverse the array, so the first line in the log is actually first
		log[left], log[right] = log[right], log[left]
	}
	return log, nil
}

// LogAdd adds a custom log message to the virtualserver log. You cannot log to the instance log
// sid - required
// level - required
// message - required
func (a Agent) LogAdd(sid int, level LogLevel, message string) error {
	req := libts.Request{
		Command:  "logadd",
		ServerID: sid,
		Args: map[string]interface{}{
			"logmsg":   message,
			"loglevel": int(level),
		},
	}
	return a.Query.Do(req, nil)
}

func parseLogEntry(l string) (*LogEntry, error) {
	sections := strings.Split(l, "|")
	if len(sections) != 5 {
		return nil, fmt.Errorf("failed to parse log entry: invalid number of log segments %d", len(sections))
	}
	le := &LogEntry{}
	err := le.Level.UnmarshalText(bytes.TrimSpace([]byte(sections[1])))
	if err != nil {
		return nil, err
	}
	le.Timestamp, err = time.Parse("2006-01-02 15:04:05.000000", sections[0])
	if err != nil {
		le.Timestamp = time.Time{} // workaround for the above mentioned bug
		// return nil, err
	}
	le.Logger = strings.TrimSpace(sections[2])
	le.Filed3 = strings.TrimSpace(sections[3])
	le.Message = strings.TrimSpace(sections[4])
	return le, nil
}
