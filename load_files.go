package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func loadChalls(filePath string) (Challs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	colIdx := make(map[string]int)
	for i, h := range header {
		colIdx[h] = i
	}

	challs := make(Challs)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id := record[colIdx["id"]]
		name := record[colIdx["name"]]
		genre := record[colIdx["category"]]
		score, err := strconv.Atoi(record[colIdx["value"]])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		challs[id] = &Chall{
			Name:   name,
			Genre:  genre,
			Score:  score,
			Solver: 0,
		}
	}
	return challs, nil
}

func loadTeams(filePath string) (map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	colIdx := make(map[string]int)
	for i, h := range header {
		colIdx[h] = i
	}

	ignoredTeams := make(map[string]bool)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		id := record[colIdx["id"]]
		hidden := record[colIdx["hidden"]]
		banned := record[colIdx["banned"]]

		if hidden == "True" || banned == "True" {
			ignoredTeams[id] = true
		}
	}
	return ignoredTeams, nil
}

// loadSolves はsolves CSVを読み込み、各チャレンジのsolver数を更新する。
// solvesはファイルサイズが大きいため、ストリーミング処理を行う。
func loadSolves(filePath string, challs Challs, ignoredTeams map[string]bool) (Challs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		challID := record[0]
		teamID := record[2]

		if ignoredTeams[teamID] {
			continue
		}

		if chall, ok := challs[challID]; ok {
			chall.Solver++
		}
	}
	return challs, nil
}
