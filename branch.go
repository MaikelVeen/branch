package main

import (
	"os/exec"

	"github.com/fatih/color"
	"github.com/tucnak/climax"
)

func GetBranchCommand() climax.Command {
	return climax.Command{
		Name:  "n",
		Brief: "creates a new branch based on a jira issue",

		Flags: []climax.Flag{
			{
				Name:     "issue",
				Short:    "i",
				Usage:    `--issue="."`,
				Help:     `The id of the jira issue`,
				Variable: true,
			},
		},
		Handle: HandleBranchCommand,
	}
}

// HandleBranchCommand is the main entry point for the new command.
func HandleBranchCommand(ctx climax.Context) int {
	issue, ok := ctx.Get("issue")
	if !ok {
		color.Red("issue argument is not optional")
		return 0
	}

	color.Blue("Echo %s", issue)

	if err := ExecGitStatus(); err != nil {
		color.Red("Failed to check git status, are you in a git repo?")
		return 0
	}

	return 0
}

func ExecGitStatus() error {
	cmd := exec.Command("git", "status")
	return cmd.Run()
}
