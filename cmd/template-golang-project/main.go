// Package main contains the main code for the application.
package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/mundoalem/template-golang-project/internal/command"
)

// The following values are set during build time through the linker flags.
var (
	// AppName is the name of current application.
	AppName string = "app"
	// Commit is the hash of the commit used to build the current binary.
	Commit string
	// BuildTime is a representation of the build process timestamp in RFC3339 format.
	BuildTime string
	// Version is the current version of the binary.
	Version string = "dev"
)

// Main function of the program, here you should only parse the command line arguments and call
// the cli.Command object which will run the requested process.
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
