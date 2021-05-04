package main

import (
	"log"
	"os"

	"github.com/egonbraun/golang-project-template/internal/command"
	"github.com/mitchellh/cli"
)

const AppName = "example"
const Version = "1.1.0"

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
