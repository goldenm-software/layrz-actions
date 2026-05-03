package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func gitLog(fromRef, toRef string) ([]Commit, error) {
	// Ensure both refs are available locally
	exec.Command("git", "fetch", "--depth=100", "origin", fromRef).Run()
	exec.Command("git", "fetch", "--depth=100", "origin", toRef).Run()

	fmt.Fprintf(os.Stderr, "debug: git log %s..%s\n", fromRef, toRef)
	out, err := exec.Command("git", "log",
		fmt.Sprintf("%s..%s", fromRef, toRef),
		"--pretty=format:%h%n%s%n%b%n---COMMIT---",
		"--no-merges",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("git log: %w", err)
	}
	fmt.Fprintf(os.Stderr, "debug: git log output length: %d\n", len(out))

	raw := string(out)
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var commits []Commit
	for _, block := range strings.Split(raw, "---COMMIT---") {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.SplitN(block, "\n", 3)
		c := Commit{}
		if len(lines) > 0 {
			c.ShortHash = strings.TrimSpace(lines[0])
		}
		if len(lines) > 1 {
			c.Subject = strings.TrimSpace(lines[1])
		}
		if len(lines) > 2 {
			c.Body = strings.TrimSpace(lines[2])
		}
		if c.ShortHash == "" {
			continue
		}
		commits = append(commits, c)
	}
	return commits, nil
}

func gitFullHash(shortHash string) (string, error) {
	out, err := exec.Command("git", "log", "-1", "--pretty=format:%H", shortHash).Output()
	if err != nil {
		return shortHash, nil
	}
	return strings.TrimSpace(string(out)), nil
}

func gitAuthorName(hash string) string {
	out, _ := exec.Command("git", "log", "-1", "--pretty=format:%an", hash).Output()
	return strings.TrimSpace(string(out))
}

func gitDiffStats(fromRef, toRef string) (files, insertions, deletions int) {
	out, err := exec.Command("git", "diff", "--shortstat", fmt.Sprintf("%s..%s", fromRef, toRef)).Output()
	if err != nil {
		return
	}
	s := string(out)
	fmt.Sscanf(s, " %d file", &files)
	if idx := strings.Index(s, "insertion"); idx >= 0 {
		fmt.Sscanf(s[strings.LastIndex(s[:idx], ",")+1:], " %d insertion", &insertions)
	}
	if idx := strings.Index(s, "deletion"); idx >= 0 {
		fmt.Sscanf(s[strings.LastIndex(s[:idx], ",")+1:], " %d deletion", &deletions)
	}
	return
}

func gitPreviousTag(currentTag string) string {
	out, err := exec.Command("git", "tag", "--sort=-version:refname", "--list", "v*").Output()
	if err != nil {
		return "main"
	}
	tags := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, t := range tags {
		if t == currentTag && i+1 < len(tags) {
			return tags[i+1]
		}
	}
	return "origin/main~1"
}
