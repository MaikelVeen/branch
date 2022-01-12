package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// This resource represents Jira issues.
type IssueResource interface {
	// GetIssue returns the details for an issue.
	GetIssue(key string) (IssueBean, error)
}

func (c *jiraClient) GetIssue(key string) (IssueBean, error) {
	issue := IssueBean{}

	path := fmt.Sprintf("rest/api/3/issues/%s", key)
	resp, err := c.B.Call(http.MethodGet, path, c.Email, c.Token, nil, &issue)
	if err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			return issue, ErrUnauthorized
		}

		if resp.StatusCode == http.StatusNotFound {
			return issue, ErrNotFound
		}
	}

	return issue, nil
}

func UnmarshalIssueBean(data []byte) (IssueBean, error) {
	var r IssueBean
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *IssueBean) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IssueBean struct {
	Expand string          `json:"expand"`
	ID     string          `json:"id"`
	Self   string          `json:"self"`
	Key    string          `json:"key"`
	Fields IssueBeanFields `json:"fields"`
}

type IssueBeanFields struct {
	Issuetype    Issuetype    `json:"issuetype"`
	Watcher      Watcher      `json:"watcher"`
	Attachment   []Attachment `json:"attachment"`
	SubTasks     []Issuelink  `json:"sub-tasks"`
	Description  Description  `json:"description"`
	Project      Project      `json:"project"`
	Comment      []Comment    `json:"comment"`
	Issuelinks   []Issuelink  `json:"issuelinks"`
	Worklog      []Worklog    `json:"worklog"`
	Updated      int64        `json:"updated"`
	Timetracking Timetracking `json:"timetracking"`
}

type Attachment struct {
	ID             int64            `json:"id"`
	Self           string           `json:"self"`
	Filename       string           `json:"filename"`
	Author         AttachmentAuthor `json:"author"`
	Created        string           `json:"created"`
	Size           int64            `json:"size"`
	MIMEType       string           `json:"mimeType"`
	Content        string           `json:"content"`
	Thumbnail      string           `json:"thumbnail"`
	MediaAPIFileID string           `json:"mediaApiFileId"`
}

type AttachmentAuthor struct {
	Self        string     `json:"self"`
	Key         string     `json:"key"`
	AccountID   string     `json:"accountId"`
	AccountType string     `json:"accountType"`
	Name        string     `json:"name"`
	AvatarUrls  AvatarUrls `json:"avatarUrls"`
	DisplayName string     `json:"displayName"`
	Active      bool       `json:"active"`
}

type Comment struct {
	Self         string        `json:"self"`
	ID           string        `json:"id"`
	Author       AuthorElement `json:"author"`
	Body         Description   `json:"body"`
	UpdateAuthor AuthorElement `json:"updateAuthor"`
	Created      string        `json:"created"`
	Updated      string        `json:"updated"`
	Visibility   Visibility    `json:"visibility"`
}

type AuthorElement struct {
	Self        string `json:"self"`
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
}

type Description struct {
	Type    string               `json:"type"`
	Version int64                `json:"version"`
	Content []DescriptionContent `json:"content"`
}

type DescriptionContent struct {
	Type    string           `json:"type"`
	Content []ContentContent `json:"content"`
}

type ContentContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Visibility struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Issuelink struct {
	ID           string     `json:"id"`
	Type         Type       `json:"type"`
	OutwardIssue *WardIssue `json:"outwardIssue,omitempty"`
	InwardIssue  *WardIssue `json:"inwardIssue,omitempty"`
}

type WardIssue struct {
	ID     string            `json:"id"`
	Key    string            `json:"key"`
	Self   string            `json:"self"`
	Fields InwardIssueFields `json:"fields"`
}

type InwardIssueFields struct {
	Status Status `json:"status"`
}

type Status struct {
	IconURL string `json:"iconUrl"`
	Name    string `json:"name"`
}

type Type struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Inward  string `json:"inward"`
	Outward string `json:"outward"`
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

type Project struct {
	Self            string          `json:"self"`
	ID              string          `json:"id"`
	Key             string          `json:"key"`
	Name            string          `json:"name"`
	AvatarUrls      AvatarUrls      `json:"avatarUrls"`
	ProjectCategory ProjectCategory `json:"projectCategory"`
	Simplified      bool            `json:"simplified"`
	Style           string          `json:"style"`
	Insight         Insight         `json:"insight"`
}

type Insight struct {
	TotalIssueCount     int64  `json:"totalIssueCount"`
	LastIssueUpdateTime string `json:"lastIssueUpdateTime"`
}

type ProjectCategory struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Timetracking struct {
	OriginalEstimate         string `json:"originalEstimate"`
	RemainingEstimate        string `json:"remainingEstimate"`
	TimeSpent                string `json:"timeSpent"`
	OriginalEstimateSeconds  int64  `json:"originalEstimateSeconds"`
	RemainingEstimateSeconds int64  `json:"remainingEstimateSeconds"`
	TimeSpentSeconds         int64  `json:"timeSpentSeconds"`
}

type Watcher struct {
	Self       string          `json:"self"`
	IsWatching bool            `json:"isWatching"`
	WatchCount int64           `json:"watchCount"`
	Watchers   []AuthorElement `json:"watchers"`
}

type Worklog struct {
	Self             string        `json:"self"`
	Author           AuthorElement `json:"author"`
	UpdateAuthor     AuthorElement `json:"updateAuthor"`
	Comment          Description   `json:"comment"`
	Updated          string        `json:"updated"`
	Visibility       Visibility    `json:"visibility"`
	Started          string        `json:"started"`
	TimeSpent        string        `json:"timeSpent"`
	TimeSpentSeconds int64         `json:"timeSpentSeconds"`
	ID               string        `json:"id"`
	IssueID          string        `json:"issueId"`
}
