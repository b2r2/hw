package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var exp = Environment{
	"BAR":   EnvValue{"bar", false},
	"FOO":   EnvValue{"   foo\nwith new line", false},
	"EMPTY": EnvValue{"", true},
	"HELLO": EnvValue{"\"hello\"", false},
	"UNSET": EnvValue{"", true},
}

func TestReadDir(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")

		require.NoError(t, err)
		require.Equal(t, exp, env)
	})
	t.Run("wrong dir case", func(t *testing.T) {
		env, err := ReadDir("/dev/null")

		require.Error(t, err)
		require.Nil(t, env)
	})
	t.Run("empty dir case", func(t *testing.T) {
		path := t.TempDir()
		err := os.Chmod(path, os.ModePerm)
		require.NoError(t, err)

		env, err := ReadDir(path)

		require.Error(t, err)
		require.Nil(t, env)
	})

	t.Run("filename contains character '='", func(t *testing.T) {
		f, err := ioutil.TempFile("./testdata/env", "out=*")
		require.NoError(t, err)

		path := "./testdata/env"
		env, err := ReadDir(path)

		require.Error(t, err)
		require.Nil(t, env)

		require.NoError(t, os.Remove(f.Name()))
	})
}
