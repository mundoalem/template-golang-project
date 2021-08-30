// This file is part of template-golang-project.
//
// template-golang-project is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// template-golang-project is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with template-golang-project. If not, see <https://www.gnu.org/licenses/>.

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
