package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("validation", func(t *testing.T) {
		const testdata = "./testdata/input.txt"
		out := "out/input.txt"
		tests := []struct {
			name string
			exp  error
			act  error
		}{
			{
				"from path is empty",
				ErrEmptyFromPath,
				Copy("", out, 0, 0),
			},
			{
				"to path is empty",
				ErrEmptyToPath,
				Copy(testdata, "", 0, 0),
			},
			{
				"equal paths",
				ErrEqualPaths,
				Copy(testdata, testdata, 0, 0),
			},
			{
				"negative limit",
				ErrNegativeLimit,
				Copy(testdata, out, 0, -1),
			},
			{
				"negative offset",
				ErrNegativeOffset,
				Copy(testdata, out, -1, 0),
			},
			{
				"offset exceeds file size",
				ErrOffsetExceedsFileSize,
				Copy(testdata, out, 1<<20, 0),
			},
			{
				"err unsupported file",
				ErrUnsupportedFile,
				Copy("/dev/urandom", out, 0, 0),
			},
			{
				"valid data",
				nil,
				Copy(testdata, out, 0, 0),
			},
		}
		for _, ts := range tests {
			require.Equal(t, ts.exp, ts.act, ts.name)
			if ts.exp == nil {
				require.Nil(t, os.Remove(out))
			}
		}
	})
}
