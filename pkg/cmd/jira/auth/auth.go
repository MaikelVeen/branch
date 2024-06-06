package auth

import "github.com/spf13/cobra"

// AuthCommand is the parent command for all authentication related commands.
type AuthCommand struct {
	Command *cobra.Command
}

func NewRootCommand() *AuthCommand {
	cmd := &AuthCommand{}
	cmd.Command = &cobra.Command{
		Use: "auth",
	}

	cmd.Command.AddCommand(NewInitCommand().Command)
	cmd.Command.AddCommand(NewShowCommand().Command)
	return cmd
}
