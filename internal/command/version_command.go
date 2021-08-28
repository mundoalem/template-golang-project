package command

import (
	"fmt"
)

type VersionCommand struct {
	Commit    string
	BuildTime string
	Version   string
}

func (c *VersionCommand) Help() string {
	return "Shows the command version and build metadata"
}

func (c *VersionCommand) Run(args []string) int {
	fmt.Println("Version:    " + c.Version)
	fmt.Println("Build Time: " + c.BuildTime)
	fmt.Println("Commit:     " + c.Commit)

	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Shows the command version and build metadata"
}
