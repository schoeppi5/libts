package webquery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/core"
)

// This file fullfills the libts.Query interface

type response struct {
	Body   []map[string]interface{}
	Status core.QueryError
}

// Do executes the given command against TeamSpeak
func (wq WebQuery) Do(request libts.Request, response interface{}) error {
	respBody, err := wq.DoRaw(request)
	err = unmarshalBody(respBody, response)
	if err != nil {
		return err
	}
	return nil
}

// DoRaw executes the given command against TeamSpeak and returns the unformated output
func (wq WebQuery) DoRaw(request libts.Request) ([]byte, error) {
	resp, err := wq.HTTPClient.Do(wq.marshalRequest(request))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// Notification is not yet implementable using webquery! This func will panic if used
func (wq WebQuery) Notification() <-chan []byte {
	panic(errors.New("webquery does not yet implement subscribing to events"))
}

func unmarshalBody(body []byte, value interface{}) error {
	response := &response{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Status.ID != 0 {
		return response.Status
	}
	// decode query response
	for i := range response.Body {
		for key := range response.Body[i] {
			if s, ok := response.Body[i][key].(string); ok {
				response.Body[i][key] = libts.QueryDecoder.Replace(s)
			}
		}
	}
	return core.UnmarshalResponse(response.Body, value)
}

func (wq WebQuery) marshalRequest(r libts.Request) *http.Request {
	var url string
	if r.ServerID != 0 {
		url = wq.url(fmt.Sprintf("%d/%s", r.ServerID, r.Command))
	} else {
		url = wq.url(r.Command)
	}
	if r.Args != nil {
		body, _ := json.Marshal(r.Args)
		req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
		return req
	}
	req, _ := http.NewRequest("GET", url, nil)
	return req
}

func (wq WebQuery) url(c string) string {
	if wq.TLS {
		return fmt.Sprintf("https://%s:%d/%s?api-key=%s", wq.Host, wq.Port, c, wq.Key)
	}
	return fmt.Sprintf("http://%s:%d/%s?api-key=%s", wq.Host, wq.Port, c, wq.Key)
}
