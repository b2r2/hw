package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var timeout string

var (
	ErrNotEnoughArguments = errors.New("not enough arguments, usage: command --timeout[optional] host port")
	ErrTooMuchArguments   = errors.New("too much arguments, usage: command --timeout[optional] host port")
)

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "timeout to connect to the server (by default 10s)")
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalln(ErrNotEnoughArguments)
	}
	if flag.NArg() > 3 {
		log.Fatalln(ErrTooMuchArguments)
	}
	timeout, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalln(err)
	}

	args := flag.Args()
	telnet := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	if err := telnet.Run(); err != nil {
		log.Fatalln(err)
	}
}
