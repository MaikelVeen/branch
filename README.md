# branch

branch is a small CLI tool to automatically create git branches based on JIRA issues.

## Installation

To install the command line tool, change to the directory of the folder and run the following command:

```bash
go install
```

The `go install` command places the executable into the $GOPATH/bin directory. The command will place generated executables into a sub-directory of $GOPATH named bin. So please make sure that this directory is in your `$PATH` environment variable.

## Usage

### Authentication

To use the tool you first have to authenticate with Jira. You need the following to successfully setup the tool:

- Email
- Domain of your Jira
- API Token (Learn more [here](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/))

Once you have gathered these inputs, you can run the following command:

```bash
branch login -e=youremail@test.com -d=test -t=yourtoken
```

### Create branch

To create a branch based on a Jira issue you can run the following command:

```bash
branch n -i=key
```

The `i` argument corresponds to the issue key. Please not that your working tree must be clean for the command to succeed, untracked files are ignored at the moment!

#### Example

![Example usage](https://github.com/github.com/MaikelVeen/branch/raw/main/example.png "Example usage")
