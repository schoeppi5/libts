package query

import "github.com/schoeppi5/libts"

// Login to a useraccount
// This is not supported by webquery
// username - required
// password - required
func (a Agent) Login(username, password string) error {
	login := libts.Request{
		Command: "login",
		Args: map[string]interface{}{
			"client_login_name":     username,
			"client_login_password": password,
		},
	}
	return a.Query.Do(login, nil)
}

// Logout from a server
// This is not supported by webquery
func (a Agent) Logout() error {
	logout := libts.Request{
		Command: "logout",
	}
	return a.Query.Do(logout, nil)
}
