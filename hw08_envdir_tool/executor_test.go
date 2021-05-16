package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh", "arg1", "arg2"}
		env := make(Environment)
		env["FOO"] = EnvValue{Value: "foo"}
		env["BAR"] = EnvValue{Value: "bar"}
		env["ADDED"] = EnvValue{Value: "added"}

		code := RunCmd(cmd, env)

		require.Equal(t, 0, code)
	})

	t.Run("nil", func(t *testing.T) {
		require.Equal(t, 1, RunCmd(nil, nil))
	})
}
