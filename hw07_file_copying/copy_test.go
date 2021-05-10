package main

import (
	"os"
	"path"
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
				"empty path",
				ErrEmptyPath,
				checkValid("", out, 0, 0),
			},
			{
				"empty to path",
				ErrEmptyToPath,
				checkValid(testdata, "", 0, 0),
			},
			{
				"equal paths",
				ErrEqualPaths,
				checkValid(testdata, testdata, 0, 0),
			},
			{
				"negative limit",
				ErrNegativeLimit,
				checkValid(testdata, out, 0, -1),
			},
			{
				"negative offset",
				ErrNegativeOffset,
				checkValid(testdata, out, -1, 0),
			},
			{
				"offset exceeds file size",
				ErrOffsetExceedsFileSize,
				checkValid(testdata, out, 1<<20, 0),
			},
			{
				"err unsupported file",
				ErrUnsupportedFile,
				checkValid("/dev/urandom", out, 0, 0),
			},
			{
				"valid data",
				nil,
				checkValid(testdata, out, 0, 0),
			},
		}
		for _, ts := range tests {
			require.Equal(t, ts.exp, ts.act, ts.name)
		}
	})
	t.Run("simple copy", func(t *testing.T) {
		testdata := "./testdata/input.txt"
		out := "./tmp/out.txt"
		require.Equal(t, nil, Copy(testdata, out, 0, 0))
		require.Nil(t, os.RemoveAll(path.Dir(out)))
	})

	t.Run("invalid data ", func(t *testing.T) {
		testdata := "./testdata/input.txt"
		require.Error(t, ErrEqualPaths, Copy(testdata, testdata, 0, 0))
	})
}
