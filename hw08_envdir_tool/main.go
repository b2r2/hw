package main

import (
	"errors"
	"log"
	"os"
)

var (
	ErrNotEnoughArgs   = errors.New("not enough arguments")
	ErrNotEnvVariables = errors.New("no variables on the directory")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln(ErrNotEnoughArgs)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(ErrNotEnvVariables)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
