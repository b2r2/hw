package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	Connect() error
	io.Closer
	Send() error
	Receive() error
	Run() error
}

type Telnet struct {
	conn    net.Conn
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *Telnet) Connect() error {
	dialer := net.Dialer{}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	t.ctx = ctx
	t.cancel = cancel
	conn, err := dialer.DialContext(t.ctx, tcp, t.address)
	if err != nil {
		cancel()
		return err
	}
	t.conn = conn
	return nil
}

func (t *Telnet) Close() error {
	return t.conn.Close()
}

func (t *Telnet) Send() error {
	scanner := bufio.NewScanner(t.in)
	for scanner.Scan() {
		msg := append(scanner.Bytes(), newLine)
		if _, err := t.conn.Write(msg); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (t *Telnet) Receive() error {
	scanner := bufio.NewScanner(t.conn)
	for scanner.Scan() {
		msg := append(scanner.Bytes(), newLine)
		if _, err := t.out.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

func (t *Telnet) Run() error {
	if err := t.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := t.Close(); err != nil {
			return
		}
	}()

	go func() {
		if err := t.Receive(); err != nil {
			fmt.Println(fmt.Errorf("cannot receive message %w", err))
		}
	}()
	go func() {
		if err := t.Send(); err != nil {
			fmt.Println(fmt.Errorf("cannot send message %w", err))
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		fmt.Println(fmt.Errorf("...Connection was closed by peer"))
		t.cancel()
	case <-t.ctx.Done():
		close(sigCh)
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{address: address, timeout: timeout, in: in, out: out}
}
