package jira

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

type jiraClient struct {
	BaseUrl string `json:"url"`
	Email   string `json:"email"`
	Token   string `json:"token"`

	apiPath string
}

type JiraClient interface {
	// SaveToKeyring saves the credentials of the current client to the
	// system keyring.
	SaveToKeyring(service, user string) error
}

// InitializeApiFromInit returns a new instance of JiraClient based on the
// passed init command values.
func InitializeApiFromInit(email, domain, token string) JiraClient {
	return &jiraClient{
		Email:   email,
		BaseUrl: buildBaseUrl(domain),
		Token:   token,
		apiPath: "rest/api/3/",
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
