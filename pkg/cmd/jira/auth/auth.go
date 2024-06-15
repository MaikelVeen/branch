package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"

	client "github.com/MaikelVeen/branch/pkg/jira"
)

var ErrAuthContextMissing = errors.New("auth context not present")

// Command is the parent command for all authentication related commands.
type Command struct {
	Command *cobra.Command
}

func NewCommand() *Command {
	cmd := &Command{}
	cmd.Command = &cobra.Command{
		Use:   "auth",
		Short: "Commands to authenticate with Jira",
	}

	cmd.Command.AddCommand(NewInitCommand().Command)
	cmd.Command.AddCommand(NewShowCommand().Command)
	return cmd
}

// ContextKey is the key used to store the auth context in the context.
type ContextKey string

const (
	DefaultContextKey ContextKey = "default-auth-context"

	keyringService = "branch_jira"
	keyringUser    = "branch"
)

// Context encapsulates the details needed to authenticate with Jira
// and the user details, fetched after authentication.
type Context struct {
	EmailAddress string `json:"emailAddress"`
	Subdomain    string `json:"subdomain"`
	Token        string `json:"token"`
	DisplayName  string `json:"displayName"`
}

// Save saves the user context to the keyring.
func (c *Context) Save() error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err = keyring.Set(keyringService, keyringUser, string(jsonData)); err != nil {
		return fmt.Errorf("failed to set keyring: %w", err)
	}

	return nil
}

// Load loads the user context from the keyring.
func LoadUserContext() (*Context, error) {
	jsonData, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil, ErrAuthContextMissing
		}
		return nil, fmt.Errorf("failed to get keyring: %w", err)
	}

	var c Context
	if err = json.Unmarshal([]byte(jsonData), &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keyring data: %w", err)
	}

	return &c, nil
}

// newClient creates a new Jira client with the given authentication context.
func newClient(authCtx *Context) (*client.Client, error) {
	baseURL := fmt.Sprintf(client.BaseURLTemplate, authCtx.Subdomain)
	c, err := client.NewClient(baseURL, client.WithBasicAuthentication(authCtx.EmailAddress, authCtx.Token))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewClientFromContext creates a new Jira client from the given context.
func NewClientFromContext(ctx context.Context) (*client.Client, error) {
	if authCtx, ok := ctx.Value(DefaultContextKey).(*Context); ok {
		return newClient(authCtx)
	}

	return nil, errors.New("no Jira authentication context found, create one with jira auth init")
}
