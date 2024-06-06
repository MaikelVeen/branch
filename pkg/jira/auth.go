package jira

type Credentials struct {
	Email string `json:"email"`
	URL   string `json:"domain"`
	Token string `json:"token"`
}
