package main

import (
	"github.com/tucnak/climax"
)

func main() {
	cli := climax.New("branch")
	cli.Brief = ""
	cli.Version = "stable"

	cli.AddCommand(GetLoginCommand())
	cli.Run()
}
