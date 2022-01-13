package cmd

import (
	"errors"

	"github.com/MaikelVeen/branch/printer"
	"github.com/MaikelVeen/branch/prompt"
	"github.com/MaikelVeen/branch/ticket"
	"github.com/pterm/pterm"
	"github.com/tucnak/climax"
)

// TODO: make configurable
const keyRingService = "branch-cli"
const keyRingUser = "branch-cli-anon"

func GetLoginCommand() climax.Command {
	return climax.Command{
		Name:   "login",
		Brief:  "authenticates with Jira",
		Handle: ExecuteLoginCommand,
	}
}

// ExecuteLoginCommand executes the login command.
func ExecuteLoginCommand(ctx climax.Context) int {
	// Print a large text with differently colored letters.
	_ = pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("branch", pterm.NewStyle(pterm.FgMagenta)),
	).Render()

	printer.Greet("Welcome to the branch command line interface!")

	// Ask the user for type of system
	systemPrompt := prompt.Prompt{
		InfoLines: []string{"Which system are you logging into ?", "Options are: jira"},
		Label:     "System",
		Validator: func(s string) error {
			if StringInSliceCaseInsensitive(s, ticket.SupportedTicketSystems) {
				return nil
			}

			return errors.New("")
		},
		Retries: 999,
		Invalid: "That is not a valid option!",
	}

	system, err := systemPrompt.Run()
	if err != nil {
		printer.Error(nil, err)
		return 1
	}

	printer.Success(system)

	return 0
}
