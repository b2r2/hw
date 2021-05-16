package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r, err := ReadDir("")
		require.Error(t, err)
		require.Nil(t, r)
	})
	t.Run("positive", func(t *testing.T) {
		exp := Environment{
			"BAR":   EnvValue{"bar", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"EMPTY": EnvValue{"", true},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, exp, env)
	})
}
