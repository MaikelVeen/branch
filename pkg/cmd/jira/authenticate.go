package jira

import "github.com/spf13/cobra"

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

func (ac *AuthenticateCommand) Execute(_ *cobra.Command, _ []string) error {
	return nil
}
