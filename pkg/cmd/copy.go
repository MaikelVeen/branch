package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/MaikelVeen/branch/pkg/git"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

type CopyCommand struct {
	Command *cobra.Command

	logger *slog.Logger
	git    *git.Commander
}

func NewCopyCommand() *CopyCommand {
	cc := &CopyCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
		git: git.NewCommander(),
	}

	cc.Command = &cobra.Command{
		Use:     "copy",
		Aliases: []string{"cp"},
		Short:   "Copies the current branch name",
		Args:    cobra.NoArgs,
		RunE:    cc.Execute,
	}

	return cc
}

func (c *CopyCommand) Execute(_ *cobra.Command, _ []string) error {
	branch, err := c.git.Branch(exec.Command, "--show-current")
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(branch))
	c.logger.Info(fmt.Sprintf("%s copied to clipboard", strings.TrimSpace(branch)))
	return nil
}
