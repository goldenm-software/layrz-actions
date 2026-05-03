package main

import (
	"bufio"
	"os"
	"strings"
)

type CoverageData struct {
	Language    string
	Icon        string
	Pct         float64
	LinesHit    int
	LinesTotal  int
	Available   bool
}

func parseLcov(path, language, icon string) CoverageData {
	d := CoverageData{Language: language, Icon: icon}

	f, err := os.Open(path)
	if err != nil {
		return d
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "DA:") {
			continue
		}
		d.LinesTotal++
		parts := strings.SplitN(line[3:], ",", 2)
		if len(parts) == 2 && parts[1] != "0" {
			d.LinesHit++
		}
	}

	d.Available = true
	if d.LinesTotal > 0 {
		d.Pct = float64(d.LinesHit) / float64(d.LinesTotal) * 100
	}
	return d
}
