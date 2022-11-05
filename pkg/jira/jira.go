package jira

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MaikelVeen/branch/pkg/prompt"
	"github.com/MaikelVeen/branch/pkg/ticket"
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

func NewAuthenticatedJira(service, user string) (ticket.TicketSystem, error) {
	j := &jiraImplemtation{
		KeyService: service,
		KeyUser:    user,
		Suffix:     "jira",
	}

	s := fmt.Sprintf("%s.%s", j.KeyService, j.Suffix)
	c, err := NewJiraClient(s, j.KeyUser)
	if err != nil {
		return nil, err
	}

	j.client = c

	return j, nil
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

func (j *jiraImplemtation) GetTicket(key string) (ticket.Ticket, error) {
	issue, err := j.client.GetIssue(key)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return ticket.Ticket{}, ticket.ErrNotUnauthorized
		}

		if errors.Is(err, ErrNotFound) {
			return ticket.Ticket{}, ticket.ErrNotFound
		}
	}

	return ticket.Ticket{
		Title: issue.Fields.Summary,
		Type:  issue.Fields.Issuetype.Name,
		Key:   issue.Key,
	}, nil
}

func (j *jiraImplemtation) GetBaseFromTicketType(typ string) string {
	// TODO: define defaults in git package
	for _, sliceItem := range []string{"bug"} {
		if strings.EqualFold(typ, sliceItem) {
			return "hotfix"
		}
	}

	return "feature"
}

func (j *jiraImplemtation) SaveCredentials() error {
	service := fmt.Sprintf("%s.%s", j.KeyService, j.Suffix)
	return j.client.SaveToKeyring(service, j.KeyUser)
}

func (j *jiraImplemtation) ValidateKey(key string) error {
	// TODO: Implement key checking
	// https://support.atlassian.com/jira-software-cloud/docs/what-is-an-issue/
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
