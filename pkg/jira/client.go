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

// Client manages communication with the Jira API.
type Client struct {
	credentials Credentials
	BaseURL    url.Url
	httpClient *http.Client
}

// NewClient returns a new client with the given options.
func NewClient(creds Credentials, opts ...func(*Client) error) (*Client, error) {
	client := &Client{
		APIKey:     apiKey,
		BaseURL:    "api.example.com",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// TODO: NewRequest. BasicAuth. RoundTripper

func WithHTTPClient(client *http.Client) func(*Client) error {
  return func(c *Client) error {
    c.httpClient = client
    return nil
  }
}