package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	cfg "github.com/MaikelVeen/branch/pkg/config"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type GetCommand struct {
	Command *cobra.Command
	logger  *slog.Logger
}

func NewGetCommand() *GetCommand {
	cmd := &GetCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	cmd.Command = &cobra.Command{
		Use:   "get <key>",
		Short: "Get the value of a configuration option",
		Args:  cobra.ExactArgs(1),
		RunE:  cmd.Execute,
	}

	return cmd
}

func (c *GetCommand) Execute(_ *cobra.Command, args []string) error {
	key := args[0]

	opt, err := ValididateKey(key)
	if err != nil {
		return err
	}

	config, err := cfg.Load()
	if err != nil {
		return err
	}

	val := opt.CurrentValue(*config)
	if val == nil {
		c.logger.Info("No value set")
		return nil
	}

	c.logger.Info(fmt.Sprintf("%s=%s", key, *val))
	return nil
}
