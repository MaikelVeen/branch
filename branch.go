package main

import "github.com/mkideal/cli"

// root command
type branchCommand struct {
	Help  bool   `cli:"h,help" usage:"show help"`
	Issue string `cli:"issue" usage:"the key of the issue"`
}

func (argv *branchCommand) AutoHelp() bool {
	return argv.Help
}

var root = &cli.Command{
	Desc: "",
	Argv: func() interface{} { return new(branchCommand) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*branchCommand)
		ctx.String("Hello, root command, I am %s\n", argv.Issue)
		return nil
	},
}
