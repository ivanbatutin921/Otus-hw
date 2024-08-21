package main

import (
	"errors"
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
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	t.conn = conn
	return err
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return errConnection
	}
	defer t.conn.Close() // close the connection when done
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return errConnection
	}
	defer t.conn.Close() // close the connection when done
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}
