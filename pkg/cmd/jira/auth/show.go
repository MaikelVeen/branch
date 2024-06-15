package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type ShowCommand struct {
	Command *cobra.Command
	logger  *slog.Logger
}

func NewShowCommand() *ShowCommand {
	cmd := &ShowCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	cmd.Command = &cobra.Command{
		Use:   "show",
		Short: "Display the current Jira auth context",
		RunE:  cmd.Execute,
		Args:  cobra.NoArgs,
	}

	return cmd
}

func (cmd *ShowCommand) Execute(_ *cobra.Command, _ []string) error {
	auth, err := LoadUserContext()
	if err != nil {
		if errors.Is(err, ErrAuthContextMissing) {
			cmd.logger.Info("No authentication context found")
			return nil
		}
		cmd.logger.Error(err.Error())
	}

	cmd.logger.Info(fmt.Sprintf("Authenticated as %s(%s)", auth.DisplayName, auth.EmailAddress))
	return nil
}
