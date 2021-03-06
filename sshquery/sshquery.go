package sshquery

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/schoeppi5/libts"
	"github.com/schoeppi5/libts/communication"

	"golang.org/x/crypto/ssh"
)

// SSHQuery uses ssh to connect to teamspeak
type SSHQuery struct {
	Host     string
	Port     int
	Username string
	Password string
	serverID int
	conn     *ssh.Client
	out      io.Writer
	in       chan []byte
	notify   chan []byte
}

// NewSSHQuery opens a new ssh connection to the host and starts a session for communicating with teamspeak
func NewSSHQuery(host string, port int, username, password string) (*SSHQuery, error) {
	config := &ssh.ClientConfig{ // configure ssh connection
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: maybe change this to knownhosts
	}
	conn, err := ssh.Dial("tcp", net.JoinHostPort(host, fmt.Sprint(port)), config)
	if err != nil {
		return nil, err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return nil, err
	}
	input, err := session.StdoutPipe() // the return from teamspeak
	if err != nil {
		conn.Close()
		return nil, err
	}
	out, err := session.StdinPipe() // for sending commands
	if err != nil {
		conn.Close()
		return nil, err
	}
	if err = session.Shell(); err != nil { // starts a new session and attaches stdin and stdout
		conn.Close()
		return nil, err
	}
	in := make(chan []byte, 5)
	sshq := &SSHQuery{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		conn:     conn,
		in:       in,
		out:      out,
		notify:   nil,
	}
	go communication.Read(input, sshq.in, sshq.notify) // split the input from teamspeak into notifys and everything else
	err = communication.ReadHeader(in)                 // slurp header
	if err != nil {
		conn.Close()
		return nil, err
	}
	go communication.KeepAlive(out, 200*time.Second) // keepAlive
	return sshq, nil
}

// Close the ssh connection to teamspeak
func (sq *SSHQuery) Close() {
	sq.logout()
	sq.quit()
}

// logout from server
func (sq *SSHQuery) logout() error {
	logout := libts.Request{
		Command: "logout",
	}
	_, err := sq.DoRaw(logout)
	if err != nil {
		return err
	}
	return nil
}

func (sq *SSHQuery) quit() error {
	quit := libts.Request{
		Command: "quit",
	}
	_, err := sq.DoRaw(quit)
	if err != nil {
		return err
	}
	return nil
}
