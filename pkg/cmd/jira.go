package cmd

import (
	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/spf13/cobra"
)

// JiraCommand is the parent command for all Jira related commands.
type JiraCommand struct {
	Command *cobra.Command
}

func NewJiraCommand() *JiraCommand {
	jc := &JiraCommand{}

	jc.Command = &cobra.Command{
		Use:   "jira",
		Short: "Interact with Jira",
	}

	jc.Command.AddCommand(jira.NewAuthenticateCommand().Command)
	return jc
}
