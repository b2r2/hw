package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const (
	tcp     = "tcp"
	newLine = byte('\n')
)

type TelnetClient interface {
	Run() error
	close() error
	connect() error
	send() error
	receive() error
}

type telnet struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *telnet) connect() error {
	conn, err := net.DialTimeout(tcp, t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

func (t *telnet) close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

func (t *telnet) send() error {
	if t.conn == nil {
		return nil
	}
	scanner := bufio.NewScanner(t.in)
	for {
		if scanner.Scan() {
			msg := append(scanner.Bytes(), newLine)
			if _, err := t.conn.Write(msg); err != nil {
				return err
			}
		} else {
			return scanner.Err()
		}
	}
}

func (t *telnet) receive() error {
	if t.conn == nil {
		return nil
	}
	_, err := io.Copy(t.out, t.conn)
	return err
}

func (t *telnet) Run() error {
	if err := t.connect(); err != nil {
		return err
	}
	defer func() {
		if err := t.close(); err != nil {
			log.Fatal(err)
		}
	}()

	errs := make(chan error, 1)
	go func() { errs <- t.receive() }()
	go func() { errs <- t.send() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		if _, err := fmt.Fprint(os.Stderr, "...Connection was closed by peer"); err != nil {
			log.Println(err)
		}
		signal.Stop(sigCh)
		close(sigCh)
	case err := <-errs:
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnet{address: address, timeout: timeout, in: in, out: out}
}
