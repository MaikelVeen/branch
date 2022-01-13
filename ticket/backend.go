// Package ticket provides generic functionality to interact with ticketing systems.
// The interfaces provided allow the cli tool to be agnostic and not bound to a single
// backend.
package ticket

import (
	"encoding/json"
	"fmt"
	"os"
)

// SupportedTicketSystem represent an identifier for a ticket system.
type SupportedTicketSystem string

const Jira SupportedTicketSystem = "jira"

// SupportedTicketSystems represent the systems that are currently supported by the tool.
var SupportedTicketSystems []string = []string{string(Jira)}

type TicketSystem interface {
	Authenticate(data interface{}) (User, error)
	// GetLoginScenario returns a function that will execute a number of prompts
	// to gather the credentials needed to authenticate with the system.
	// The return value interface{} of the returned function represent the login data.
	GetLoginScenario() LoginScenario
	GetTicketName(key string) (string, error)
	LoadCredentials() (interface{}, error)
	SaveCredentials() error
	ValidateKey(key string) error
}

type LoginScenario func() (interface{}, error)

type User struct {
	DisplayName string
	Email       string
	System      SupportedTicketSystem
}

// SaveToDisk dumps the current user to the file system
// so that it can be retrieved later.
func (u *User) SaveToDisk() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/.branch-cli", dirname))
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(dataBytes)
	if err != nil {
		return err
	}

	return f.Sync()
}
