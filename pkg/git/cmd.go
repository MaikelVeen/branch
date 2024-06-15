// Package git encapsulates git command and branch verification functionality.
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	StatusCmd   string = "status"
	BranchCmd   string = "branch"
	CheckoutCmd string = "checkout"
)

// ExecContext is a function that returns an external command being prepared or run
// in either a real or simulated shell.
type ExecContext = func(command string, args ...string) *exec.Cmd

// Commander implements git commands.
type Commander struct{}

// NewCommander returns a new GitCommander.
func NewCommander() *Commander {
	return &Commander{}
}

func executewithOutput(ctx ExecContext, name string, args ...string) (string, error) {
	cmd := []string{name}
	cmd = append(cmd, args...)

	out, err := ctx("git", cmd...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Status executes `git status` and returns
// the output and an error if the command execution fails.
//
//  https://git-scm.com/docs/git-status
func (g *Commander) Status(ctx ExecContext, args ...string) (string, error) {
	return executewithOutput(ctx, StatusCmd, args...)
}

// Branch executes `git branch <b>` where b represents
// the branch name. Returns an error if command execution fails.
//
// https://git-scm.com/docs/git-branch
func (g *Commander) Branch(ctx ExecContext, b string) error {
	cmd := ctx("git", "branch", b)
	return cmd.Run()
}

// Checkout executes `git branch <b>` where b represents
// the branch name. Returns an error if command execution fails.
//
// https://git-scm.com/docs/git-checkout
func (g *Commander) Checkout(ctx ExecContext, b string) error {
	cmd := ctx("git", "checkout", b)
	return cmd.Run()
}

// DiffIndex compares a tree `t` to the working tree or index.
// Returns an error when there is a diff.
//
// https://git-scm.com/docs/git-diff-index
func (g *Commander) DiffIndex(ctx ExecContext, t string) error {
	cmd := ctx("git", "diff-index", "--quiet", t)
	return cmd.Run()
}

// ShowRef list references in a local repository.
// This function can be used to check if a local branch exists or not.
//
// https://git-scm.com/docs/git-show-ref
func (g *Commander) ShowRef(ctx ExecContext, b string) error {
	pattern := fmt.Sprintf("refs/heads/%s", b)
	cmd := ctx("git", "show-ref", "--verify", "--quiet", pattern)
	return cmd.Run()
}

// ShortSymbolicRef executes `git symbolic-ref --short HEAD`
// Returns which branch head the given symbolic ref refers to and outputs its path as first
// return value. Any error is returned as second return value.
//
// https://git-scm.com/docs/git-symbolic-ref
func (g *Commander) ShortSymbolicRef(ctx ExecContext) (string, error) {
	out, err := ctx("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
