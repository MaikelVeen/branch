package cmd

import (
	"errors"
	"fmt"

	"github.com/MaikelVeen/branch/pkg/printer"
	"github.com/MaikelVeen/branch/pkg/ticket"
	"github.com/spf13/cobra"
)

type loginCmd struct {
	cmd *cobra.Command

	systems []ticket.SystemType
}

func newLoginCommand() *loginCmd {
	lc := &loginCmd{}

	lc.cmd = &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(1),
		Short: "Authenticates with a ticket system.",
		RunE:  lc.runLoginCommand,
	}

	return lc
}

func (c *loginCmd) RegisterSystem(sys ticket.SystemType) {
	c.systems = append(c.systems, sys)
}

func (c *loginCmd) runLoginCommand(_ *cobra.Command, args []string) error {
	sys := args[0]

	if err := c.isValidSystem(sys); err != nil {
		return err
	}

	system := getNewTicketSystem(ticket.SystemType(sys))
	loginScenario := system.LoginScenario()

	i := 3

	for i > 0 {
		// Execute the login scenario of the system.
		loginData, err := loginScenario()
		if err != nil {
			printer.Error(nil, err)
			return err
		}

		// Authenticate with the login data.
		user, err := system.Authenticate(loginData)
		if err != nil {
			if errors.Is(err, ticket.ErrNotUnauthorized) {
				printer.Warning("Those credentials are invalid, please try again.")
			} else {
				printer.Error(nil, err)
				return err
			}

			i--
			continue
		}

		if err = user.SaveToDisk(); err != nil {
			printer.Error(nil, err)
			return err
		}

		printer.Success(fmt.Sprintf("Authenticated successfully as %s (%s)", user.DisplayName, user.Email))
		return nil
	}

	printer.Warning("Aborting...")
	return fmt.Errorf("could not authenticate with %s", sys)
}

func (c *loginCmd) isValidSystem(sys string) error {
	for _, rs := range c.systems {
		if string(rs) == sys {
			return nil
		}
	}

	return fmt.Errorf("%s is not a registered ticket system", sys)
}
