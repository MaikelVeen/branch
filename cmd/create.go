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
			{
				Name:     "base",
				Short:    "b",
				Usage:    `--base="."`,
				Help:     `Overrides the default base branch`,
				Variable: true,
			},
		},
		Handle: ExecuteCreateCommand,
	}
}

func ExecuteCreateCommand(ctx climax.Context) int {
	// TODO: make a better distinction between the default base `develop` and the
	// "base" of the new branch.
	key, ok := ctx.Get("key")
	if !ok {
		printer.Warning("--key --k flag is not optional, example -k=abc-123")
		return 1
	}

	customBranch, customBranchSet := ctx.Get("base")

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
		return 1
	}

	ticket, err := system.GetTicket(key)
	if err != nil {
		printer.Error(nil, err)
		return 1
	}

	// If a custom base is set use that otherwise get from ticket.
	var base string
	if customBranchSet {
		base = customBranch
	} else {
		base = system.GetBaseFromTicketType(ticket.Type)
	}

	branch := git.GetBranchName(base, ticket.Key, ticket.Title)

	err = checkoutOrCreateBranch(branch, g)
	if err != nil {
		printer.Error(nil, err)
		return 1
	}

	printer.Success(fmt.Sprintf("checked out %s", branch))

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

// checkoutOrCreateBranch checks if current branch equals `b`, if true returns nil.
// Then checks if `b` exists, if not creates it and checks it out
func checkoutOrCreateBranch(b string, g git.GitCommander) error {
	current, err := g.ExecuteShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if current == b {
		return nil
	}

	// TODO: return pretty errors, or just the errors that the command returns
	err = g.ExecuteShowRef(exec.Command, b)
	if err != nil {
		// ShowRef returns error when branch does not exist.
		err = g.ExecuteBranch(exec.Command, b)
		if err != nil {
			return err
		}
	}

	err = g.ExecuteCheckout(exec.Command, b)
	if err != nil {
		return err
	}

	return nil
}
