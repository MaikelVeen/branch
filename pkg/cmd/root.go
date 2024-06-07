package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "branch",
	Short: "branch is a VSC and Jira swiss army knife",
	Long:  "branch offers multiple commands to make your life easier when working with version control systems and Jira.",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		authCtx, err := auth.LoadUserContext()
		if err != nil {
			if errors.Is(err, auth.ErrAuthContextMissing) {
				return nil
			}
			return err
		}

		ctx := context.WithValue(cmd.Context(), auth.DefaultContextKey, authCtx)
		cmd.SetContext(ctx)

		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
		}),
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewCreateCommand().cmd)
	rootCmd.AddCommand(newPullRequestCommand().cmd)
	rootCmd.AddCommand(jira.NewCommand().Command)
}

func runParentPersistentPreRun(cmd *cobra.Command, args []string) {
	if parent := cmd.Parent(); parent != nil {
		if parent.PersistentPreRun != nil {
			parent.PersistentPreRun(parent, args)
		}
	}
}
