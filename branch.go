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
	//TODO: replace color.red calls with generic error call with debug option.

	_, ok := ctx.Get("issue")
	if !ok {
		color.Red("issue argument is not optional")
		return 0
	}

	if err := ExecGitStatus(); err != nil {
		color.Red("Git status command failed, are you in a git repo?")
		return 0
	}

	if err := ExecCleanTreeCheck(); err != nil {
		color.Red("Working tree is not clean; commit or stash changes and try again")
		return 0
	}

	con, err := ExecBranchCheck()
	if err != nil {
		color.Red("Failed to check what branch you are one")
	}

	if !con {
		return 0
	}

	color.Blue("almost there")

	return 0
}

// ExecGitStatus starts the git status command and waits for it to complete.
func ExecGitStatus() error {
	cmd := exec.Command("git", "status")
	return cmd.Run()
}

// ExecBranchCheck checks what the current branch of the user is.
//
// When the user is not on develop he is prompted for a confirmation.
func ExecBranchCheck() (bool, error) {
	out, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return false, err
	}

	branch := string(out)
	if branch != "develop" {
		color.Yellow("You are currently not on the develop branch")
		c := UserConfirmation("Do you wish to continue?", 2)

		return c, err
	}

	return true, err
}

// ExecCleanTreeCheck returns an error when the user does
// not have a clean working tree.
func ExecCleanTreeCheck() error {
	cmd := exec.Command("git", "diff-index", "--quiet", "HEAD")
	return cmd.Run()
}
