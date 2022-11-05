package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	. "github.com/MaikelVeen/branch/pkg/git"
	"github.com/spf13/cobra"
)

type pullRequestCmd struct {
	cmd *cobra.Command
}

func newPullRequestCommand() *pullRequestCmd {
	pr := &pullRequestCmd{}

	pr.cmd = &cobra.Command{
		Use:   "pr",
		Short: "creates a new PR",
		RunE:  pr.runPullRequestCommand,
	}

	return pr
}

func (c *pullRequestCmd) runPullRequestCommand(cmd *cobra.Command, args []string) error {
	commander := NewCommander()

	if err := c.validOrigin(commander); err != nil {
		return err
	}

	return nil
}

// validOrigin checks wether the current checked out branch
// has an origin. If there is no origin then an error is returned.
func (c *pullRequestCmd) validOrigin(git *Commander) error {
	status, err := git.Status(exec.Command, "-sb")
	if err != nil {
		return fmt.Errorf("could not get git status: %v", err)
	}

	lines := strings.Split(status, "\n")
	if !strings.Contains(lines[0], "origin") {
		return errors.New("current branch does not have an origin, aborting pr")
	}

	return nil
}
