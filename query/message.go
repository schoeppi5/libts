package query

import "github.com/schoeppi5/libts"

// SendOfflineMessage on server sid to client cluid with header subject and body message
// sid - required
// cluid - required
// subject - required
// message - required
func (a Agent) SendOfflineMessage(sid int, cluid, subject, message string) error {
	req := libts.Request{
		Command:  "messageadd",
		ServerID: sid,
		Args: map[string]interface{}{
			"cluid":   cluid,
			"subject": subject,
			"message": message,
		},
	}
	return a.Query.Do(req, nil)
}

// DeleteOfflineMessage on server sid with message id id from your inbox
// I don't really know why this functionality even exists. I wasn't able to figure out, how to even send offline messages to the server query
// sid - required
// id - required
func (a Agent) DeleteOfflineMessage(sid int, id int) error {
	req := libts.Request{
		Command:  "messagedel",
		ServerID: sid,
		Args: map[string]interface{}{
			"msgid": id,
		},
	}
	return a.Query.Do(req, nil)
}

// GetOfflineMessage on server sid with message id id from your inbox
// I don't really know why this functionality even exists. I wasn't able to figure out, how to even send offline messages to the server query
// sid - required
// id - required
func (a Agent) GetOfflineMessage(sid, id int) error {
	req := libts.Request{
		Command:  "messageget",
		ServerID: sid,
		Args: map[string]interface{}{
			"msgid": id,
		},
	}
	return a.Query.Do(req, nil)
}

// Not implemented
// I don't really know why this functionality even exists. I wasn't able to figure out, how to even send offline messages to the server query
// func (a Agent) ListOfflineMessages(sid int)

// Not implemented
// I don't really know why this functionality even exists. I wasn't able to figure out, how to even send offline messages to the server query
// func (a Agent) UpdateReadFlagOfflineMessage(sid int, id int, read true)
