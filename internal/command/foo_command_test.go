package command

import "testing"

func TestFooCommand_Help(t *testing.T) {
	cmd := FooCommand{}
	msg := cmd.Help()

	if msg != "This is Foo" {
		t.Log("Help() returned unexpected value", msg)
		t.Fail()
	}
}

func TestFooCommand_Run(t *testing.T) {
	cmd := FooCommand{}
	ret := cmd.Run([]string{})

	if ret != 0 {
		t.Log("Run() returned a value different than 0")
		t.Fail()
	}

	ret = cmd.Run([]string{"bar"})

	if ret != 0 {
		t.Log("Run() returned a value different than 0")
		t.Fail()
	}
}

func TestFooCommand_Synopsis(t *testing.T) {
	cmd := FooCommand{}
	msg := cmd.Synopsis()

	if msg != "This is Foo" {
		t.Log("Synopsis() returned unexpected value", msg)
		t.Fail()
	}
}
