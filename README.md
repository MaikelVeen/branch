# Branch 

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

`branch` is a command line utility that provides useful shortcuts when working with git and Jira.

# Installation


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
branch config set template "{{issue.key}}/{{issue.summary}}"
```

Create a new branch:

```bash
branch create issue-key
```