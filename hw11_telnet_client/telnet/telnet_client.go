package telnet

import (
	"errors"
	"io"
)

var ErrConnection = errors.New("connection error")

type ClientInterface interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}
