package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	ErrNotEnoughArgs   = errors.New("not enough arguments: usage command path args")
	ErrNotEnvVariables = errors.New("no variables on the directory")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln(ErrNotEnoughArgs)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(fmt.Errorf("%v: %w", ErrNotEnvVariables, err))
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
