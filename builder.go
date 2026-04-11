package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type challEntry struct {
	ID    string
	Chall *Chall
}

func buildMdTable(challs Challs) string {
	var sb strings.Builder
	sb.WriteString("| ID | Name | Genre | Score | Solver |\n")
	sb.WriteString("|---|---|---|---|---|\n")

	entries := make([]challEntry, 0, len(challs))
	for id, chall := range challs {
		entries = append(entries, challEntry{ID: id, Chall: chall})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Chall.Genre != entries[j].Chall.Genre {
			return entries[i].Chall.Genre < entries[j].Chall.Genre
		}
		return entries[i].Chall.Solver < entries[j].Chall.Solver
	})

	for _, e := range entries {
		fmt.Fprintf(&sb, "| %s | %s | %s | %d | %d |\n",
			e.ID, e.Chall.Name, e.Chall.Genre, e.Chall.Score, e.Chall.Solver)
	}
	return sb.String()
}

func buildMdSections(challs Challs) string {
	var sb strings.Builder

	genreSet := make(map[string]bool)
	var genres []string
	for _, chall := range challs {
		if !genreSet[chall.Genre] {
			genreSet[chall.Genre] = true
			genres = append(genres, chall.Genre)
		}
	}
	sort.Strings(genres)

	for _, genre := range genres {
		fmt.Fprintf(&sb, "## %s\n\n", genre)

		var genreChalls []challEntry
		for id, chall := range challs {
			if chall.Genre == genre {
				genreChalls = append(genreChalls, challEntry{ID: id, Chall: chall})
			}
		}

		sort.Slice(genreChalls, func(i, j int) bool {
			return genreChalls[i].Chall.Solver < genreChalls[j].Chall.Solver
		})

		for _, e := range genreChalls {
			fmt.Fprintf(&sb, "### %s (%dpt / %d solves)\n\n",
				e.Chall.Name, e.Chall.Score, e.Chall.Solver)
		}
	}
	return sb.String()
}

func buildJSON(challs Challs) (string, error) {
	data, err := json.Marshal(challs)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
