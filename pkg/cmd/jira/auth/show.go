package auth

import (
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
	ac := &ShowCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	ac.Command = &cobra.Command{
		Use:   "show",
		Short: "Display the current Jira auth context",
		RunE:  ac.Execute,
		Args:  cobra.NoArgs,
	}

	return ac
}

func (ac *ShowCommand) Execute(_ *cobra.Command, _ []string) error {
	auth, err := LoadUserContext()
	if err != nil {
		ac.logger.Error(err.Error())
	}

	if auth == nil {
		ac.logger.Info("No authentication context found")
		return nil
	}

	ac.logger.Info(fmt.Sprintf("Authenticated as %s(%s)", auth.DisplayName, auth.EmailAddress))
	return nil
}
