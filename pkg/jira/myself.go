package jira

import (
	"context"
	"net/http"
)

type MyselfResourceService struct {
	client *Client
}

func (m *MyselfResourceService) Myself(ctx context.Context) (*User, error) {
	req, err := m.client.NewRequest(ctx, http.MethodGet, "rest/api/3/myself", nil)
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = m.client.Do(req, user)
	if err != nil {
		return nil, err
	}

	return user, nil
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
