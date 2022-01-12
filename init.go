package main

import (
	"encoding/json"

	"github.com/mkideal/cli"
	"github.com/zalando/go-keyring"
)

const service = "branch-cli"
const user = "branch-cli-anon"

type initCommand struct {
	Help   bool   `cli:"h,help" usage:"show help" json:"-"`
	Email  string `cli:"e" usage:"The email associated with your Jira Account" json:"email"`
	Domain string `cli:"d" usage:"The domain of your Jira the part in the url before: atlassian.net" json:"domain"`
	Token  string `cli:"t" usage:"The API token of your Jira" json:"token"`
}

func (argv *initCommand) AutoHelp() bool {
	return argv.Help
}

var InitCommand = &cli.Command{
	Name: "init",
	Desc: "Init will initialize the cli for further use",
	Argv: func() interface{} { return new(initCommand) },
	Fn:   executeInit,
}

//TODO validate fields of struct.
func executeInit(ctx *cli.Context) error {
	//argv := ctx.Argv().(*initCommand)

	//err := saveCredsToKeychain(nil)
	//fmt.Printf("%s", err)
	return nil
}

func CheckValid(cmd *initCommand) (*JiraClient, error) {
	return nil, nil
}

func SaveToKeychain(cmd *initCommand) error {
	dataBytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	data := string(dataBytes)

	err = keyring.Set(service, user, data)
	if err != nil {
		return err
	}

	return nil
}
