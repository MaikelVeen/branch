package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/MaikelVeen/branch/pkg/cmd/jira"
	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = ".branchconfig"
	envPrefix             = "BRANCH"
)

var rootCmd = &cobra.Command{
	Use:   "branch",
	Short: "branch is a VSC and Jira swiss army knife",
	Long:  "branch offers multiple commands to make your life easier when working with version control systems and Jira.",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if err := initializeConfig(cmd); err != nil {
			return err
		}

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

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath("$HOME")
	v.AddConfigPath(".")

	var cfgNotFoundError viper.ConfigFileNotFoundError
	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &cfgNotFoundError) {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable).
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
