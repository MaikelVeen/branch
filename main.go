package main

import (
	"github.com/tucnak/climax"

	"github.com/MaikelVeen/branch/cmd"
)

func main() {
	cli := climax.New("branch")
	cli.Brief = "branch is a small CLI tool to automatically create git branches based on tickets."
	cli.Version = "stable"

	cli.AddCommand(cmd.GetLoginCommand())
	cli.AddCommand(cmd.GetCreateCommand())

	cli.Run()
}
