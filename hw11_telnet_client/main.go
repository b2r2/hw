package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"time"
)

var timeout time.Duration

var (
	ErrNotEnoughArguments = errors.New("not enough arguments, usage: command --timeout[optional] host port")
	ErrTooMuchArguments   = errors.New("too much arguments, usage: command --timeout[optional] host port")
)

func main() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout to connect to the server (by default 10s")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalln(ErrNotEnoughArguments)
	}
	if flag.NArg() > 3 {
		log.Fatalln(ErrTooMuchArguments)
	}

	args := flag.Args()
	telnet := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	if err := telnet.Run(); err != nil {
		log.Fatalln(err)
	}
}
