package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

const defTimeout = 10

func init() {
	flag.DurationVar(
		&timeout, "timeout", time.Second*defTimeout, "connection timeout [default=10s]",
	)
}

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fatal(err)
	}

	defer func(client TelnetClient) {
		if err := client.Close(); err != nil {
			fatal(err)
		}
	}(client)

	log.Printf("Simple telnet: connected to %s \n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)

	defer func() { cancel() }()

	go func() {
		if err := client.Send(); err != nil {
			fatal(err)
		}
		log.Print("...EOF")
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			fatal(err)
		}
		log.Print("...Connection was closed by peer")
	}()

	<-ctx.Done()
}
