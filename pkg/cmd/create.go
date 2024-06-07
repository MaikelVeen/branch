package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/MaikelVeen/branch/pkg/git"
	"github.com/MaikelVeen/branch/pkg/jira"
	"github.com/charmbracelet/huh"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

const baseBranch = "develop"

type CreateCommand struct {
	cmd    *cobra.Command
	logger *slog.Logger
	client *jira.Client
}

func NewCreateCommand() *CreateCommand {
	cc := &CreateCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	cc.cmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new git branch based on a ticket identifier",
		RunE:    cc.runCreateCommand,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd.Context())
			if err != nil {
				cc.logger.Warn("a valid auth context is needed for `create`. Run `branch jira auth init` to authenticate.")
				return err
			}

			cc.client = c
			return nil
		},
	}

	return cc
}

func (c *CreateCommand) runCreateCommand(_ *cobra.Command, args []string) error {
	// key := args[0]

	commander := git.NewCommander()

	err := checkPreconditions(commander)
	if err != nil {
		return err
	}

	err = checkBaseBranch(commander, baseBranch)
	if err != nil {
		return nil
	}

	branch := "test"
	err = checkoutOrCreateBranch(branch, commander)
	if err != nil {
		return err
	}

	c.logger.Info(fmt.Sprintf("checked out %s", branch))
	return nil
}

func checkPreconditions(git *git.Commander) error {
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
