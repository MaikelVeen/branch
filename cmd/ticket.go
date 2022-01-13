package cmd

import (
	"github.com/MaikelVeen/branch/jira"
	"github.com/MaikelVeen/branch/ticket"
)

const keyRingService = "branch-cli"
const keyRingUser = "branch-cli-anon"

func GetTicketSystem(s ticket.SupportedTicketSystem) ticket.TicketSystem {
	switch s {
	case ticket.Jira:
		return jira.NewJira(keyRingService, keyRingUser)
	}

	return nil
}
