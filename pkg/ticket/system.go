// Package ticket provides generic functionality to interact with ticketing systems.
// The interfaces provided allow the cli tool to be agnostic and not bound to a single
// backend.
package ticket

// System represent an identifier for a ticket system.
type System string

// Jira is a proprietary issue tracking product developed by Atlassian.
const Jira System = "jira"

type TicketSystem interface {
	Authenticate(data interface{}) (User, error)

	// LoginScenario returns a function that will execute a number of prompts
	// to gather the credentials needed to authenticate with the system.
	// The return value interface{} of the returned function represent the login data.
	LoginScenario() LoginScenario

	// Ticker returns the details of a ticket in the ticket system.
	Ticket(key string) (Ticket, error)

	GetBaseFromTicketType(typ string) string

	SaveCredentials() error

	ValidateKey(key string) error
}

type LoginScenario func() (interface{}, error)
