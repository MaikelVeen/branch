package jira

import "github.com/MaikelVeen/branch/ticket"

type jiraImplemtation struct{}

func NewJira() ticket.TicketSystem {
	return &jiraImplemtation{}
}

func (j *jiraImplemtation) Authenticate(username, password string) error {
	return nil
}

func (j *jiraImplemtation) GetTicketName(key string) (string, error) {
	return "", nil
}

func (j *jiraImplemtation) LoadCredentials(s, u string) interface{} {
	return nil
}

func (j *jiraImplemtation) SaveCredentials(s, u string) interface{} {
	return nil
}

func (j *jiraImplemtation) ValidateKey(key string) error {
	return nil
}
