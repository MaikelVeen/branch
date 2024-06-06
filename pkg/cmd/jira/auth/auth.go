package auth

import "github.com/spf13/cobra"

// RootCommand is the parent command for all authentication related commands.
type RootCommand struct {
	Command *cobra.Command
}

func NewRootCommand() *RootCommand {
	cmd := &RootCommand{}
	cmd.Command = &cobra.Command{
		Use: "auth",
	}

	cmd.Command.AddCommand(NewInitCommand().Command)
	cmd.Command.AddCommand(NewShowCommand().Command)
	return cmd
}
