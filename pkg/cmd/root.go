package cmd

import (
	"fmt"
	"os"

	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "branch",
	Short: "branch is a VSC and Jira swiss army knife",
	Long:  "branch offers multiple commands to make your life easier when working with version control systems and Jira.",
}

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
}
