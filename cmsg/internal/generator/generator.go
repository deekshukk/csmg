package generator

import (
    "fmt"
    "strings"
    "github/deekshukk/cmsg/internal/git"
)

func GenerateMessage(changes []git.Change) string {
    if len(changes) == 0 {
        return "chore: update files"
    }
    
    // Detect type based on files
    commitType := detectType(changes)
    scope := detectScope(changes)
    
    // Build message
    var message strings.Builder
    message.WriteString(commitType)
    if scope != "" {
        message.WriteString("/" + scope)
    }
    message.WriteString(": ")
    message.WriteString(generateDescription(changes))
    
    return message.String()
}

func detectType(changes []git.Change) string {
    for _, change := range changes {
        if strings.Contains(change.FilePath, "test") {
            return "test"
        }
        if strings.Contains(change.FilePath, "README") || 
           strings.Contains(change.FilePath, ".md") {
            return "docs"
        }
        if change.Status == "A" {
            return "feat"
        }
    }
    return "chore"
}

func detectScope(changes []git.Change) string {
    // Extract common directory
    if len(changes) == 0 {
        return ""
    }
    
    firstPath := changes[0].FilePath
    parts := strings.Split(firstPath, "/")
    if len(parts) > 1 {
        return parts[0]
    }
    return ""
}

func generateDescription(changes []git.Change) string {
    if len(changes) == 1 {
        return fmt.Sprintf("update %s", changes[0].FilePath)
    }
    return fmt.Sprintf("update %d files", len(changes))
}

func GeneratePRDescription(commits []git.Commit, diffOutput string) string {
    if len(commits) == 0 {
        return "No commits found on this branch."
    }

    title := commits[0].Subject

    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("## %s\n\n", title))

    sb.WriteString("### Summary\n\n")
    for _, c := range commits {
        sb.WriteString(fmt.Sprintf("- %s\n", c.Subject))
    }

    changedFiles := parseDiffFiles(diffOutput)
    if len(changedFiles) > 0 {
        sb.WriteString("\n### Files Changed\n\n")
        for _, f := range changedFiles {
            sb.WriteString(fmt.Sprintf("- `%s`\n", f))
        }
    }

    sb.WriteString("\n### Test Plan\n\n")
    sb.WriteString(generateTestPlan(changedFiles))

    return sb.String()
}

func parseDiffFiles(diffOutput string) []string {
    var files []string
    for _, line := range strings.Split(diffOutput, "\n") {
        if line == "" {
            continue
        }
        parts := strings.Fields(line)
        if len(parts) >= 2 {
            files = append(files, parts[1])
        }
    }
    return files
}

func generateTestPlan(files []string) string {
    hasTests := false
    for _, f := range files {
        if strings.Contains(f, "test") {
            hasTests = true
            break
        }
    }

    var steps strings.Builder
    steps.WriteString("- [ ] Run existing test suite\n")
    if !hasTests {
        steps.WriteString("- [ ] Add unit tests for new logic\n")
    }
    steps.WriteString("- [ ] Verify changes locally\n")
    steps.WriteString("- [ ] Check for regressions\n")
    return steps.String()
}