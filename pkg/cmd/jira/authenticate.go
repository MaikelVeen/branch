package jira

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type AuthenticateCommand struct {
	Command *cobra.Command
}

func NewAuthenticateCommand() *AuthenticateCommand {
	ac := &AuthenticateCommand{}
	ac.Command = &cobra.Command{
		Use:   "authenticate",
		Short: "Authenticate with Jira",
		RunE:  ac.Execute,
	}

	return ac
}

type AuthDetails struct {
	Email  string `json:"email"`
	Domain string `json:"domain"`
	Token  string `json:"token"`
}

func (ac *AuthenticateCommand) Execute(_ *cobra.Command, _ []string) error {
	details := &AuthDetails{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your email").
				Value(&details.Email),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your Jira domain").
				Value(&details.Domain),
		),
		huh.NewGroup(
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Title("Enter your API token").
				Description("You can generate this from your Jira account settings. See: https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/").
				Value(&details.Token),
		),
	)

	return form.Run()
}
