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
	"strings"
	"testing"
)

func TestVersionCommand_Help(t *testing.T) {
	cmd := VersionCommand{}
	msg := strings.TrimSpace(cmd.Help())

	if msg == "" {
		t.Log("Help() should not return an empty string", msg)
		t.Fail()
	}
}

func TestVersionCommand_Run(t *testing.T) {
	cmd := VersionCommand{}
	ret := cmd.Run([]string{})

	if ret != 0 {
		t.Log("Run() returned a value different than 0")
		t.Fail()
	}
}

func TestVersionCommand_Synopsis(t *testing.T) {
	cmd := VersionCommand{}
	msg := strings.TrimSpace(cmd.Synopsis())

	if msg == "" {
		t.Log("Synopsis() should not return an empty string", msg)
		t.Fail()
	}
}
