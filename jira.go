package main

type JiraClient struct {
	Url    string
	Email  string
	Domain string
}

type JiraApi interface {
	GetClientCreds() JiraClient
	SaveToKeyring() error
}

// NewJiraClient returns a new instance of JiraApi based on the
// passed init command struct.
func InitializeApiFromInit(c initCommand) *JiraApi {
	return nil
}

// NewJiraApi returns a new instance of JiraApi with credentials
// gathered from the local keyring.
func NewJiraApi() *JiraApi {
	return nil
}
