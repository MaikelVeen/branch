package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "branch",
	Short: "branch is a VSC and Jira swiss army knife",
	Long:  "branch offers multiple commands to make your life easier when working with version control systems and Jira.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		authCtx, err := auth.LoadUserContext()
		if err != nil {
			return err
		}

		if authCtx == nil {
			return nil
		}

		ctx := context.WithValue(cmd.Context(), auth.DefaultContextKey, authCtx)
		cmd.SetContext(ctx)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewCreateCommand().cmd)
	rootCmd.AddCommand(newPullRequestCommand().cmd)
	rootCmd.AddCommand(jira.NewCommand().Command)
}
