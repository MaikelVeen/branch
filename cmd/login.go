package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MaikelVeen/branch/printer"
	"github.com/MaikelVeen/branch/prompt"
	"github.com/MaikelVeen/branch/ticket"
	"github.com/pterm/pterm"
	"github.com/tucnak/climax"
)

func GetLoginCommand() climax.Command {
	return climax.Command{
		Name:   "login",
		Brief:  "authenticates with ticket system",
		Handle: ExecuteLoginCommand,
	}
}

// ExecuteLoginCommand executes the login command.
func ExecuteLoginCommand(ctx climax.Context) int {
	// Print a large text with differently colored letters.
	_ = pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("branch", pterm.NewStyle(pterm.FgGreen)),
	).Render()

	printer.Greet("Welcome to the branch command line interface!")

	systemType, err := getSystemType()
	if err != nil {
		printer.Error(nil, err)
		return 1
	}

	system := GetNewTicketSystem(systemType)
	loginScenario := system.GetLoginScenario()

	i := 3

	for i > 0 {
		// Execute the login scenario of the system.
		loginData, err := loginScenario()
		if err != nil {
			printer.Error(nil, err)
			return 1
		}

		// Authenticate with the login data.
		user, err := system.Authenticate(loginData)
		if err != nil {
			if errors.Is(err, ticket.ErrNotUnauthorized) {
				printer.Warning("Those credentials are invalid, please try again.")
			} else {
				printer.Error(nil, err)
				return 1
			}

			i--
			continue
		}

		err = user.SaveToDisk()
		if err != nil {
			printer.Error(nil, err)
			return 1
		}

		printer.Success(fmt.Sprintf("Authenticated successfully as %s (%s)", user.DisplayName, user.Email))
		return 0
	}

	printer.Warning("Aborting...")
	return 1
}

// getSystemType returns which SupportedTicketSystem the user is trying to login to.
func getSystemType() (ticket.SupportedTicketSystem, error) {
	// Ask the user for type of system.
	systemPrompt := prompt.Prompt{
		InfoLines: []string{"Which system are you logging into ?", "Options are: jira"},
		Label:     "System",
		Validator: func(s string) error {
			for _, sliceItem := range ticket.SupportedTicketSystems {
				if strings.EqualFold(s, sliceItem) {
					return nil
				}
			}

			return errors.New("")
		},
		Invalid: "That is not a valid option!",
	}

	// Run the first prompt.
	system, err := systemPrompt.Run()
	if err != nil {
		return "", err
	}

	return ticket.SupportedTicketSystem(system), err
}
