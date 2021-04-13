package command

import (
	"fmt"
)

type FooCommand struct {
}

func (c *FooCommand) Help() string {
	return "This is Foo"
}

func (c *FooCommand) Run(args []string) int {
	if len(args) <= 0 {
		return 0
	}

	for _, arg := range args {
		fmt.Println(arg)
	}

	return 0
}

func (c *FooCommand) Synopsis() string {
	return "This is Foo"
}
