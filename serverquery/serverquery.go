package serverquery

import (
	"fmt"
	"net"
	"time"

	"github.com/schoeppi5/libts/core"

	"github.com/schoeppi5/libts"
)

// ServerQuery wraps around ts3.Client and has some additional methods
type ServerQuery struct {
	Host     string
	Port     int
	Username string
	Password string
	conn     net.Conn
	serverID int
	in       chan []byte
	notify   chan []byte
}

// NewServerQuery establishes the tcp connection to teamspeak and logs the query in once cennected
func NewServerQuery(host string, port int, username string, password string) (*ServerQuery, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprint(port)), 10*time.Second)
	if err != nil {
		return nil, err
	}
	in := make(chan []byte)
	notify := make(chan []byte)
	go core.Demultiplexer(conn, in, notify) // split the input from teamspeak into notifys and everything else
	sq := ServerQuery{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		in:       in,
		conn:     conn,
		notify:   notify,
	}
	// slurp the header
	err = core.ReadHeader(in)
	if err != nil {
		conn.Close()
		return nil, err
	}
	go core.KeepAlive(conn, 200*time.Second)
	sq.login(username, password)
	return &sq, nil
}

// Close the tcp connection
func (sq *ServerQuery) Close() {
	sq.logout()
	sq.quit()
	sq.conn.Close()
}

// Login to a useraccount
func (sq *ServerQuery) login(username, password string) error {
	login := libts.Request{
		Command: "login",
		Args: map[string]interface{}{
			"client_login_name":     username,
			"client_login_password": password,
		},
	}
	_, err := sq.DoRaw(login)
	if err != nil {
		return err
	}
	return nil
}

// logout from server
func (sq *ServerQuery) logout() error {
	logout := libts.Request{
		Command: "logout",
	}
	_, err := sq.DoRaw(logout)
	if err != nil {
		return err
	}
	return nil
}

func (sq *ServerQuery) quit() error {
	quit := libts.Request{
		Command: "quit",
	}
	_, err := sq.DoRaw(quit)
	if err != nil {
		return err
	}
	return nil
}
