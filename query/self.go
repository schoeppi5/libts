package query

import "github.com/schoeppi5/libts"

// WhoAmI returns the client struct for the serverquery
func (a Agent) WhoAmI(sid int) (*Client, error) {
	whoami := libts.Request{
		ServerID: sid,
		Command:  "whoami",
	}
	i := struct {
		ID int `mapstructure:"client_id"`
	}{}
	err := a.Query.Do(whoami, &i)
	if err != nil {
		return nil, err
	}
	return a.Client(sid, i.ID)
}
