package issue

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type RootIssueCommand struct {
	Command *cobra.Command
	logger  *slog.Logger
}

func NewCommand() *RootIssueCommand {
	ac := &RootIssueCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	ac.Command = &cobra.Command{
		Use:     "issue",
		Aliases: []string{"i"},
		Short:   "Commands to interact with Jira issues",
	}

	ac.Command.AddCommand(NewGetCommand().Command)
	return ac
}
