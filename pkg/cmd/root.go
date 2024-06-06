package cmd

import (
	"fmt"
	"os"

	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/MaikelVeen/branch/pkg/ticket"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "branch",
	Short: "branch is a CLI tool with version control enhancements",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCreateCommand().cmd)
	rootCmd.AddCommand(newPullRequestCommand().cmd)
	rootCmd.AddCommand(jira.NewRootCommand().Command)

	lc := newLoginCommand()
	lc.RegisterSystem(ticket.Jira)
	rootCmd.AddCommand(lc.cmd)
}
