package main

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/MaikelVeen/branch/jira"
	"github.com/mkideal/cli"
)

// TODO: make configurable
const keyRingService = "branch-cli"
const keyRingUser = "branch-cli-anon"

type loginCommand struct {
	Help   bool   `cli:"h,help" usage:"show help" json:"-"`
	Email  string `cli:"e" usage:"The email associated with your Jira Account" json:"email"`
	Domain string `cli:"d" usage:"The domain of your Jira the part in the url before: atlassian.net" json:"domain"`
	Token  string `cli:"t" usage:"The API token of your Jira account" json:"token"`
}

func (argv *loginCommand) AutoHelp() bool {
	return argv.Help
}

var LoginCommand = &cli.Command{
	Name: "login",
	Desc: "Login will validate the passed credentials and save them to the keyring",
	Argv: func() interface{} { return new(loginCommand) },
	Fn:   ExecuteLoginCommand,
}

// Validate implements cli.Validator interface
func (argv *loginCommand) Validate(ctx *cli.Context) error {
	if !validEmail(argv.Email) {
		return fmt.Errorf("%s is not a valid email address", argv.Email)
	}

	if argv.Domain == "" {
		return errors.New("domain cannot be empty")
	}

	if argv.Token == "" {
		return errors.New("token cannot be empty")
	}

	return nil
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ExecuteLoginCommand is the main entry point for the login command.
//
// The login command first initializes an jira api client, verifies
// the credentials and saves them to the keyring
func ExecuteLoginCommand(ctx *cli.Context) error {
	arguments := ctx.Argv().(*loginCommand)

	client := jira.InitializeApiFromInit(arguments.Email, arguments.Domain, arguments.Token)

	return client.SaveToKeyring(keyRingService, keyRingUser)
}
