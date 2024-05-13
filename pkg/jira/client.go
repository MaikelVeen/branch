package jira

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	BaseURL = "https://%s.atlassian.net/"
)

var ErrUnauthorized = errors.New("received 401")
var ErrNotFound = errors.New("received 404")

type jiraClient struct {
	BaseURL string `json:"url"`
	Email   string `json:"email"`
	Token   string `json:"token"`

	B Backend `json:"-"`
}

type Client interface {
	// SaveToKeyring saves the credentials of the current client to the
	// system keyring.
	SaveToKeyring(service, user string) error

	MyselfResource
	IssueResource
}

// InitializeAPIFromInit returns a new instance of JiraClient based on the
// passed init command values.
func InitializeAPIFromInit(email, domain, token string) Client {
	url := fmt.Sprintf(BaseURL, domain)
	b := NewBackend(url)

	return &jiraClient{
		Email:   email,
		BaseURL: url,
		Token:   token,
		B:       b,
	}
}

// NewJiraClient returns a new instance of JiraApi with credentials
// gathered from the local keyring.
func NewJiraClient(service, user string) (Client, error) {
	client := &jiraClient{}

	creds, err := keyring.Get(service, user)
	if err != nil {
		return client, err
	}

	err = json.Unmarshal([]byte(creds), &client)
	if err != nil {
		return client, err
	}

	client.B = NewBackend(client.BaseURL)

	return client, nil
}

// SaveToKeyring saves the credentials of the current client to the
// system keyring.
func (c *jiraClient) SaveToKeyring(service, user string) error {
	dataBytes, err := json.Marshal(c)
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
