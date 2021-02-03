package query

import (
	"time"

	"github.com/schoeppi5/libts"
)

// Scope of an APIKey
type Scope string

// UnmarshalText fulfills the TextUnmarshaler interface
func (s *Scope) UnmarshalText(text []byte) error {
	*s = Scope(text)
	return nil
}

const (
	// ManageScope includes all commands
	ManageScope Scope = Scope("manage")
	// WriteScope includes more commands than the read, but fewer than the manage scope
	// See an overview [here]("https://community.teamspeak.com/t/webquery-discussion-help-3-12-0-onwards/7184")
	WriteScope Scope = Scope("write")
	// ReadScope is the least permissive scope
	// See an overview [here]("https://community.teamspeak.com/t/webquery-discussion-help-3-12-0-onwards/7184")
	ReadScope Scope = Scope("read")
)

// APIKey represents the output of ever APIKeyAdd or APIKeyList
// If this is a result of APIKeyList Key is empty
type APIKey struct {
	Key        string `mapstructure:"apikey"`
	ID         int    `mapstructure:"id"`
	ServerID   int    `mapstructure:"sid"`
	ClientDBID int    `mapstructure:"cldbid"`
	Scope      Scope  `mapstructure:"scope"`
	Created    int64  `mapstructure:"created_at"`
	Expires    int64  `mapstructure:"expires_at"`
}

// APIKeyAdd a new key on server sid for user cldbid with the scope s and lifetime in days
// sid - required
// scope - required
// lifetime - optional - Default is 14 days (-1)
// cldbid - optional - Default is invoker (0)
func (a Agent) APIKeyAdd(sid int, s Scope, lifetime int, cldbid int) (*APIKey, error) {
	if lifetime == -1 {
		lifetime = 14
	}
	req := libts.Request{
		Command:  "apikeyadd",
		ServerID: sid,
		Args: map[string]interface{}{
			"scope":    s,
			"lifetime": lifetime,
		},
	}
	if cldbid != 0 {
		req.Args["cldbid"] = cldbid
	}
	apiKey := &APIKey{}
	err := a.Query.Do(req, apiKey)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	apiKey.Created = now.Unix()
	apiKey.Expires = now.Add((24 * time.Hour) * time.Duration(lifetime)).Unix()
	return apiKey, nil
}

// APIKeyDel the key id from server sid
// sid - required
// id - required
func (a Agent) APIKeyDel(sid int, id int) error {
	req := libts.Request{
		Command:  "apikeydel",
		ServerID: sid,
		Args: map[string]interface{}{
			"id": id,
		},
	}
	return a.Query.Do(req, nil)
}

// APIKeyList all keys on server sid for client cldbid. Start at offset and return count keys total
// sid - required
// cldbid - optional - Default is invoker (0) - (-1) for all clients
// offset - optional
// count - optional - Default is all (0)
func (a Agent) APIKeyList(sid int, cldbid int, offset int, count int) ([]APIKey, error) {
	req := libts.Request{
		Command:  "apikeylist",
		ServerID: sid,
		Args: map[string]interface{}{
			"start": offset,
		},
	}
	if cldbid == -1 {
		req.Args["cldbid"] = "*"
	} else if cldbid != 0 {
		req.Args["cldbid"] = cldbid
	}
	if count != 0 {
		req.Args["duration"] = count
	}
	println(req.String())
	apiKeys := []APIKey{}
	err := a.Query.Do(req, &apiKeys)
	if err != nil {
		return nil, err
	}
	return apiKeys, nil
}
