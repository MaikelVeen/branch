package config

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	ErrInvalidKey = errors.New("invalid key")
)

// Command is the parent command for all configuration related commands.
type Command struct {
	Command *cobra.Command
}

func NewCommand() *Command {
	cmd := &Command{}
	cmd.Command = &cobra.Command{
		Use:   "config",
		Short: "Commands to configure branch",
	}

	cmd.Command.AddCommand(NewSetCommand().Command)
	cmd.Command.AddCommand(NewGetCommand().Command)
	return cmd
}
