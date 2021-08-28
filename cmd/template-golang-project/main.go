package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/mundoalem/template-golang-project/internal/command"
)

const AppName = "template-golang-project"

// The following values are set during build time by the linker
var (
	Commit    string
	BuildTime string
	Version   string
)

type Empty struct{}

func main() {
	c := cli.NewCLI(AppName, Version)
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Commit:    Commit,
				BuildTime: BuildTime,
				Version:   Version,
			}, nil
		},
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
