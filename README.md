```
 ________      ________      ________      ________       ________      ___  ___     
|\   __  \    |\   __  \    |\   __  \    |\   ___  \    |\   ____\    |\  \|\  \    
\ \  \|\ /_   \ \  \|\  \   \ \  \|\  \   \ \  \\ \  \   \ \  \___|    \ \  \\\  \   
 \ \   __  \   \ \   _  _\   \ \   __  \   \ \  \\ \  \   \ \  \        \ \   __  \  
  \ \  \|\  \   \ \  \\  \|   \ \  \ \  \   \ \  \\ \  \   \ \  \____    \ \  \ \  \ 
   \ \_______\   \ \__\\ _\    \ \__\ \__\   \ \__\\ \__\   \ \_______\   \ \__\ \__\
    \|_______|    \|__|\|__|    \|__|\|__|    \|__| \|__|    \|_______|    \|__|\|__|
                                                                                     
```
# branch

branch is a small CLI tool to automatically create git branches based on tickets from issue/tickets systems like Jira.

## Installation

To install the command line tool, change to the directory of the folder and run the following command:

```bash
go install
```

The `go install` command places the executable into the $GOPATH/bin directory. The command will place generated executables into a sub-directory of $GOPATH named bin. So please make sure that this directory is in your `$PATH` environment variable.

## Usage

```
branch is a small CLI tool to automatically create git branches based on tickets.

Usage:

        branch command [arguments]

The commands are:

        login       authenticates with ticket system
        c           creates a new branch based on a ticket

Use "branch help [command]" for more information about a command.
```

### Authentication

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

### Create branch

To create a branch based on a Jira issue you can run the following command:

```bash
branch c -k=key
```

The `k` argument corresponds to the issue key. Please not that your working tree must be clean for the command to succeed, untracked files are ignored at the moment!