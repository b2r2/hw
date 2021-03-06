package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.connect())
			defer func() { require.NoError(t, client.close()) }()

			in.WriteString("hello\n")
			err = client.send()
			require.NoError(t, err)

			err = client.receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
	t.Run("invalid address", func(t *testing.T) {
		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)
		var in io.ReadCloser
		var out io.Writer
		err = NewTelnetClient("google.com:", timeout, in, out).Run()
		require.Error(t, err)
	})
	t.Run("missing address", func(t *testing.T) {
		timeout, err := time.ParseDuration("1s")
		require.NoError(t, err)
		var in io.ReadCloser
		var out io.Writer
		err = NewTelnetClient("", timeout, in, out).Run()
		require.Error(t, err)
	})
}
