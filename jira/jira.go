package jira

import (
	"errors"
	"fmt"

	"github.com/MaikelVeen/branch/prompt"
	"github.com/MaikelVeen/branch/ticket"
	"github.com/fatih/color"
)

type jiraImplemtation struct {
	KeyService string
	KeyUser    string
	Suffix     string
	client     JiraClient
}

func NewJira(service, user string) ticket.TicketSystem {
	return &jiraImplemtation{
		KeyService: service,
		KeyUser:    user,
		Suffix:     "jira",
	}
}

type JiraCredentials struct {
	Email  string
	Token  string
	Domain string
}

func (j *jiraImplemtation) Authenticate(data interface{}) (ticket.User, error) {
	tu := ticket.User{
		System: ticket.Jira,
	}

	// Type assert passed interface.
	creds, ok := data.(JiraCredentials)
	if !ok {
		return tu, ticket.ErrTypeAssertionAuth
	}

	// Set the client
	j.client = InitializeApiFromInit(creds.Email, creds.Domain, creds.Token)

	user, err := j.client.GetCurrentUser()
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return tu, ticket.ErrNotUnauthorized
		}
		return tu, err
	}

	// Save the credentials on the user system.
	err = j.SaveCredentials()
	if err != nil {
		return tu, ticket.ErrCredentialSaving
	}

	// Populate struct
	tu.DisplayName = user.DisplayName
	tu.Email = user.EmailAddress

	return tu, nil
}

func (j *jiraImplemtation) GetTicketName(key string) (string, error) {
	return "", nil
}

func (j *jiraImplemtation) LoadCredentials() (interface{}, error) {
	return nil, nil
}

func (j *jiraImplemtation) SaveCredentials() error {
	service := fmt.Sprintf("%s.%s", j.KeyService, j.Suffix)
	return j.client.SaveToKeyring(service, j.KeyUser)
}

func (j *jiraImplemtation) ValidateKey(key string) error {
	return nil
}

func (j *jiraImplemtation) GetLoginScenario() ticket.LoginScenario {
	return func() (interface{}, error) {
		// Get the email address from the user.
		emailPrompt := prompt.Prompt{
			Label:       "Email",
			Validator:   prompt.ValidateEmail,
			Invalid:     "That is not a valid email!",
			LabelColour: color.FgHiGreen,
		}

		email, err := emailPrompt.Run()
		if err != nil {
			return "", err
		}

		// Get the domain from the user.
		domainPrompt := prompt.Prompt{
			Label:       "Domain",
			LabelColour: color.FgHiGreen,
		}

		domain, err := domainPrompt.Run()
		if err != nil {
			return "", err
		}

		// Get the token from the user.
		tokenPrompt := prompt.Prompt{
			Label:       "Token",
			HideEntered: true,
			LabelColour: color.FgHiGreen,
		}

		token, err := tokenPrompt.Run()
		if err != nil {
			return "", err
		}

		return JiraCredentials{
			Email:  email,
			Domain: domain,
			Token:  token,
		}, nil
	}

}
