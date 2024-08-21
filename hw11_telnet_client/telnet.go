package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var errConnection = fmt.Errorf("connection error")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	addres  string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		addres:  address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.addres, t.timeout)
	t.conn = conn
	return err
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return errConnection
	}
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return errConnection
	}
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *telnetClient) Close() error {
	if t.conn != nil{
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}
