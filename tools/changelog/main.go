package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gh "github.com/goldenm-software/layrz-actions/tools/internal/github"
)

func main() {
	actor := os.Getenv("GITHUB_ACTOR")
	if os.Getenv("EXCLUDE_DEPENDABOT") == "true" && actor == "dependabot[bot]" {
		os.Exit(0)
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	fromRef, toRef, isPR := determineRefs()

	commits, err := gitLog(fromRef, toRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "git log error: %v\n", err)
	}

	client := gh.NewClient(token, repo)

	authorCache := map[string]string{}

	for i := range commits {
		c := &commits[i]

		full, _ := gitFullHash(c.ShortHash)
		c.FullHash = full

		if login, ok := authorCache[full]; ok {
			c.Author = formatContributor(login, gitAuthorName(c.ShortHash))
		} else {
			login, err := client.GetCommitAuthorLogin(full)
			if err != nil {
				login = ""
			}
			authorCache[full] = login
			c.Author = formatContributor(login, gitAuthorName(c.ShortHash))
		}

		for _, trailer := range extractCoauthorTrailers(c.Body) {
			name, email := parseCoauthor(trailer)
			if email == "" {
				continue
			}
			var login string
			if m := noreplyRegex.FindStringSubmatch(email); m != nil {
				login = m[1]
			} else {
				login, _ = client.SearchUserByEmail(email)
			}
			c.CoAuthors = append(c.CoAuthors, formatContributor(login, name))
		}
	}

	excludeTypes := map[string]bool{}
	for _, t := range strings.Split(os.Getenv("EXCLUDE_TYPES"), ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			excludeTypes[t] = true
		}
	}

	buckets := categorize(commits, excludeTypes, repo)

	header := os.Getenv("COMMENT_HEADER")
	if header == "" {
		header = "## 📋 Changelog Summary"
	}

	intro := "This PR includes the following changes:"
	if !isPR {
		intro = "This release includes the following changes:"
	}

	files, insertions, deletions := gitDiffStats(fromRef, toRef)
	changelog := buildChangelog(buckets, header, intro, files, insertions, deletions)

	fmt.Print(changelog)

	if os.Getenv("POST_COMMENT") == "true" && isPR {
		prNum, err := strconv.Atoi(os.Getenv("PR_NUMBER"))
		if err != nil || prNum == 0 {
			return
		}
		if err := client.PostOrUpdateComment(prNum, "📋 Changelog Summary", changelog); err != nil {
			fmt.Fprintf(os.Stderr, "post comment: %v\n", err)
			os.Exit(1)
		}
	}
}

func determineRefs() (fromRef, toRef string, isPR bool) {
	fromRef = os.Getenv("FROM_REF")
	toRef = os.Getenv("TO_REF")
	isPR = os.Getenv("IS_PR") == "true"

	if fromRef != "" && toRef != "" {
		return
	}

	eventName := os.Getenv("GITHUB_EVENT_NAME")

	if eventName == "pull_request" {
		isPR = true
		fromRef = os.Getenv("PR_BASE_SHA")
		toRef = os.Getenv("PR_HEAD_SHA")
		return
	}

	ref := os.Getenv("GITHUB_REF")
	if strings.HasPrefix(ref, "refs/tags/") {
		currentTag := os.Getenv("GITHUB_REF_NAME")
		toRef = currentTag
		fromRef = gitPreviousTag(currentTag)
		return
	}

	sha := os.Getenv("GITHUB_SHA")
	fromRef = sha
	toRef = sha + "~1"
	return
}
