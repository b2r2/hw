package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	ds, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment, len(ds))
	for _, i := range ds {
		n := i.Name()
		if i.IsDir() || strings.Contains(n, "=") {
			continue
		}
		v, err := getEnv(filepath.Join(dir, n))
		if err != nil {
			return nil, err
		}
		env[n] = v
	}

	return env, nil
}

func getEnv(path string) (EnvValue, error) {
	f, err := os.Open(path)
	if err != nil {
		return EnvValue{}, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	s := bufio.NewScanner(f)
	s.Scan()
	if err := s.Err(); err != nil {
		return EnvValue{}, err
	}
	return handleEnv(s.Bytes()), nil
}

func handleEnv(b []byte) EnvValue {
	env := EnvValue{}
	b = bytes.ReplaceAll(
		bytes.TrimRightFunc(b, unicode.IsSpace),
		[]byte{'\x00'},
		[]byte("\n"),
	)
	if len(b) == 0 {
		env.NeedRemove = true
	}
	env.Value = string(b)
	return env
}
