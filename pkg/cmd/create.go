package cmd

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/MaikelVeen/branch/pkg/git"
	"github.com/MaikelVeen/branch/pkg/printer"
	"github.com/MaikelVeen/branch/pkg/ticket"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

const baseBranch = "develop"

type CreateCommand struct {
	cmd *cobra.Command
}

func NewCreateCommand() *CreateCommand {
	cc := &CreateCommand{}

	cc.cmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new git branch based on a ticket identifier",
		RunE:    cc.runCreateCommand,
	}

	return cc
}

func (c *CreateCommand) runCreateCommand(_ *cobra.Command, args []string) error {
	key := args[0]

	// Get an authenticated ticket system.
	system, err := getSystem()
	if err != nil {
		printer.Error(nil, err)
		return err
	}

	commander := git.NewCommander()

	// Check the preconditions.
	err = checkPreconditions(key, commander, system)
	if err != nil {
		printer.Warning(err.Error())
	}

	printer.Print("Key is valid and working from a clean tree")

	err = checkBaseBranch(commander, baseBranch)
	if err != nil {
		printer.Error(nil, err)
		return errors.New("could not check current branch")
	}

	ticket, err := system.Ticket(key)
	if err != nil {
		printer.Error(nil, err)
		return errors.New("could not get ticket")
	}

	base := system.GetBaseFromTicketType(ticket.Type)
	branch := git.GetBranchName(base, ticket.Key, ticket.Title)

	err = checkoutOrCreateBranch(branch, commander)
	if err != nil {
		printer.Error(nil, err)
		return errors.New("could not checkout")
	}

	printer.Success(fmt.Sprintf("checked out %s", branch))

	return nil
}

// getSystem returns a ticket system based on the local saved user.
func getSystem() (ticket.System, error) {
	// Load the current user from the disk.
	u, err := ticket.LoadFromDisk()
	if err != nil {
		return nil, err
	}

	system, err := getAuthenticatedTicketSystem(u.System)
	if err != nil {
		return nil, err
	}

	return system, nil
}

// checkPreconditions returns an error when one of the following checks fails:
// validity of the key, in git repo and working tree clean.
func checkPreconditions(key string, git *git.Commander, s ticket.System) error {
	if err := s.ValidateKey(key); err != nil {
		return err
	}

	if _, err := git.Status(exec.Command); err != nil {
		return errors.New("checking git status failed, are you in a git repo?")
	}

	if err := git.DiffIndex(exec.Command, "HEAD"); err != nil {
		return errors.New("working tree is not clean, aborting")
	}

	return nil
}

// checkBaseBranch checks if the configured base branch is currently
// set and ask if the user wants to switch if that is not the case.
func checkBaseBranch(git *git.Commander, base string) error {
	b, err := git.ShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if b != base {
		var switchBase bool
		switchForm := huh.NewConfirm().
			Title("Switch to base branch?").
			Description("Do you want to switch to the base branch?").
			Value(&switchBase)

		if err := switchForm.Run(); err != nil {
			return err
		}

		if switchBase {
			if err = git.Checkout(exec.Command, base); err != nil {
				return fmt.Errorf("could not checkout the %s branch", base)
			}
		}
	}

	return nil
}

// checkoutOrCreateBranch checks if current branch equals `b`, if true returns nil.
// Then checks if `b` exists, if not creates it and checks it out.
func checkoutOrCreateBranch(b string, git *git.Commander) error {
	current, err := git.ShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if current == b {
		return nil
	}

	// TODO: return pretty errors, or just the errors that the command returns
	err = git.ShowRef(exec.Command, b)
	if err != nil {
		// ShowRef returns error when branch does not exist.
		err = git.Branch(exec.Command, b)
		if err != nil {
			return err
		}
	}

	err = git.Checkout(exec.Command, b)
	if err != nil {
		return err
	}

	return nil
}

func getAuthenticatedTicketSystem(_ ticket.SystemType) (ticket.System, error) {
	panic("not implemented")
}
