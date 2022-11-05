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

	if err := c.hasRemote(commander); err != nil {
		return err
	}

	return nil
}

// hasRemote returns an error if the current branch does not have a remote.
func (c *pullRequestCmd) hasRemote(git *Commander) error {
	remote, err := git.Remote(exec.Command)
	if err != nil {
		return fmt.Errorf("could not get git origin: %v", err)

	}

	status, err := git.Status(exec.Command, "-sb")
	if err != nil {
		return fmt.Errorf("could not get git status: %v", err)
	}

	lines := strings.Split(status, "\n")
	if !strings.Contains(lines[0], strings.TrimSpace(remote)) {
		return errors.New("current branch does not have an origin, aborting pr")
	}

	return nil
}
