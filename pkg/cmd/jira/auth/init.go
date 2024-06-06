package auth

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"

	client "github.com/MaikelVeen/branch/pkg/jira"
)

const (
	keyringService = "branch_jira"
	keyringUser    = "branch"
)

type InitCommand struct {
	Command *cobra.Command
	logger  *slog.Logger
}

func NewInitCommand() *InitCommand {
	ac := &InitCommand{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.Kitchen,
			}),
		),
	}

	ac.Command = &cobra.Command{
		Use:   "init",
		Short: "Initialize the Jira authentication",
		RunE:  ac.Execute,
	}

	return ac
}

func (ac *InitCommand) Execute(cmd *cobra.Command, _ []string) error {
	auth := &Context{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your email").
				Value(&auth.EmailAddress),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your Jira subdomain").
				Value(&auth.Subdomain),
		),
		huh.NewGroup(
			huh.NewInput().
				EchoMode(huh.EchoModePassword).
				Title("Enter your API token").
				Description("You can generate this from your Jira account settings").
				Value(&auth.Token),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	baseURL := fmt.Sprintf(client.BaseURLTemplate, auth.Subdomain)
	c, err := client.NewClient(baseURL, client.WithBasicAuthentication(auth.EmailAddress, auth.Token))
	if err != nil {
		return err
	}

	user, err := c.Myself.Myself(cmd.Context())
	if err != nil {
		return err
	}
	auth.User = user

	// Save the user context.
	if err := auth.Save(); err != nil {
		return err
	}

	ac.logger.Info(fmt.Sprintf("Successfully authenticated as %s, saved credentials to keyring", user.DisplayName))
	return nil
}

// Context encapsulates the details needed to authenticate with Jira
// and the user details, fetched after authentication.
type Context struct {
	EmailAddress string       `json:"emailAddress"`
	Subdomain    string       `json:"subdomain"`
	Token        string       `json:"token"`
	User         *client.User `json:"user"` // TODO: Only store display name.
}

// Save saves the user context to the keyring.
func (c *Context) Save() error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := keyring.Set(keyringService, keyringUser, string(jsonData)); err != nil {
		return fmt.Errorf("failed to set keyring: %w", err)
	}

	return nil
}

// Load loads the user context from the keyring.
func LoadUserContext() (*Context, error) {
	jsonData, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get keyring: %w", err)
	}

	var c Context
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keyring data: %w", err)
	}

	return &c, nil
}
