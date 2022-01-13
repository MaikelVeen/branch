package cmd

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/MaikelVeen/branch/jira"
	"github.com/fatih/color"
	"github.com/tucnak/climax"
)

// TODO: make configurable
const keyRingService = "branch-cli"
const keyRingUser = "branch-cli-anon"

func GetLoginCommand() climax.Command {
	return climax.Command{
		Name:  "login",
		Brief: "authenticates with Jira",

		Flags: []climax.Flag{
			{
				Name:     "email",
				Short:    "e",
				Usage:    `--email="."`,
				Help:     `The email associated with your Jira Account`,
				Variable: true,
			},
			{
				Name:     "domain",
				Short:    "d",
				Usage:    `--domain="."`,
				Help:     `The domain of your Jira the part in the url before: atlassian.net" json:"domain`,
				Variable: true,
			},
			{
				Name:     "token",
				Short:    "t",
				Usage:    `--token="."`,
				Help:     `The API token of your Jira account`,
				Variable: true,
			},
		},
		Handle: HandleLoginCommand,
	}
}

type LoginCommand struct {
	Email  string
	Token  string
	Domain string
}

func LoginCommandFromCliCtx(ctx climax.Context) (*LoginCommand, error) {
	//TODO: make more dry, tag based reflective lookup ?
	cmd := &LoginCommand{}

	if email, ok := ctx.Get("email"); !ok {
		color.Red("email not set in argument list")

		return cmd, errors.New("email not set")
	} else {
		if !validEmail(email) {
			color.Red("email is not valid")

			return cmd, errors.New("email not valid")
		}
		cmd.Email = email
	}

	if token, ok := ctx.Get("token"); !ok {
		color.Red("token not set in argument list")

		return cmd, errors.New("token not set")
	} else {
		cmd.Token = token
	}

	if domain, ok := ctx.Get("domain"); !ok {
		color.Red("email not set in argument list")

		return cmd, errors.New("email not set")
	} else {
		cmd.Domain = domain
	}

	return cmd, nil
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ExecuteLoginCommand is the main entry point for the login command.
//
// The login command first initializes an jira api client, verifies
// the credentials and saves them to the keyring
func HandleLoginCommand(ctx climax.Context) int {
	arguments, err := LoginCommandFromCliCtx(ctx)
	if err != nil {
		return 1
	}

	client := jira.InitializeApiFromInit(arguments.Email, arguments.Domain, arguments.Token)

	user, err := client.GetCurrentUser()
	if err != nil {
		if errors.Is(err, jira.ErrUnauthorized) {
			color.Red("Invalid credentials")
		} else {
			color.Red("Unknown error")
		}

		return 1
	}

	err = client.SaveToKeyring(keyRingService, keyRingUser)
	if err != nil {
		color.Red("Credentials valid, but could not be saved to keyring")
		return 1
	}

	color.Green(fmt.Sprintf("Authenticated successfully as %s (%s)", user.DisplayName, user.EmailAddress))
	return 0
}
