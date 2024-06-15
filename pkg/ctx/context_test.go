package ctx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MaikelVeen/branch/pkg/ctx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type testContextKey string

func TestRootCommandContextWithPreRun(t *testing.T) {
	t.Parallel()

	rootCmdWithContext := func() *cobra.Command {
		rootCmd := &cobra.Command{
			Use: "root",
		}
		rootCtx := context.WithValue(context.Background(), testContextKey("key"), "rootValue")
		rootCmd.SetContext(rootCtx)
		return rootCmd
	}

	childCmdWithContext := func(rootCmd *cobra.Command) *cobra.Command {
		childCmd := &cobra.Command{
			Use: "child",
		}
		rootCmd.AddCommand(childCmd)
		return childCmd
	}

	grandchildCmdWithContext := func(childCmd *cobra.Command) *cobra.Command {
		grandchildCmd := &cobra.Command{
			Use: "grandchild",
		}
		childCmd.AddCommand(grandchildCmd)
		return grandchildCmd
	}

	grandchildCmdWithPreRun := func(rootCmd *cobra.Command) *cobra.Command {
		grandchildCmd := &cobra.Command{
			Use: "grandchild",
		}
		childCmd := childCmdWithContext(rootCmd)
		childCmd.AddCommand(grandchildCmd)
		return grandchildCmd
	}

	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    context.Context
		wantErr bool
	}{
		{
			name: "NilCommand",
			args: args{
				cmd:  nil,
				args: []string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "RootContext",
			args: args{
				cmd:  rootCmdWithContext(),
				args: []string{},
			},
			want:    rootCmdWithContext().Context(),
			wantErr: false,
		},
		{
			name: "ChildContext",
			args: args{
				cmd:  childCmdWithContext(rootCmdWithContext()),
				args: []string{},
			},
			want:    rootCmdWithContext().Context(),
			wantErr: false,
		},
		{
			name: "GrandchildContext",
			args: args{
				cmd:  grandchildCmdWithContext(childCmdWithContext(rootCmdWithContext())),
				args: []string{},
			},
			want:    rootCmdWithContext().Context(),
			wantErr: false,
		},
		{
			name: "PersistentPreRunE",
			args: args{
				cmd: func() *cobra.Command {
					rootCmd := rootCmdWithContext()
					rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
						rootCmd.SetContext(context.WithValue(rootCmd.Context(), testContextKey("preRunKey"), "preRunValue"))
						return nil
					}
					return grandchildCmdWithPreRun(rootCmd)
				}(),
				args: []string{},
			},
			want: func() context.Context {
				ctx := context.Background()
				ctx = context.WithValue(ctx, testContextKey("key"), "rootValue")
				return context.WithValue(ctx, testContextKey("preRunKey"), "preRunValue")
			}(),
			wantErr: false,
		},
		{
			name: "PersistentPreRunE with error",
			args: args{
				cmd: func() *cobra.Command {
					rootCmd := rootCmdWithContext()
					rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
						return errors.New("error")
					}
					return grandchildCmdWithPreRun(rootCmd)
				}(),
				args: []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ctx.RootCommandContextWithPreRun(tt.args.cmd, tt.args.args)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			require.Equal(t, tt.want, got)
		})
	}
}
