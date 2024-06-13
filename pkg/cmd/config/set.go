package config

import (
	cfg "github.com/MaikelVeen/branch/pkg/config"
	"github.com/spf13/cobra"
)

// SetCommand is the command to update the configuration.
type SetCommand struct {
	Command *cobra.Command
}

func NewSetCommand() *SetCommand {
	cmd := &SetCommand{}
	cmd.Command = &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update the configuration",
		Args:  cobra.ExactArgs(2),
		RunE:  cmd.Execute,
	}

	return cmd
}

func (c *SetCommand) Execute(_ *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	opt, err := ValididateKey(key)
	if err != nil {
		return err
	}

	config, err := cfg.Load()
	if err != nil {
		return err
	}

	if err = opt.SetValue(config, value); err != nil {
		return err
	}

	if err = config.Save(); err != nil {
		return err
	}

	return nil
}

func ValididateKey(key string) (*cfg.Option, error) {
	if _, ok := cfg.Options[key]; !ok {
		return nil, ErrInvalidKey
	}

	return cfg.Options[key], nil
}
