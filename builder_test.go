package main

import (
	"encoding/json"
	"strings"
	"testing"
)

var emptyChalls = Challs{}

var singleGenreChalls = Challs{
	"c1": {Name: "A", Genre: "Web", Score: 100, Solver: 15},
	"c2": {Name: "B", Genre: "Web", Score: 200, Solver: 3},
	"c3": {Name: "C", Genre: "Web", Score: 150, Solver: 8},
}

var sameSolverChalls = Challs{
	"c1": {Name: "A", Genre: "Crypto", Score: 100, Solver: 5},
	"c2": {Name: "B", Genre: "Crypto", Score: 200, Solver: 5},
}

var sampleChalls = Challs{
	"chall1": {Name: "Challenge 1", Genre: "Crypto", Score: 100, Solver: 10},
	"chall2": {Name: "Challenge 2", Genre: "Web", Score: 200, Solver: 5},
	"chall3": {Name: "Challenge 3", Genre: "Crypto", Score: 150, Solver: 8},
}

func TestBuildMdTable_SortedByGenre(t *testing.T) {
	result := buildMdTable(sampleChalls)
	expected := "| ID | Name | Genre | Score | Solver |\n" +
		"|---|---|---|---|---|\n" +
		"| chall3 | Challenge 3 | Crypto | 150 | 8 |\n" +
		"| chall1 | Challenge 1 | Crypto | 100 | 10 |\n" +
		"| chall2 | Challenge 2 | Web | 200 | 5 |\n"
	if result != expected {
		t.Errorf("got:\n%s\nwant:\n%s", result, expected)
	}
}

func TestBuildMdSections_GroupedByGenre(t *testing.T) {
	result := buildMdSections(sampleChalls)
	expected := "## Crypto\n\n" +
		"### Challenge 3 (150pt / 8 solves)\n\n" +
		"### Challenge 1 (100pt / 10 solves)\n\n" +
		"## Web\n\n" +
		"### Challenge 2 (200pt / 5 solves)\n\n"
	if result != expected {
		t.Errorf("got:\n%s\nwant:\n%s", result, expected)
	}
}

func TestBuildJSON(t *testing.T) {
	result := buildJSON(sampleChalls)
	expected, _ := json.Marshal(sampleChalls)
	if result != string(expected) {
		t.Errorf("got: %s\nwant: %s", result, string(expected))
	}
}

func TestBuildMdTable_EmptyInput(t *testing.T) {
	result := buildMdTable(emptyChalls)
	expected := "| ID | Name | Genre | Score | Solver |\n" +
		"|---|---|---|---|---|\n"
	if result != expected {
		t.Errorf("got:\n%s\nwant:\n%s", result, expected)
	}
}

func TestBuildMdSections_EmptyInput(t *testing.T) {
	result := buildMdSections(emptyChalls)
	if result != "" {
		t.Errorf("expected empty string, got: %s", result)
	}
}

func TestBuildJSON_EmptyInput(t *testing.T) {
	result := buildJSON(emptyChalls)
	if result != "{}" {
		t.Errorf("expected '{}', got: %s", result)
	}
}

func TestBuildMdTable_SingleGenreSortBySolver(t *testing.T) {
	result := buildMdTable(singleGenreChalls)
	expected := "| ID | Name | Genre | Score | Solver |\n" +
		"|---|---|---|---|---|\n" +
		"| c2 | B | Web | 200 | 3 |\n" +
		"| c3 | C | Web | 150 | 8 |\n" +
		"| c1 | A | Web | 100 | 15 |\n"
	if result != expected {
		t.Errorf("got:\n%s\nwant:\n%s", result, expected)
	}
}

func TestBuildMdSections_SingleGenreSortBySolver(t *testing.T) {
	result := buildMdSections(singleGenreChalls)
	expected := "## Web\n\n" +
		"### B (200pt / 3 solves)\n\n" +
		"### C (150pt / 8 solves)\n\n" +
		"### A (100pt / 15 solves)\n\n"
	if result != expected {
		t.Errorf("got:\n%s\nwant:\n%s", result, expected)
	}
}

func TestBuildMdTable_SameSolverCount(t *testing.T) {
	result := buildMdTable(sameSolverChalls)
	lines := strings.Split(result, "\n")
	var dataLines []string
	for _, l := range lines {
		if strings.HasPrefix(l, "| c") {
			dataLines = append(dataLines, l)
		}
	}
	if len(dataLines) != 2 {
		t.Errorf("expected 2 data lines, got %d", len(dataLines))
	}
	hasA := false
	hasB := false
	for _, l := range dataLines {
		if strings.Contains(l, "| A |") {
			hasA = true
		}
		if strings.Contains(l, "| B |") {
			hasB = true
		}
	}
	if !hasA || !hasB {
		t.Errorf("expected both A and B in output, got:\n%s", result)
	}
}
