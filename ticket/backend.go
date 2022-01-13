// Package ticket provides generic functionality to interact with ticketing systems.
// The interfaces provided allow the cli tool to be agnostic and not bound to a single
// backend.
package ticket

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
	LoadCredentials() interface{}
	SaveCredentials() interface{}
	ValidateKey(key string) error
}

type LoginScenario func() (interface{}, error)

type User struct {
	DisplayName string
	Email       string
}
