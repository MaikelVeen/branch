package jira

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	client "github.com/MaikelVeen/branch/pkg/jira"
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

type AuthenticationDetails struct {
	Email     string
	Subdomain string
	Token     string
}

func (ac *AuthenticateCommand) Execute(cmd *cobra.Command, _ []string) error {
	details := &AuthenticationDetails{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your email").
				Value(&details.Email),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your Jira subdomain").
				Value(&details.Subdomain),
		),
		huh.NewGroup(
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Title("Enter your API token").
				Description("You can generate this from your Jira account settings").
				Value(&details.Token),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	// Get a Jira client.
	baseURL := fmt.Sprintf(client.BaseURLTemplate, details.Subdomain)
	c, err := client.NewClient(baseURL, client.WithBasicAuthentication(details.Email, details.Token))
	if err != nil {
		return err
	}

	// Get the current user.
	_, err = c.Myself.Myself(cmd.Context())
	if err != nil {
		return err
	}

	return nil
}
