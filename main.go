package main

import (
	"fmt"
	"os"

	"github.com/mkideal/cli"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(child),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

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

type initCommand struct {
	Help   bool   `cli:"h,help" usage:"show help" json:"-"`
	Email  string `cli:"e" usage:"The email associated with your Jira Account" json:"email"`
	Domain string `cli:"d" usage:"The domain of your Jira the part in the url before: atlassian.net" json:"domain"`
	Token  string `cli:"t" usage:"The API token of your Jira" json:"token"`
}

func (argv *initCommand) AutoHelp() bool {
	return argv.Help
}

var child = &cli.Command{
	Name: "init",
	Desc: "Init will initialize the cli for further use",
	Argv: func() interface{} { return new(initCommand) },
	Fn:   executeInit,
}

func executeInit(ctx *cli.Context) error {
	ctx.String(ctx.Color().Blue("yes"))
	return nil
}

func saveCredsToKeychain(cmd initCommand) error {
	return nil
}
