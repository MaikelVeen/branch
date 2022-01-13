// Package git encapsulates git command and branch verification functionality.
package git

// Git is an interface for executing git commands locally.
// This interface exists to enable mocking for during testing if needed.
type Git interface {
	ExecuteStatus() error
	ExecuteBranch(b string) error
	ExecuteCheckout(b string) error
	ExecuteDiffIndex(t string) error
	ExecuteShowRef(b string) error
}

type GitImplemtation struct{}

func NewGit() Git {
	return &GitImplemtation{}
}

// ExecuteStatus executes `git status` and returns
// and error if the command execution fails.
//
// https://git-scm.com/docs/git-status
func (g *GitImplemtation) ExecuteStatus() error {
	return nil
}

// ExecuteBranch executes `git branch <b>` where b represents
// the branch name. Returns an error if command execution fails.
//
// https://git-scm.com/docs/git-branch
func (g *GitImplemtation) ExecuteBranch(b string) error {
	return nil
}

// ExecuteCheckout executes `git branch <b>` where b represents
// the branch name. Returns an error if command execution fails.
//
// https://git-scm.com/docs/git-checkout
func (g *GitImplemtation) ExecuteCheckout(b string) error {
	return nil
}

// ExecuteDiffIndex compares a tree `t` to the working tree or index.
// Returns an error when there is a diff.
//
// https://git-scm.com/docs/git-diff-index
func (g *GitImplemtation) ExecuteDiffIndex(t string) error {
	return nil
}

// ExecuteShowRef
//
// https://git-scm.com/docs/git-show-ref
func (g *GitImplemtation) ExecuteShowRef(b string) error {
	return nil
}
