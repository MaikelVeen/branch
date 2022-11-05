package git

// This file uses a modified version of the `fakeExecCommandSuccess` outlined in:
// https://jamiethompson.me/posts/Unit-Testing-Exec-Command-In-Golang/
// Special thanks to him for the method used to test these commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteStatus(t *testing.T) {
	t.Parallel()
	cmd := NewCommander()

	t.Run("shell cmd success returns no err", func(t *testing.T) {
		cmdCtx := getFakeCommand(t, "TestShellProcessSuccess", "git status")
		err := cmd.ExecuteStatus(cmdCtx)

		assert.NoError(t, err)
	})

	t.Run("shell cmd failure returns err", func(t *testing.T) {
		cmdCtx := getFakeCommand(t, "TestShellProcessFail", "git status")
		err := cmd.ExecuteStatus(cmdCtx)

		assert.Error(t, err)
	})
}

func TestExecuteBranch(t *testing.T) {
	t.Parallel()
	cmd := NewCommander()
	b := "feature"

	t.Run("shell cmd success returns no err", func(t *testing.T) {
		exp := fmt.Sprintf("git branch %s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessSuccess", exp)
		err := cmd.ExecuteBranch(cmdCtx, b)

		assert.NoError(t, err)
	})

	t.Run("shell cmd failure returns err", func(t *testing.T) {
		exp := fmt.Sprintf("git branch %s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessFail", exp)
		err := cmd.ExecuteBranch(cmdCtx, b)

		assert.Error(t, err)
	})
}

func TestExecuteCheckout(t *testing.T) {
	t.Parallel()
	cmd := NewCommander()
	b := "feature"

	t.Run("shell cmd success returns no err", func(t *testing.T) {
		exp := fmt.Sprintf("git checkout %s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessSuccess", exp)
		err := cmd.ExecuteCheckout(cmdCtx, b)

		assert.NoError(t, err)
	})

	t.Run("shell cmd failure returns err", func(t *testing.T) {
		exp := fmt.Sprintf("git checkout %s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessFail", exp)
		err := cmd.ExecuteCheckout(cmdCtx, b)

		assert.Error(t, err)
	})
}

func TestExecuteDiffIndex(t *testing.T) {
	t.Parallel()
	cmd := NewCommander()
	b := "HEAD"

	t.Run("shell cmd success returns no err", func(t *testing.T) {
		exp := fmt.Sprintf("git diff-index --quiet %s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessSuccess", exp)
		err := cmd.ExecuteDiffIndex(cmdCtx, b)

		assert.NoError(t, err)
	})

	t.Run("shell cmd failure returns err", func(t *testing.T) {
		exp := fmt.Sprintf("git diff-index --quiet %s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessFail", exp)
		err := cmd.ExecuteDiffIndex(cmdCtx, b)

		assert.Error(t, err)
	})
}

func TestExecuteShowRef(t *testing.T) {
	t.Parallel()
	cmd := NewCommander()
	b := "feature"

	t.Run("shell cmd success returns no err", func(t *testing.T) {
		exp := fmt.Sprintf("git show-ref --verify --quiet refs/heads/%s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessSuccess", exp)
		err := cmd.ExecuteShowRef(cmdCtx, b)

		assert.NoError(t, err)
	})

	t.Run("shell cmd failure returns err", func(t *testing.T) {
		exp := fmt.Sprintf("git show-ref --verify --quiet refs/heads/%s", b)

		cmdCtx := getFakeCommand(t, "TestShellProcessFail", exp)
		err := cmd.ExecuteShowRef(cmdCtx, b)

		assert.Error(t, err)
	})
}

func TestExecuteShortSymbolicRef(t *testing.T) {
	t.Parallel()
	cmd := NewCommander()

	t.Run("shell cmd success returns no err", func(t *testing.T) {
		cmdCtx := getFakeCommand(t, "TestShellProcessSuccessSymbolicRef", "git symbolic-ref --short HEAD")
		branch, err := cmd.ExecuteShortSymbolicRef(cmdCtx)

		assert.NoError(t, err)
		assert.Equal(t, "master", branch)
	})

	t.Run("shell cmd failure returns err", func(t *testing.T) {
		cmdCtx := getFakeCommand(t, "TestShellProcessFail", "git symbolic-ref --short HEAD")
		branch, err := cmd.ExecuteShortSymbolicRef(cmdCtx)

		assert.Error(t, err)
		assert.Empty(t, branch)
	})
}

func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestShellProcessSuccessSymbolicRef(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, "master")
	os.Exit(0)
}

func TestShellProcessFail(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	os.Exit(1)
}

// getFakeCommand is a function that initializes a new exec.Cmd and returns it as closure.
// The command will call the passed shell substitute rather than the command it is provided.
//
// It will also pass through the command and its arguments as an argument to shell substitute
// expectedCommand can be used to check if the passed command is what is expected from the func under test.
//
// Since the closure has access to the outer func, the execContext is able to do asserts on
func getFakeCommand(t *testing.T, shellSub, expectedCommand string) ExecContext {
	// Modified from https://github.com/jthomperoo/test-exec-command-golang/blob/master/funshell/funshell_test.go#L57
	test := fmt.Sprintf("-test.run=%s", shellSub)

	return func(command string, args ...string) *exec.Cmd {
		cs := []string{test, "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)

		commandString := strings.Join(cs[2:], " ")
		assert.Equal(t, expectedCommand, commandString)

		exp := fmt.Sprintf("GO_TEST_EXPECTED_CMD=%s", expectedCommand)
		cmd.Env = []string{"GO_TEST_PROCESS=1", exp}
		return cmd
	}
}
