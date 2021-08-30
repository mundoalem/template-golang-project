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
