```
 ________      ________      ________      ________       ________      ___  ___     
|\   __  \    |\   __  \    |\   __  \    |\   ___  \    |\   ____\    |\  \|\  \    
\ \  \|\ /_   \ \  \|\  \   \ \  \|\  \   \ \  \\ \  \   \ \  \___|    \ \  \\\  \   
 \ \   __  \   \ \   _  _\   \ \   __  \   \ \  \\ \  \   \ \  \        \ \   __  \  
  \ \  \|\  \   \ \  \\  \|   \ \  \ \  \   \ \  \\ \  \   \ \  \____    \ \  \ \  \ 
   \ \_______\   \ \__\\ _\    \ \__\ \__\   \ \__\\ \__\   \ \_______\   \ \__\ \__\
    \|_______|    \|__|\|__|    \|__|\|__|    \|__| \|__|    \|_______|    \|__|\|__|
                                                                                     
```
---

branch is a macOS CLI tool with version control addons. 

# Installation

## Manual
To install the command line tool, change to the directory of the folder and run the following command:

```bash
go install
```

The `go install` command places the executable into the $GOPATH/bin directory. The command will place generated executables into a sub-directory of $GOPATH named bin. So please make sure that this directory is in your `$PATH` environment variable.

## Homebrew

The tool is also available on Homebrew using the following commands:

```
brew tap maikelveen/branch
brew install branch
```

# Usage

```
branch is a CLI tool with version control enhancements

Usage:
  branch [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Creates a new git branch based on a ticket identifier
  help        Help about any command
  login       Authenticates with a ticket system.
  pr          Creates a new GitHub Pull Request

Flags:
  -h, --help   help for branch

Use "branch [command] --help" for more information about a command.
```

## Authentication

To use the tool you first have to authenticate with the ticket system. 
### Jira 

When using Jira as ticket system you will need to gather the following information:

- Email
- Domain of your Jira
- API Token (Learn more [here](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/))

Once you have gathered these inputs, you can run the following command and you will be prompted to enter the information.

```bash
branch login 
```

---
## Create a branch

To create a branch based on a ticket you can run the following command:

```bash
branch c -k=key
```

The `k` argument corresponds to the key/identifier of the ticket/issue. Please not that your working tree must be clean for the command to succeed, untracked files are ignored at the moment!