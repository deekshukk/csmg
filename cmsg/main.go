package main

import (
    "fmt"
    "os"
    "os/exec"
    "github/deekshukk/cmsg/internal/git"
    "github/deekshukk/cmsg/internal/generator"
)

func main() {
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "--pr":
            runPR()
            return
        case "push":
            runPush()
            return
        }
    }
    runCommit()
}

func runCommit() {
    changes, err := git.GetStagedChanges()
    if err != nil {
        fmt.Println("Error: Not in a git repository")
        os.Exit(1)
    }

    if len(changes) == 0 {
        fmt.Println("No staged changes. Use 'git add' first.")
        os.Exit(1)
    }

    fmt.Println("Analyzing staged changes...")
    for _, change := range changes {
        fmt.Printf("  %s %s\n", change.Status, change.FilePath)
    }

    message := generator.GenerateMessage(changes)
    fmt.Printf("\nSuggested commit message:\n%s\n", message)
}

func runPush() {
    branch, err := git.CurrentBranch()
    if err != nil || branch == "" || branch == "main" || branch == "master" {
        fmt.Println("Error: could not determine branch, or you are on main/master.")
        os.Exit(1)
    }

    fmt.Printf("Pushing branch %q to origin...\n", branch)
    pushCmd := exec.Command("git", "push", "-u", "origin", branch)
    pushCmd.Stdout = os.Stdout
    pushCmd.Stderr = os.Stderr
    if err := pushCmd.Run(); err != nil {
        fmt.Println("Error: git push failed.")
        os.Exit(1)
    }

    commits, err := git.GetBranchCommits()
    if err != nil || len(commits) == 0 {
        fmt.Println("Push succeeded. No branch commits found to generate PR description.")
        return
    }

    diffOutput, _ := git.GetBranchDiff()
    description := generator.GeneratePRDescription(commits, diffOutput)
    title := commits[0].Subject

    fmt.Println("Creating PR...")
    prCmd := exec.Command("gh", "pr", "create", "--draft", "--title", title, "--body", description)
    prCmd.Stdout = os.Stdout
    prCmd.Stderr = os.Stderr
    if err := prCmd.Run(); err != nil {
        fmt.Println("Error: gh pr create failed. Is the GitHub CLI installed and authenticated?")
        os.Exit(1)
    }
}

func runPR() {
    commits, err := git.GetBranchCommits()
    if err != nil {
        fmt.Println("Error: could not read branch commits")
        os.Exit(1)
    }

    if len(commits) == 0 {
        fmt.Println("No commits found on this branch compared to main/master.")
        os.Exit(1)
    }

    diffOutput, err := git.GetBranchDiff()
    if err != nil {
        diffOutput = ""
    }

    description := generator.GeneratePRDescription(commits, diffOutput)
    fmt.Printf("Suggested PR description:\n\n%s\n", description)
}