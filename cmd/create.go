package cmd

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/MaikelVeen/branch/git"
	"github.com/MaikelVeen/branch/printer"
	"github.com/MaikelVeen/branch/prompt"
	"github.com/MaikelVeen/branch/ticket"
	"github.com/tucnak/climax"
)

// TODO: make this configurable
const baseBranch = "develop"

func GetCreateCommand() climax.Command {
	return climax.Command{
		Name:  "c",
		Brief: "creates a new branch based on a ticket",

		Flags: []climax.Flag{
			{
				Name:     "key",
				Short:    "k",
				Usage:    `--key="."`,
				Help:     `The key/id of the ticket`,
				Variable: true,
			},
		},
		Handle: ExecuteCreateCommand,
	}
}

func ExecuteCreateCommand(ctx climax.Context) int {
	key, ok := ctx.Get("key")
	if !ok {
		printer.Warning("--key --k flag is not optional, example -k=abc-123")
		return 1
	}

	// Get an authenticated ticket system.
	system, err := getSystem()
	if err != nil {
		printer.Error(nil, err)
		return 1
	}

	g := git.NewGitGitCommander()

	// Check the preconditions.
	err = checkPreconditions(key, g, system)
	if err != nil {
		printer.Warning(err.Error())
	}

	printer.Print("Key is valid and working from a clean tree")

	err = checkBaseBranch(g, baseBranch)
	if err != nil {
		printer.Error(nil, err)
	}

	return 0
}

// getSystem returns a ticket system based on the local saved user.
func getSystem() (ticket.TicketSystem, error) {
	// Load the current user from the disk.
	u, err := ticket.LoadFromDisk()
	if err != nil {
		return nil, err
	}

	system, err := GetAuthenticatedTicketSystem(u.System)
	if err != nil {
		return nil, err
	}

	return system, nil
}

// checkPreconditions returns an error when one of the following checks fails:
// validity of the key, in git repo and working tree clean.
func checkPreconditions(key string, g git.GitCommander, s ticket.TicketSystem) error {
	if err := s.ValidateKey(key); err != nil {
		return err
	}

	if err := g.ExecuteStatus(exec.Command); err != nil {
		return errors.New("Checking git status failed, are you in a git repo?")
	}

	if err := g.ExecuteDiffIndex(exec.Command, "HEAD"); err != nil {
		return errors.New("Working tree is not clean, aborting...")
	}

	return nil
}

// checkBaseBranch checks if the configured base branch is currently
// set and ask if the user wants to switch if that is not the case.
func checkBaseBranch(g git.GitCommander, base string) error {
	b, err := g.ExecuteShortSymbolicRef(exec.Command)
	if err != nil {
		return err
	}

	if b != base {
		switchPrompt := prompt.GetConfirmationPrompt("Do you want to switch ? [y/n]", []string{"You are not on the develop branch"})
		val, err := switchPrompt.Run()
		if err != nil {
			return err
		}

		s := strings.ToLower(strings.TrimSpace(val))[0] == 'y'
		if s {
			err := g.ExecuteCheckout(exec.Command, base)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/*import (
	"errors"
	"os/exec"
	"regexp"
	"strings"

	"github.com/MaikelVeen/branch/jira"
	"github.com/fatih/color"
	"github.com/tucnak/climax"
)

func GetBranchCommand() climax.Command {
	return climax.Command{
		Name:  "c",
		Brief: "creates a new branch based on a jira issue",

		Flags: []climax.Flag{
			{
				Name:     "issue",
				Short:    "i",
				Usage:    `--issue="."`,
				Help:     `The id of the jira issue`,
				Variable: true,
			},
		},
		Handle: HandleBranchCommand,
	}
}

// HandleBranchCommand is the main entry point for the new command.
func HandleBranchCommand(ctx climax.Context) int {
	//TODO: replace color.red calls with generic error call with debug option.

	key, ok := ctx.Get("issue")
	if !ok {
		color.Red("issue argument is not optional")
		return 1
	}

	client, err := jira.NewJiraApi(keyRingService, keyRingUser)
	if err != nil {
		color.Red("Could not instantiate jira api client")
		color.Red(err.Error())
		return 1
	}

	if err := ExecGitStatus(); err != nil {
		color.Red("Git status command failed, are you in a git repo?")
		return 1
	}

	if err := ExecCleanTreeCheck(); err != nil {
		color.Red("Working tree is not clean; commit or stash changes and try again")
		return 1
	}

	con, err := ExecBranchCheck()
	if err != nil {
		color.Red("Failed to check what branch you are one")
		return 1
	}

	if !con {
		return 0
	}

	issue, err := client.GetIssue(key)
	if err != nil {
		if errors.Is(err, jira.ErrUnauthorized) {
			color.Red("Invalid credentials")

			return 1
		}

		if errors.Is(err, jira.ErrNotFound) {
			color.Yellow("Issue %s was not found", key)

			return 0
		}
	}

	branchName, err := GetBranchNameFromIssue(issue)
	if err != nil {
		color.Red("Could not build branch name")
		return 1
	}

	if err := ExecBranch(branchName); err != nil {
		color.Red("Could not create branch %s", branchName)
		return 1
	}

	if err := ExecCheckout(branchName); err != nil {
		color.Red("Could not checkout %s", branchName)
		return 1
	}

	color.Green("created and switched to %s", branchName)

	return 0
}

// ExecGitStatus starts the git status command and waits for it to complete.
func ExecGitStatus() error {
	cmd := exec.Command("git", "status")
	return cmd.Run()
}

// ExecBranchCheck checks what the current branch of the user is.
//
// When the user is not on develop he is prompted for a confirmation.
func ExecBranchCheck() (bool, error) {
	out, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").Output()
	if err != nil {
		return false, err
	}

	branch := string(out)
	if branch != "develop" {
		color.Yellow("You are currently not on the develop branch")
		// TODO: Do not ask to continue but to switch.
		c := PromptConfirmation("Do you wish to continue?", 2)

		return c, err
	}

	return true, err
}

// TODO: check if the branch already exists
func ExecBranch(branch string) error {
	cmd := exec.Command("git", "branch", branch)
	return cmd.Run()
}

func ExecCheckout(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	return cmd.Run()
}

// ExecCleanTreeCheck returns an error when the user does
// not have a clean working tree.
func ExecCleanTreeCheck() error {
	cmd := exec.Command("git", "diff-index", "--quiet", "HEAD")
	return cmd.Run()
}

//TODO: this function needs some love
func GetBranchNameFromIssue(issue jira.IssueBean) (string, error) {
	base := getBranchBase(issue)

	// TODO trim whitespace on end.
	filtered, err := removeSpecialChars(issue.Fields.Summary)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.ToLower(filtered), " ") // TODO: limit to ~12 entries
	hyphenated := strings.Join(parts, "-")

	// TODO: check if string would be a valid branch name
	return base + issue.Key + "-" + hyphenated, nil
}

func removeSpecialChars(s string) (string, error) {
	re, err := regexp.Compile(`[^\w!(-\/_ )]`)
	if err != nil {
		return "", nil
	}

	return re.ReplaceAllString(s, ""), nil
}

func getBranchBase(issue jira.IssueBean) string {
	if StringInSliceCaseInsensitive(issue.Fields.Issuetype.Name, []string{"bug"}) {
		return "hotfix/"
	}

	return "feature/"
}
*/
