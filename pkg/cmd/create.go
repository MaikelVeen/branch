package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/MaikelVeen/branch/pkg/git"
	"github.com/MaikelVeen/branch/pkg/jira"
	"github.com/charmbracelet/huh"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ArgBase          = "base"
	ArgBaseShort     = "b"
	ArgTemplate      = "template"
	ArgTemplateShort = "t"
)

type CreateCommand struct {
	Command *cobra.Command

	logger *slog.Logger
	git    *git.Commander

	Template   string
	BaseBranch string // TODO: Make configurable.
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

	cc.Command = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a new git branch based on a ticket identifier",
		RunE:    cc.Execute,
	}

	flagset := cc.Command.Flags()

	flagset.StringVarP(
		&cc.Template,
		ArgTemplate,
		ArgTemplateShort,
		"{{.type}}/{{.key}}-{{.summary}}",
		"Template to use for branch name",
	)
	_ = viper.BindPFlag(ArgTemplate, flagset.Lookup(ArgTemplate))

	flagset.StringVarP(
		&cc.BaseBranch,
		ArgBase,
		ArgBaseShort,
		"main",
		"Base branch to create the new branch from",
	)
	_ = viper.BindPFlag(ArgBase, flagset.Lookup(ArgBase))

	return cc
}

func (c *CreateCommand) Execute(cmd *cobra.Command, args []string) error {
	client, err := auth.NewClientFromContext(cmd.Context())
	if err != nil {
		c.logger.Warn("a valid auth context is needed for `create`. Run `branch jira auth init` to authenticate.")
		return err
	}

	if err = c.checkPreconditions(); err != nil {
		return err
	}

	if err = c.checkBaseBranch(c.BaseBranch); err != nil {
		return err
	}

	key := args[0]
	issue, err := client.Issue.GetIssue(cmd.Context(), key)
	if err != nil {
		c.logger.Error(fmt.Errorf("failed to get issue: %w", err).Error())
		return err
	}

	branch, err := BranchNameFromTemplate(c.Template, issue)
	if err != nil {
		return err
	}

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
func (c *CreateCommand) checkoutOrCreateBranch(b string) error {
	current, err := c.git.ShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if current == b {
		return nil
	}

	// TODO: return pretty errors, or just the errors that the command returns
	err = c.git.ShowRef(exec.Command, b)
	if err != nil {
		// ShowRef returns error when branch does not exist.
		_, err = c.git.Branch(exec.Command, b)
		if err != nil {
			return err
		}
	}

	err = c.git.Checkout(exec.Command, b)
	if err != nil {
		return err
	}

	return nil
}

// BranchNameFromTemplate generates a branch name from a given template and Jira issue.
func BranchNameFromTemplate(tmpl string, issue *jira.Issue) (string, error) {
	t, err := template.New("branchName").Parse(tmpl)
	if err != nil {
		return "", err
	}

	params := map[string]string{
		"key":     issue.Key,
		"type":    strings.ToLower(issue.Fields.Issuetype.Name),
		"summary": git.FormatAsValidRef(issue.Fields.Summary),
	}

	var b strings.Builder
	if err = t.Execute(&b, params); err != nil {
		return "", err
	}

	return b.String(), nil
}
