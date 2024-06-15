package cmd

import "github.com/spf13/cobra"

type CopyCommand struct {
	Command *cobra.Command
}

func NewCopyCommand() *CopyCommand {
	cc := &CopyCommand{}

	cc.Command = &cobra.Command{
		Use:     "copy",
		Aliases: []string{"cp"},
		Short:   "Copies the current branch name",
		Args:    cobra.NoArgs,
		RunE:    cc.Execute,
	}

	return cc
}

func (cc *CopyCommand) Execute(_ *cobra.Command, _ []string) error {
	return nil
}
