package main

import (
	"bufio"
	"encoding/xml"
	"os"
	"strings"
)

// CoverageData represents parsed coverage information for a language.
type CoverageData struct {
	Language   string
	Icon       string
	Pct        float64
	LinesHit   int
	LinesTotal int
	Available  bool
}

// jacocoReport represents the root element of a JaCoCo XML report.
type jacocoReport struct {
	Counters []jacocoCounter `xml:"counter"`
}

// jacocoCounter represents a counter element in a JaCoCo report.
type jacocoCounter struct {
	Type    string `xml:"type,attr"`
	Missed  int    `xml:"missed,attr"`
	Covered int    `xml:"covered,attr"`
}

// parseLcov parses an lcov.info file and returns coverage data.
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

// parseJacoco parses a JaCoCo XML report file and returns coverage data.
func parseJacoco(path, language, icon string) CoverageData {
	d := CoverageData{Language: language, Icon: icon}

	f, err := os.Open(path)
	if err != nil {
		return d
	}
	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.Strict = false
	var report jacocoReport
	if err := decoder.Decode(&report); err != nil {
		return d
	}

	for _, counter := range report.Counters {
		if counter.Type == "LINE" {
			d.LinesHit = counter.Covered
			d.LinesTotal = counter.Missed + counter.Covered
			d.Available = true
			if d.LinesTotal > 0 {
				d.Pct = float64(d.LinesHit) / float64(d.LinesTotal) * 100
			}
			return d
		}
	}

	return d
}
