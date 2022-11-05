package cmd

import (
	"github.com/spf13/cobra"
)

type pullRequestCmd struct {
	cmd *cobra.Command
}

func newPullRequestCommand() *pullRequestCmd {
	pr := &pullRequestCmd{}

	pr.cmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "creates a new git branch based on a ticket identifier",
		RunE:    pr.runPullRequestCommand,
	}

	return pr
}

func (c *pullRequestCmd) runPullRequestCommand(cmd *cobra.Command, args []string) error {
	return nil
}
