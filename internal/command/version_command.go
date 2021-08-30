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

// Package command provides all the commands used by the main function
package command

import (
	"fmt"
)

// Version command implements the interface cli.Command. It outputs to the console metadata about
// the program.
type VersionCommand struct {
	Commit    string
	BuildTime string
	Version   string
}

// Help returns a string describing the Version command in more details to the user
func (c *VersionCommand) Help() string {
	return "Shows the command version and build metadata"
}

// Run is the core function of the Version command which runs the requested process
func (c *VersionCommand) Run(args []string) int {
	fmt.Println("Version:    " + c.Version)
	fmt.Println("Build Time: " + c.BuildTime)
	fmt.Println("Commit:     " + c.Commit)

	return 0
}

// Synopsys returns a string containing a short message about the Version command
func (c *VersionCommand) Synopsis() string {
	return "Shows the command version and other related information"
}
