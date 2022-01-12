package jira

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Backend is an interface for making calls against a Jira service.
// This interface exists to enable mocking for during testing if needed.
type Backend interface {
	// Call is the implementation for invoking Jira APIs.
	//
	// Body is marshalled into the request body and the response body is decoded into v.
	Call(method, path, username, password string, body interface{}, v interface{}) (*http.Response, error)
}

func NewBackend(u string) Backend {
	urlParsed, _ := url.Parse(u)

	return &BackendImplementation{
		URL:        urlParsed,
		httpClient: http.DefaultClient,
	}
}

type BackendImplementation struct {
	URL        *url.URL
	httpClient *http.Client
}

// Call is the implementation for invoking Jira APIs.
func (s *BackendImplementation) Call(method, path, username, password string, body interface{}, v interface{}) (*http.Response, error) {
	auth := s.BasicAuth(username, password)

	req, err := s.NewRequest(method, path, auth, "application/json", body)
	if err != nil {
		return nil, err
	}

	return s.Do(req, v)
}

// BasicAuth returns the username and password user-id/password pair, encoded using Base64.
// See: https://datatracker.ietf.org/doc/html/rfc7617
func (s *BackendImplementation) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// NewRequest is used by Call to generate an http.Request. It handles encoding
// parameters and attaching the appropriate headers.
func (s *BackendImplementation) NewRequest(method, path, auth, contentType string, body interface{}) (*http.Request, error) {
	parsedPath, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed creating new request %w", err)
	}

	// TODO: add the body

	fullPath := s.URL.ResolveReference(parsedPath)

	// Body is set later by `Do`.
	req, err := http.NewRequest(method, fullPath.String(), nil)
	if err != nil {

		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// Do is used by Call to execute an API request and parse the response. It uses
// the backend's HTTP client to execute the request and unmarshals the response
// into v.
func (s *BackendImplementation) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return resp, errors.New("error during call")
	}

	// Decode the body into the receiving interface.
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}
