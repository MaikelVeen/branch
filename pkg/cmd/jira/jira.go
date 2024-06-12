package jira

import (
	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/spf13/cobra"
)

// RootJiraCommand is the parent command for all Jira related commands.
type RootJiraCommand struct {
	Command *cobra.Command
}

func NewCommand() *RootJiraCommand {
	jc := &RootJiraCommand{}
	jc.Command = &cobra.Command{
		Use:   "jira",
		Short: "Interact with Jira",
	}

	jc.Command.AddCommand(auth.NewCommand().Command)
	return jc
}
