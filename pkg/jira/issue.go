package jira

import (
	"context"
	"fmt"
	"net/http"
)

type IssueResourceService struct {
	client *Client
}

func (i *IssueResourceService) GetIssue(ctx context.Context, key string) (*Issue, error) {
	req, err := i.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("rest/api/3/issue/%s", key), nil)
	if err != nil {
		return nil, err
	}

	issue := new(Issue)
	err = i.client.Do(req, issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

type Issue struct {
	Expand string      `json:"expand"`
	ID     string      `json:"id"`
	Self   string      `json:"self"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Issuetype Issuetype `json:"issuetype"`
	Updated   string    `json:"updated"`
	Summary   string    `json:"summary"`
}

type Issuetype struct {
	Self           string `json:"self"`
	ID             string `json:"id"`
	Description    string `json:"description"`
	IconURL        string `json:"iconUrl"`
	Name           string `json:"name"`
	Subtask        bool   `json:"subtask"`
	AvatarID       int64  `json:"avatarId"`
	HierarchyLevel int64  `json:"hierarchyLevel"`
}
