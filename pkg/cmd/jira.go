package cmd

import (
	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/spf13/cobra"
)

// JiraCommand is the parent command for all Jira related commands.
type JiraCommand struct {
	cmd *cobra.Command
}

// NewJiraCommand returns a new JiraCommand.
func NewJiraCommand() *JiraCommand {
	jc := &JiraCommand{}

	jc.cmd = &cobra.Command{
		Use:   "jira",
		Short: "Interact with Jira",
	}

	jc.cmd.AddCommand(jira.NewAuthenticateCommand().Command)
	return jc
}

func (jc *JiraCommand) Command() *cobra.Command {
	return jc.cmd
}
