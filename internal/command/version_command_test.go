package command

import "testing"

func TestVersionCommand_Help(t *testing.T) {
	cmd := VersionCommand{}
	msg := cmd.Help()

	if msg != "Shows the command version and build metadata" {
		t.Log("Help() returned unexpected value", msg)
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
	msg := cmd.Synopsis()

	if msg != "Shows the command version" {
		t.Log("Synopsis() returned unexpected value", msg)
		t.Fail()
	}
}
