package jira

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

var ErrUnauthorized = errors.New("received 401")

type jiraClient struct {
	BaseURL string `json:"url"`
	Email   string `json:"email"`
	Token   string `json:"token"`

	B Backend
}

type JiraClient interface {
	// SaveToKeyring saves the credentials of the current client to the
	// system keyring.
	SaveToKeyring(service, user string) error

	MyselfResource
}

// InitializeApiFromInit returns a new instance of JiraClient based on the
// passed init command values.
func InitializeApiFromInit(email, domain, token string) JiraClient {
	url := buildBaseUrl(domain)
	b := NewBackend(url)

	return &jiraClient{
		Email:   email,
		BaseURL: url,
		Token:   token,
		B:       b,
	}
}

func buildBaseUrl(domain string) string {
	return fmt.Sprintf("https://%s.atlassian.net/", domain)
}

// NewJiraApi returns a new instance of JiraApi with credentials
// gathered from the local keyring.
func NewJiraApi() JiraClient {
	return &jiraClient{}
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
