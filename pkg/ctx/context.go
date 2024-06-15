package ctx

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
)

// RootCommandContextWithPreRun traverses backwards from the given command to its root
// command and extracts the context of the root command. It also executes the
// PersistentPreRunE function of the root command if it exists.
func RootCommandContextWithPreRun(cmd *cobra.Command, args []string) (context.Context, error) {
	if cmd == nil {
		return nil, errors.New("cmd is nil")
	}

	// Find the root command
	rootCmd := cmd
	for rootCmd.Parent() != nil {
		rootCmd = rootCmd.Parent()
	}

	// Execute PersistentPreRunE of the root command if it exists
	if rootCmd.PersistentPreRunE != nil {
		if err := rootCmd.PersistentPreRunE(rootCmd, args); err != nil {
			return nil, err
		}
	}

	return rootCmd.Context(), nil
}
