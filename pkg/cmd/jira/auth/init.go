package auth

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type InitCommand struct {
	Command *cobra.Command

	logger *slog.Logger
}

func NewInitCommand() *InitCommand {
	cmd := &InitCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	cmd.Command = &cobra.Command{
		Use:   "init",
		Short: "Initialize the Jira authentication",
		RunE:  cmd.Execute,
		Args:  cobra.NoArgs,
	}

	return cmd
}

func (ac *InitCommand) Execute(cmd *cobra.Command, _ []string) error {
	auth := &Context{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your email").
				Value(&auth.EmailAddress),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your Jira subdomain").
				Value(&auth.Subdomain),
		),
		huh.NewGroup(
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Title("Enter your API token").
				Description("You can generate this from your Jira account settings").
				Value(&auth.Token),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	c, err := newClient(auth)
	if err != nil {
		return err
	}

	user, err := c.Myself.Myself(cmd.Context())
	if err != nil {
		return err
	}
	auth.DisplayName = user.DisplayName

	if err = auth.Save(); err != nil {
		return err
	}

	ac.logger.Info(fmt.Sprintf("Successfully authenticated as %s, saved credentials to keyring", user.DisplayName))
	return nil
}
