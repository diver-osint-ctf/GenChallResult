package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var challengesPath, solvesPath, teamsPath string
	flag.StringVar(&challengesPath, "c", "", "Path to challenges CSV")
	flag.StringVar(&challengesPath, "challenges", "", "Path to challenges CSV")
	flag.StringVar(&solvesPath, "s", "", "Path to solves CSV")
	flag.StringVar(&solvesPath, "solves", "", "Path to solves CSV")
	flag.StringVar(&teamsPath, "t", "", "Path to teams CSV (optional)")
	flag.StringVar(&teamsPath, "teams", "", "Path to teams CSV (optional)")
	flag.Parse()

	if challengesPath == "" || solvesPath == "" {
		fmt.Println("Usage: genChallResult -c <challenges.csv> -s <solves.csv> [-t <teams.csv>]")
		fmt.Println("Example: genChallResult -c ~/Download/HogeCTF-challenges.csv -s ~/Download/HogeCTF-solves.csv -t ~/Download/HogeCTF-teams.csv")
		os.Exit(1)
	}

	challs, err := loadChalls(challengesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading challenges: %v\n", err)
		os.Exit(1)
	}

	ignoredTeams := make(map[string]bool)
	if teamsPath != "" {
		ignoredTeams, err = loadTeams(teamsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading teams: %v\n", err)
			os.Exit(1)
		}
	}

	challs, err = loadSolves(solvesPath, challs, ignoredTeams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading solves: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll("out", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("out/summary.md", []byte(buildMdTable(challs)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing summary.md: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile("out/sections.md", []byte(buildMdSections(challs)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing sections.md: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile("out/challs.json", []byte(buildJSON(challs)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing challs.json: %v\n", err)
		os.Exit(1)
	}
}
