package telnet

import (
	"io"
	"net"
	"time"
)

type TelnetSimpleClient struct {
	Address string
	Timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewSimpleTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetSimpleClient {
	return &TelnetSimpleClient{
		in:      in,
		out:     out,
		Address: address,
		Timeout: timeout,
	}
}

func (c *TelnetSimpleClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.Address, c.Timeout)
	c.conn = conn
	return err
}

func (c *TelnetSimpleClient) Send() error {
	if c.conn == nil {
		return ErrConnection
	}
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *TelnetSimpleClient) Receive() error {
	if c.conn == nil {
		return ErrConnection
	}
	_, err := io.Copy(c.out, c.conn)
	return err
}

func (c *TelnetSimpleClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}