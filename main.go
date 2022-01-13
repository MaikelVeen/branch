package main

import (
	"github.com/tucnak/climax"

	"github.com/MaikelVeen/branch/cmd"
)

func main() {
	cli := climax.New("branch")
	cli.Brief = ""
	cli.Version = "stable"

	cli.AddCommand(cmd.GetLoginCommand())
	cli.AddCommand(cmd.GetBranchCommand())
	cli.Run()
}
