package auth

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
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
		Short: "Show the Jira authentication",
		RunE:  ac.Execute,
	}

	return ac
}

func (ac *ShowCommand) Execute(cmd *cobra.Command, _ []string) error {
	auth, err := LoadUserContext()
	// Check for ErrNotFound.
	if err != nil {
		if err == keyring.ErrNotFound {
			ac.logger.Info("No Jira authentication found")
			return nil
		}

		ac.logger.Error(err.Error())
	}

	ac.logger.Info(fmt.Sprintf("Authenticated as %s(%s)", auth.User.DisplayName, auth.EmailAddress))
	return nil
}
