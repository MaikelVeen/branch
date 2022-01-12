package jira

import (
	"encoding/json"
	"net/http"
)

// This resource represents information about the current user, such as basic details, group membership, application roles, preferences, and locale. Use it to
// get, create, update, and delete (restore default) values of the user's preferences and locale.
type MyselfResource interface {
	// GetCurrentUser returns details for the current user.
	GetCurrentUser() (User, error)
}

func UnmarshalUser(data []byte) (User, error) {
	var r User
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *User) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type User struct {
	Self             string           `json:"self"`
	Key              string           `json:"key"`
	AccountID        string           `json:"accountId"`
	AccountType      string           `json:"accountType"`
	Name             string           `json:"name"`
	EmailAddress     string           `json:"emailAddress"`
	AvatarUrls       AvatarUrls       `json:"avatarUrls"`
	DisplayName      string           `json:"displayName"`
	Active           bool             `json:"active"`
	TimeZone         string           `json:"timeZone"`
	Locale           string           `json:"locale"`
	Groups           Groups           `json:"groups"`
	ApplicationRoles ApplicationRoles `json:"applicationRoles"`
	Expand           string           `json:"expand"`
}

type ApplicationRoles struct {
	Size           int64                  `json:"size"`
	Items          []ApplicationRolesItem `json:"items"`
	PagingCallback Callback               `json:"pagingCallback"`
	Callback       Callback               `json:"callback"`
	MaxResults     int64                  `json:"max-results"`
}

type Callback struct {
}

type ApplicationRolesItem struct {
	Key                  string   `json:"key"`
	Groups               []string `json:"groups"`
	Name                 string   `json:"name"`
	DefaultGroups        []string `json:"defaultGroups"`
	SelectedByDefault    bool     `json:"selectedByDefault"`
	Defined              bool     `json:"defined"`
	NumberOfSeats        int64    `json:"numberOfSeats"`
	RemainingSeats       int64    `json:"remainingSeats"`
	UserCount            int64    `json:"userCount"`
	UserCountDescription string   `json:"userCountDescription"`
	HasUnlimitedSeats    bool     `json:"hasUnlimitedSeats"`
	Platform             bool     `json:"platform"`
}

type AvatarUrls struct {
	The16X16 string `json:"16x16"`
	The24X24 string `json:"24x24"`
	The32X32 string `json:"32x32"`
	The48X48 string `json:"48x48"`
}

type Groups struct {
	Size           int64        `json:"size"`
	Items          []GroupsItem `json:"items"`
	PagingCallback Callback     `json:"pagingCallback"`
	Callback       Callback     `json:"callback"`
	MaxResults     int64        `json:"max-results"`
}

type GroupsItem struct {
	Name string `json:"name"`
	Self string `json:"self"`
}

func (c *jiraClient) GetCurrentUser() (User, error) {
	User := User{}

	resp, err := c.B.Call(http.MethodGet, "rest/api/3/myself", c.Email, c.Token, nil, &User)
	if err != nil {
		if resp.StatusCode == http.StatusUnauthorized {
			return User, ErrUnauthorized
		}
	}

	return User, nil
}
