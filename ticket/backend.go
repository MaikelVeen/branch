// Package ticket provides generic functionality to interact with ticketing systems.
// The interfaces provided allow the cli tool to be agnostic and not bound to a single
// backend.
package ticket

type TicketSystem interface {
	Authenticate(data interface{}) error
	GetTicketName(key string) (string, error)
	LoadCredentials(s, u string) interface{}
	SaveCredentials(s, u string) interface{}
	ValidateKey(key string) error
	GetLoginScenario() LoginScenario
}

type LoginScenario func() (interface{}, error)

type SupportedTicketSystem string

const Jira SupportedTicketSystem = "jira"

func GetTicketSystem(s SupportedTicketSystem) TicketSystem {
	return nil
}

var SupportedTicketSystems []string = []string{string(Jira)}
