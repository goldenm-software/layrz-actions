package main

import (
	"fmt"
	"os"
	"strconv"

	gh "github.com/goldenm-software/layrz-actions/tools/internal/github"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}
	repo := os.Getenv("GITHUB_REPOSITORY")

	var data []CoverageData

	if os.Getenv("ENABLE_PYTHON") == "true" {
		data = append(data, parseLcov(os.Getenv("PYTHON_LCOV_PATH"), "Python", "🐍"))
	}
	if os.Getenv("ENABLE_GO") == "true" {
		data = append(data, parseLcov(os.Getenv("GO_LCOV_PATH"), "Go", "🔷"))
	}
	if os.Getenv("ENABLE_DART") == "true" {
		data = append(data, parseLcov(os.Getenv("DART_LCOV_PATH"), "Dart/Flutter", "🎯"))
	}
	if os.Getenv("ENABLE_CPP") == "true" {
		data = append(data, parseLcov(os.Getenv("CPP_LCOV_PATH"), "C++", "⚙️"))
	}

	report := buildReport(data)
	fmt.Print(report)

	prNum, err := strconv.Atoi(os.Getenv("PR_NUMBER"))
	if err != nil || prNum == 0 {
		return
	}

	client := gh.NewClient(token, repo)
	if err := client.PostOrUpdateComment(prNum, coverageMarker, report); err != nil {
		fmt.Fprintf(os.Stderr, "post comment: %v\n", err)
		os.Exit(1)
	}
}
