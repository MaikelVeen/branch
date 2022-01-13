package cmd

import (
	"github.com/MaikelVeen/branch/jira"
	"github.com/MaikelVeen/branch/ticket"
)

const keyRingService = "branch-cli"
const keyRingUser = "branch-cli-anon"

func GetNewTicketSystem(s ticket.SupportedTicketSystem) ticket.TicketSystem {
	switch s {
	case ticket.Jira:
		return jira.NewJira(keyRingService, keyRingUser)
	}

	return nil
}

func GetAuthenticatedTicketSystem(s ticket.SupportedTicketSystem) (ticket.TicketSystem, error) {
	switch s {
	case ticket.Jira:
		return jira.NewAuthenticatedJira(keyRingService, keyRingUser)
	}

	return nil, nil
}
