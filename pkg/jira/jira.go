package jira

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURLTemplate is the template for the Jira API base URL.
	BaseURLTemplate = "https://%s.atlassian.net"

	// DefaultTimeout is the default timeout for the HTTP client.
	DefaultTimeout = 30 * time.Second
)

// Client manages communication with the Jira API.
type Client struct {
	client   *http.Client
	username string
	token    string

	// BaseURL is the base URL for the Jira API.
	BaseURL url.URL

	// UserAgent is the User-Agent string used when making API requests
	// to the Jira API.
	UserAgent string

	// Services used for talking to different parts of the Jira API.

	Issue  *IssueResourceService
	Myself *MyselfResourceService
}

// NewClient returns a new client with the given options.
func NewClient(baseURL string, opts ...func(*Client) error) (*Client, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL: *url,
		client:  &http.Client{Timeout: DefaultTimeout},
	}

	for _, opt := range opts {
		err = opt(client)
		if err != nil {
			return nil, err
		}
	}

	client.Issue = &IssueResourceService{client: client}
	client.Myself = &MyselfResourceService{client: client}

	return client, nil
}

// WithHTTPClient returns an option to set the HTTP client for the client.
func WithHTTPClient(client *http.Client) func(*Client) error {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

// WithUserAgent returns an option to set the User-Agent for the client.
func WithUserAgent(userAgent string) func(*Client) error {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithBasicAuthentication returns an option to set the username and token for the client.
func WithBasicAuthentication(username, token string) func(*Client) error {
	return func(c *Client) error {
		c.username = username
		c.token = token
		return nil
	}
}

// BasicAuthentication returns the username and token user-id/password pair, encoded using Base64.
// See: https://datatracker.ietf.org/doc/html/rfc7617
func BasicAuthentication(username, token string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, token)))
	return fmt.Sprintf("Basic %s", encoded)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	if c.token != "" {
		req.Header.Set("Authorization", BasicAuthentication(c.username, c.token))
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return err
}
