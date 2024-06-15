# Branch 

`branch` is a command line utility that provides useful shortcuts when working with git and Jira.

## Features
- Create Branches: Quickly create branches based on ticket identifiers and templates.
- Jira Integration: Authenticate and interact with Jira from the command line.

# Installation

## Homebrew

```
brew tap maikelveen/branch
brew install branch
```

# Quick start

Authenticate with Jira:

```bash
branch jira auth init
```

Configure the branch template:

```bash
branch config set template "{{.key}}/{{.summary}}"
```

Create a new branch:

```bash
branch create issue-key
```