package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseJacoco(t *testing.T) {
	// Create a temporary directory for test fixtures
	tmpDir := t.TempDir()

	// Test 1: Valid JaCoCo XML file
	validXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE report PUBLIC "-//JACOCO//DTD Report 1.1//EN" "report.dtd">
<report name="example">
  <counter type="INSTRUCTION" missed="10" covered="90"/>
  <counter type="LINE" missed="5" covered="15"/>
  <counter type="METHOD" missed="2" covered="8"/>
</report>`

	validPath := filepath.Join(tmpDir, "valid.xml")
	if err := os.WriteFile(validPath, []byte(validXML), 0644); err != nil {
		t.Fatalf("failed to write valid XML: %v", err)
	}

	data := parseJacoco(validPath, "Kotlin/Android", "🤖")
	if !data.Available {
		t.Errorf("expected Available=true, got false")
	}
	if data.LinesHit != 15 {
		t.Errorf("expected LinesHit=15, got %d", data.LinesHit)
	}
	if data.LinesTotal != 20 {
		t.Errorf("expected LinesTotal=20, got %d", data.LinesTotal)
	}
	expectedPct := 75.0
	if data.Pct != expectedPct {
		t.Errorf("expected Pct=%.1f, got %.1f", expectedPct, data.Pct)
	}
	if data.Language != "Kotlin/Android" {
		t.Errorf("expected Language=Kotlin/Android, got %s", data.Language)
	}
	if data.Icon != "🤖" {
		t.Errorf("expected Icon=🤖, got %s", data.Icon)
	}

	// Test 2: Missing file
	data = parseJacoco("/nonexistent/path/file.xml", "Kotlin/Android", "🤖")
	if data.Available {
		t.Errorf("expected Available=false for missing file, got true")
	}

	// Test 3: Zero coverage
	zeroXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE report PUBLIC "-//JACOCO//DTD Report 1.1//EN" "report.dtd">
<report name="example">
  <counter type="LINE" missed="10" covered="0"/>
</report>`

	zeroPath := filepath.Join(tmpDir, "zero.xml")
	if err := os.WriteFile(zeroPath, []byte(zeroXML), 0644); err != nil {
		t.Fatalf("failed to write zero XML: %v", err)
	}

	data = parseJacoco(zeroPath, "Kotlin/Android", "🤖")
	if !data.Available {
		t.Errorf("expected Available=true for zero coverage, got false")
	}
	if data.Pct != 0.0 {
		t.Errorf("expected Pct=0.0 for zero coverage, got %.1f", data.Pct)
	}

	// Test 4: Full coverage
	fullXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE report PUBLIC "-//JACOCO//DTD Report 1.1//EN" "report.dtd">
<report name="example">
  <counter type="LINE" missed="0" covered="100"/>
</report>`

	fullPath := filepath.Join(tmpDir, "full.xml")
	if err := os.WriteFile(fullPath, []byte(fullXML), 0644); err != nil {
		t.Fatalf("failed to write full XML: %v", err)
	}

	data = parseJacoco(fullPath, "Kotlin/Android", "🤖")
	if !data.Available {
		t.Errorf("expected Available=true for full coverage, got false")
	}
	if data.Pct != 100.0 {
		t.Errorf("expected Pct=100.0 for full coverage, got %.1f", data.Pct)
	}
}

func TestParseLcov(t *testing.T) {
	// Create a temporary directory for test fixtures
	tmpDir := t.TempDir()

	// Test 1: Valid lcov.info file
	validLcov := `DA:1,1
DA:2,0
DA:3,1
DA:4,1
DA:5,0
end_of_record`

	validPath := filepath.Join(tmpDir, "valid.info")
	if err := os.WriteFile(validPath, []byte(validLcov), 0644); err != nil {
		t.Fatalf("failed to write valid lcov: %v", err)
	}

	data := parseLcov(validPath, "Swift", "🍎")
	if !data.Available {
		t.Errorf("expected Available=true, got false")
	}
	if data.LinesHit != 3 {
		t.Errorf("expected LinesHit=3, got %d", data.LinesHit)
	}
	if data.LinesTotal != 5 {
		t.Errorf("expected LinesTotal=5, got %d", data.LinesTotal)
	}
	expectedPct := 60.0
	if data.Pct != expectedPct {
		t.Errorf("expected Pct=%.1f, got %.1f", expectedPct, data.Pct)
	}

	// Test 2: Missing file
	data = parseLcov("/nonexistent/path/file.info", "Swift", "🍎")
	if data.Available {
		t.Errorf("expected Available=false for missing file, got true")
	}
}
