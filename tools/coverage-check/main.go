package main

import (
	"fmt"
	"os"

	"github.com/goldenm-software/layrz-actions/tools/internal/env"
)

type langCheck struct {
	name       string
	icon       string
	lcovPath   string
	jacocoPath string
	threshold  float64
}

func main() {
	defaultThreshold := env.Float("DEFAULT_THRESHOLD", 80.0)

	langs := []langCheck{
		{
			name:      "C++",
			icon:      "⚙️",
			lcovPath:  os.Getenv("CPP_LCOV_PATH"),
			threshold: env.Float("CPP_THRESHOLD", defaultThreshold),
		},
		{
			name:      "Dart/Flutter",
			icon:      "🎯",
			lcovPath:  os.Getenv("DART_LCOV_PATH"),
			threshold: env.Float("DART_THRESHOLD", defaultThreshold),
		},
		{
			name:      "Go",
			icon:      "🔷",
			lcovPath:  os.Getenv("GO_LCOV_PATH"),
			threshold: env.Float("GO_THRESHOLD", defaultThreshold),
		},
		{
			name:       "Kotlin/Android",
			icon:       "🤖",
			jacocoPath: os.Getenv("KOTLIN_JACOCO_PATH"),
			threshold:  env.Float("KOTLIN_THRESHOLD", defaultThreshold),
		},
		{
			name:      "Python",
			icon:      "🐍",
			lcovPath:  os.Getenv("PYTHON_LCOV_PATH"),
			threshold: env.Float("PYTHON_THRESHOLD", defaultThreshold),
		},
		{
			name:      "Swift",
			icon:      "🍎",
			lcovPath:  os.Getenv("SWIFT_LCOV_PATH"),
			threshold: env.Float("SWIFT_THRESHOLD", defaultThreshold),
		},
	}

	enabledKeys := map[string]string{
		"C++":            "ENABLE_CPP",
		"Dart/Flutter":   "ENABLE_DART",
		"Go":             "ENABLE_GO",
		"Kotlin/Android": "ENABLE_KOTLIN",
		"Python":         "ENABLE_PYTHON",
		"Swift":          "ENABLE_SWIFT",
	}

	failed := false

	for _, l := range langs {
		if os.Getenv(enabledKeys[l.name]) != "true" {
			continue
		}

		var data CoverageData
		if l.jacocoPath != "" {
			data = parseJacoco(l.jacocoPath, l.name, l.icon)
		} else {
			data = parseLcov(l.lcovPath, l.name, l.icon)
		}

		if !data.Available {
			path := l.lcovPath
			if l.jacocoPath != "" {
				path = l.jacocoPath
			}
			fmt.Fprintf(os.Stderr, "%s %s: no coverage data found at %q\n", l.icon, l.name, path)
			failed = true
			continue
		}

		if data.Pct < l.threshold {
			fmt.Fprintf(os.Stderr, "%s %s: coverage %.1f%% is below threshold %.1f%%\n",
				l.icon, l.name, data.Pct, l.threshold)
			failed = true
		} else {
			fmt.Printf("%s %s: coverage %.1f%% meets threshold %.1f%% ✓\n",
				l.icon, l.name, data.Pct, l.threshold)
		}
	}

	if failed {
		os.Exit(1)
	}
}
