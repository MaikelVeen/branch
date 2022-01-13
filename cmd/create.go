package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/MaikelVeen/branch/git"
	"github.com/MaikelVeen/branch/printer"
	"github.com/MaikelVeen/branch/prompt"
	"github.com/MaikelVeen/branch/ticket"
	"github.com/tucnak/climax"
)

// TODO: make this configurable
const baseBranch = "develop"

func GetCreateCommand() climax.Command {
	return climax.Command{
		Name:  "c",
		Brief: "creates a new branch based on a ticket",

		Flags: []climax.Flag{
			{
				Name:     "key",
				Short:    "k",
				Usage:    `--key="."`,
				Help:     `The key/id of the ticket`,
				Variable: true,
			},
		},
		Handle: ExecuteCreateCommand,
	}
}

func ExecuteCreateCommand(ctx climax.Context) int {
	key, ok := ctx.Get("key")
	if !ok {
		printer.Warning("--key --k flag is not optional, example -k=abc-123")
		return 1
	}

	// Get an authenticated ticket system.
	system, err := getSystem()
	if err != nil {
		printer.Error(nil, err)
		return 1
	}

	g := git.NewGitGitCommander()

	// Check the preconditions.
	err = checkPreconditions(key, g, system)
	if err != nil {
		printer.Warning(err.Error())
	}

	printer.Print("Key is valid and working from a clean tree")

	err = checkBaseBranch(g, baseBranch)
	if err != nil {
		printer.Error(nil, err)
	}

	ticket, err := system.GetTicket(key)
	if err != nil {
		printer.Error(nil, err)
	}

	branch := git.GetBranchName(ticket.Key, system.GetBaseFromTicketType(ticket.Type), ticket.Title)
	printer.Success(branch)

	return 0
}

// getSystem returns a ticket system based on the local saved user.
func getSystem() (ticket.TicketSystem, error) {
	// Load the current user from the disk.
	u, err := ticket.LoadFromDisk()
	if err != nil {
		return nil, err
	}

	system, err := GetAuthenticatedTicketSystem(u.System)
	if err != nil {
		return nil, err
	}

	return system, nil
}

// checkPreconditions returns an error when one of the following checks fails:
// validity of the key, in git repo and working tree clean.
func checkPreconditions(key string, g git.GitCommander, s ticket.TicketSystem) error {
	if err := s.ValidateKey(key); err != nil {
		return err
	}

	if err := g.ExecuteStatus(exec.Command); err != nil {
		return errors.New("Checking git status failed, are you in a git repo?")
	}

	if err := g.ExecuteDiffIndex(exec.Command, "HEAD"); err != nil {
		return errors.New("Working tree is not clean, aborting...")
	}

	return nil
}

// checkBaseBranch checks if the configured base branch is currently
// set and ask if the user wants to switch if that is not the case.
func checkBaseBranch(g git.GitCommander, base string) error {
	b, err := g.ExecuteShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if b != base {
		// Construct a confirmation prompt
		info := fmt.Sprintf("You are not on the %s branch", base)
		switchPrompt := prompt.GetConfirmationPrompt("Do you want to switch ? [y/n]", []string{info})

		// Run the prompt.
		val, err := switchPrompt.Run()
		if err != nil {
			return err
		}

		// If return value is yes, checkout base branch.
		s := strings.ToLower(strings.TrimSpace(val))[0] == 'y'
		if s {
			err := g.ExecuteCheckout(exec.Command, base)
			if err != nil {
				return fmt.Errorf("Could not checkout the %s branch", base)
			}
		}
	}

	return nil
}
