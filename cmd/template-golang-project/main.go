package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/mundoalem/template-golang-project/internal/command"
)

const AppName = "template-golang-project"
const Version = "1.0.0"

func main() {
	c := cli.NewCLI(AppName, Version)
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"foo": func() (cli.Command, error) {
			return &command.FooCommand{}, nil
		},
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
