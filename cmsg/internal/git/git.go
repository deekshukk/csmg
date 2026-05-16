package git

import (
    "os/exec"
    "strings"
)

type Change struct {
    FilePath string
    Status   string // Added, Modified, Deleted
}

func GetStagedChanges() ([]Change, error) {
    cmd := exec.Command("git", "diff", "--staged", "--name-status")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    
    var changes []Change
    lines := strings.Split(string(output), "\n")
    
    for _, line := range lines {
        if line == "" {
            continue
        }
        parts := strings.Fields(line)
        if len(parts) >= 2 {
            changes = append(changes, Change{
                Status:   parts[0],
                FilePath: parts[1],
            })
        }
    }
    
    return changes, nil
}

func GetDiff() (string, error) {
    cmd := exec.Command("git", "diff", "--staged")
    output, err := cmd.Output()
    return string(output), err
}

type Commit struct {
    Hash    string
    Subject string
}

func GetBranchCommits() ([]Commit, error) {
    base, err := findBase()
    if err != nil {
        return nil, err
    }
    cmd := exec.Command("git", "log", "--pretty=format:%H %s", base+"..HEAD")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    var commits []Commit
    for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
        if line == "" {
            continue
        }
        parts := strings.SplitN(line, " ", 2)
        if len(parts) == 2 {
            commits = append(commits, Commit{Hash: parts[0], Subject: parts[1]})
        }
    }
    return commits, nil
}

func GetBranchDiff() (string, error) {
    base, err := findBase()
    if err != nil {
        return "", err
    }
    cmd := exec.Command("git", "diff", base+"..HEAD", "--name-status")
    output, err := cmd.Output()
    return string(output), err
}

func CurrentBranch() (string, error) {
    cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}

func findBase() (string, error) {
    for _, candidate := range []string{"main", "master"} {
        cmd := exec.Command("git", "rev-parse", "--verify", candidate)
        if err := cmd.Run(); err == nil {
            return candidate, nil
        }
    }
    // Fall back to first commit
    cmd := exec.Command("git", "rev-list", "--max-parents=0", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(output)), nil
}