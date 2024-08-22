package telnet

import (
	"errors"
	"io"
)

var ErrConnection = errors.New("connection error")

type TelnetClientInterface interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}