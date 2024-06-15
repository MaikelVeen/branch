package jira

import (
	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/spf13/cobra"
)

// Command is the parent command for all Jira related commands.
type Command struct {
	Command *cobra.Command
}

func NewCommand() *Command {
	jc := &Command{}
	jc.Command = &cobra.Command{
		Use:   "jira",
		Short: "Interact with Jira",
	}

	jc.Command.AddCommand(auth.NewCommand().Command)
	return jc
}
