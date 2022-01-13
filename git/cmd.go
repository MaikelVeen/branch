// Package git encapsulates git command and branch verification functionality.
package git

import (
	"fmt"
	"os/exec"
)

// Git is an interface for executing git commands locally.
// This interface exists to enable mocking for during testing if needed.
type GitCommander interface {
	ExecuteStatus(ctx ExecContext) error
	ExecuteBranch(ctx ExecContext, b string) error
	ExecuteCheckout(ctx ExecContext, b string) error
	ExecuteDiffIndex(ctx ExecContext, t string) error
	ExecuteShowRef(ctx ExecContext, b string) error
	ExecuteShortSymbolicRef(ctx ExecContext) (string, error)
}

// ExecContext is a function that returns an external command being prepared or run
// in either a real or simulated shell.
type ExecContext = func(command string, args ...string) *exec.Cmd

// gitImplementation is the internal implemtation for executing git commands locally.
type gitImplemtation struct{}

// NewGitGitCommander returns a new GitCommander.
func NewGitGitCommander() GitCommander {
	return &gitImplemtation{}
}

// ExecuteStatus executes `git status` and returns
// and error if the command execution fails.
//
// https://git-scm.com/docs/git-status
func (g *gitImplemtation) ExecuteStatus(ctx ExecContext) error {
	cmd := ctx("git", "status")
	return cmd.Run()
}

// ExecuteBranch executes `git branch <b>` where b represents
// the branch name. Returns an error if command execution fails.
//
// https://git-scm.com/docs/git-branch
func (g *gitImplemtation) ExecuteBranch(ctx ExecContext, b string) error {
	cmd := ctx("git", "branch", b)
	return cmd.Run()
}

// ExecuteCheckout executes `git branch <b>` where b represents
// the branch name. Returns an error if command execution fails.
//
// https://git-scm.com/docs/git-checkout
func (g *gitImplemtation) ExecuteCheckout(ctx ExecContext, b string) error {
	cmd := ctx("git", "checkout", b)
	return cmd.Run()
}

// ExecuteDiffIndex compares a tree `t` to the working tree or index.
// Returns an error when there is a diff.
//
// https://git-scm.com/docs/git-diff-index
func (g *gitImplemtation) ExecuteDiffIndex(ctx ExecContext, t string) error {
	cmd := ctx("git", "diff-index", "--quiet", t)
	return cmd.Run()
}

// ExecuteShowRef list references in a local repository.
// This function can be used to check if a local branch exists or not.
//
// https://git-scm.com/docs/git-show-ref
func (g *gitImplemtation) ExecuteShowRef(ctx ExecContext, b string) error {
	pattern := fmt.Sprintf("refs/heads/%s", b)
	cmd := ctx("git", "show-ref", "--verify", "--quiet", pattern)
	return cmd.Run()
}

// ExecuteShortSymbolicRef executes `git symbolic-ref --short HEAD`
// Returns which branch head the given symbolic ref refers to and outputs its path as first
// return value. Any error is returned as second return value.
//
// https://git-scm.com/docs/git-symbolic-ref
func (g *gitImplemtation) ExecuteShortSymbolicRef(ctx ExecContext) (string, error) {
	out, err := ctx("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
