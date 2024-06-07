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
	git    *git.Commander
}

func NewCreateCommand() *CreateCommand {
	cc := &CreateCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
		git: git.NewCommander(),
	}

	cc.cmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new git branch based on a ticket identifier",
		RunE:    cc.Execute,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			runParentPersistentPreRun(cmd, args)

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

func (c *CreateCommand) Execute(_ *cobra.Command, _ []string) error {
	// key := args[0]

	err := c.checkPreconditions()
	if err != nil {
		return err
	}

	if err = c.checkBaseBranch(baseBranch); err != nil {
		return err
	}

	// TODO: Get issue and construct branch name.
	branch := "test"
	if err = c.checkoutOrCreateBranch(branch); err != nil {
		return err
	}

	c.logger.Info(fmt.Sprintf("checked out %s", branch))
	return nil
}

func (c *CreateCommand) checkPreconditions() error {
	if _, err := c.git.Status(exec.Command); err != nil {
		return errors.New("checking git status failed, are you in a git repo?")
	}

	if err := c.git.DiffIndex(exec.Command, "HEAD"); err != nil {
		return errors.New("working tree is not clean, aborting")
	}

	return nil
}

// checkBaseBranch checks if the configured base branch is currently
// set and ask if the user wants to switch if that is not the case.
func (c *CreateCommand) checkBaseBranch(base string) error {
	b, err := c.git.ShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if b != base {
		var switchBase bool
		if err = huh.NewConfirm().
			Title("Switch to base branch?").
			Description("Do you want to switch to the base branch?").
			Value(&switchBase).
			Run(); err != nil {
			return err
		}

		if switchBase {
			if err = c.git.Checkout(exec.Command, base); err != nil {
				return fmt.Errorf("could not checkout the %s branch", base)
			}
		}
	}

	return nil
}

// checkoutOrCreateBranch checks if current branch equals `b`, if true returns nil.
// Then checks if `b` exists, if not creates it and checks it out.
func (cmd *CreateCommand) checkoutOrCreateBranch(b string) error {
	current, err := cmd.git.ShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if current == b {
		return nil
	}

	// TODO: return pretty errors, or just the errors that the command returns
	err = cmd.git.ShowRef(exec.Command, b)
	if err != nil {
		// ShowRef returns error when branch does not exist.
		err = cmd.git.Branch(exec.Command, b)
		if err != nil {
			return err
		}
	}

	err = cmd.git.Checkout(exec.Command, b)
	if err != nil {
		return err
	}

	return nil
}
