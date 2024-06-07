package issue

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/MaikelVeen/branch/pkg/cmd/jira/auth"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

type GetCommand struct {
	Command *cobra.Command
	logger  *slog.Logger
}

func NewGetCommand() *GetCommand {
	gc := &GetCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	gc.Command = &cobra.Command{
		Use:   "get",
		Short: "Get a Jira issue",
		RunE:  gc.Execute,
		Args:  cobra.ExactArgs(1),
	}

	return gc
}

func (gc *GetCommand) Execute(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	c, err := auth.NewClientFromContext(ctx)
	if err != nil {
		gc.logger.Error(fmt.Errorf("failed to get client: %w", err).Error())
		return err
	}

	key := args[0]
	issue, err := c.Issue.GetIssue(ctx, key)
	if err != nil {
		gc.logger.Error(fmt.Errorf("failed to get issue: %w", err).Error())
		return err
	}

	// Print issue
	fmt.Println(issue)

	return nil
}
