# cmsg

A CLI tool that generates commit messages from staged changes and automates GitHub pull request creation with auto-generated descriptions.

## Features

- Generate a conventional commit message from your staged files
- Generate a PR description from your branch commits
- Push a branch and create a draft PR in one command

## Installation

```sh
git clone https://github.com/deekshukk/cmsg
cd cmsg
go build -o cmsg .
```

Move the binary somewhere on your `$PATH`:

```sh
mv cmsg /usr/local/bin/cmsg
```

## Prerequisites

- [Git](https://git-scm.com)
- [GitHub CLI](https://cli.github.com) (`gh`) — required for `cmsg push`
  - Authenticate with `gh auth login` before using

## Usage

### Generate a commit message

Stages your files first, then run:

```sh
git add .
cmsg
```

Output:
```
Analyzing staged changes...
  A internal/generator/generator.go
  M main.go

Suggested commit message:
feat/internal: update 2 files
```

### Preview a PR description

```sh
cmsg --pr
```

Generates a Markdown PR description based on all commits on your branch vs `main`/`master`.

### Push and create a draft PR

```sh
cmsg push
```

This will:
1. Push your branch to `origin`
2. Generate a PR description from your commits
3. Create a **draft** PR on GitHub with the description pre-filled

## How it works

| Command | What it reads | What it produces |
|---|---|---|
| `cmsg` | `git diff --staged` | Conventional commit message |
| `cmsg --pr` | `git log main..HEAD` + `git diff main..HEAD` | Markdown PR description |
| `cmsg push` | Same as `--pr` | Pushes branch + creates draft PR via `gh` |

The PR description includes:
- **Summary** — one bullet per commit
- **Files Changed** — list of modified files
- **Test Plan** — checklist (adds "add unit tests" if no test files were touched)

## Project Structure

```
cmsg/
├── main.go                      # CLI entry point
├── internal/
│   ├── generator/generator.go   # Commit message + PR description logic
│   └── git/git.go               # Git command wrappers
```
