package main

import (
	"os"
	"testing"
)

func writeTempCSV(t *testing.T, content string) string {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatal(err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func TestLoadChalls_FileNotFound(t *testing.T) {
	_, err := loadChalls("/nonexistent/path.csv")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadTeams_FileNotFound(t *testing.T) {
	_, err := loadTeams("/nonexistent/path.csv")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadSolves_FileNotFound(t *testing.T) {
	_, err := loadSolves("/nonexistent/path.csv", Challs{}, nil)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadChalls_InvalidScore(t *testing.T) {
	path := writeTempCSV(t,
		"id,name,description,max_attempts,value,category,type,state,requirements\n"+
			"1,Valid,,0,100,Web,standard,visible,\n"+
			"2,Invalid Score,,0,notanumber,Crypto,standard,visible,\n"+
			"3,Also Valid,,0,200,Web,standard,visible,\n")
	defer os.Remove(path)

	challs, err := loadChalls(path)
	if err != nil {
		t.Fatal(err)
	}

	if challs["1"] == nil {
		t.Error("chall 1 should be loaded")
	}
	if challs["2"] != nil {
		t.Error("chall 2 should be skipped due to invalid score")
	}
	if challs["3"] == nil {
		t.Error("chall 3 should be loaded")
	}
}

func TestLoadChalls_EmptyFile(t *testing.T) {
	path := writeTempCSV(t, "")
	defer os.Remove(path)

	_, err := loadChalls(path)
	if err == nil {
		t.Error("expected error for empty file (no header)")
	}
}

func TestLoadChalls(t *testing.T) {
	path := writeTempCSV(t,
		"id,name,description,max_attempts,value,category,type,state,requirements\n"+
			"1,Test Challenge,,0,100,Web,standard,visible,\n"+
			"2,Another Challenge,,0,200,Crypto,standard,visible,\n")
	defer os.Remove(path)

	challs, err := loadChalls(path)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		id, name, genre string
		score, solver   int
	}{
		{"1", "Test Challenge", "Web", 100, 0},
		{"2", "Another Challenge", "Crypto", 200, 0},
	}
	for _, tt := range tests {
		c := challs[tt.id]
		if c == nil {
			t.Fatalf("chall %s not found", tt.id)
		}
		if c.Name != tt.name {
			t.Errorf("chall %s: name = %q, want %q", tt.id, c.Name, tt.name)
		}
		if c.Genre != tt.genre {
			t.Errorf("chall %s: genre = %q, want %q", tt.id, c.Genre, tt.genre)
		}
		if c.Score != tt.score {
			t.Errorf("chall %s: score = %d, want %d", tt.id, c.Score, tt.score)
		}
		if c.Solver != tt.solver {
			t.Errorf("chall %s: solver = %d, want %d", tt.id, c.Solver, tt.solver)
		}
	}
}

func TestLoadTeams_HiddenAndBanned(t *testing.T) {
	path := writeTempCSV(t,
		"id,oauth_id,name,email,password,secret,website,affiliation,country,bracket_id,hidden,banned,captain_id,created\n"+
			"1,,admin,,hash,,,,,,True,False,1,2026-01-01\n"+
			"2,,team1,,hash,,,,,,False,False,2,2026-01-01\n"+
			"3,,banned_team,,hash,,,,,,False,True,3,2026-01-01\n")
	defer os.Remove(path)

	ignoredTeams, err := loadTeams(path)
	if err != nil {
		t.Fatal(err)
	}

	if !ignoredTeams["1"] {
		t.Error("team 1 (hidden) should be ignored")
	}
	if ignoredTeams["2"] {
		t.Error("team 2 (normal) should not be ignored")
	}
	if !ignoredTeams["3"] {
		t.Error("team 3 (banned) should be ignored")
	}
}

func TestLoadTeams_NoHiddenOrBanned(t *testing.T) {
	path := writeTempCSV(t,
		"id,oauth_id,name,email,password,secret,website,affiliation,country,bracket_id,hidden,banned,captain_id,created\n"+
			"1,,team1,,hash,,,,,,False,False,1,2026-01-01\n"+
			"2,,team2,,hash,,,,,,False,False,2,2026-01-01\n")
	defer os.Remove(path)

	ignoredTeams, err := loadTeams(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(ignoredTeams) != 0 {
		t.Errorf("expected empty set, got %d entries", len(ignoredTeams))
	}
}

func TestLoadSolves_CountSolves(t *testing.T) {
	path := writeTempCSV(t,
		"challenge_id,user_id,team_id,id,ip,provided,type,date\n"+
			"1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n"+
			"1,2,2,2,127.0.0.1,flag1,correct,2026-01-01\n"+
			"2,1,1,3,127.0.0.1,flag2,correct,2026-01-01\n")
	defer os.Remove(path)

	challs := Challs{
		"1": {Name: "C1", Genre: "Web", Score: 100, Solver: 0},
		"2": {Name: "C2", Genre: "Crypto", Score: 200, Solver: 0},
	}

	result, err := loadSolves(path, challs, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result["1"].Solver != 2 {
		t.Errorf("chall 1: solver = %d, want 2", result["1"].Solver)
	}
	if result["2"].Solver != 1 {
		t.Errorf("chall 2: solver = %d, want 1", result["2"].Solver)
	}
}

func TestLoadSolves_ExcludeIgnoredTeams(t *testing.T) {
	path := writeTempCSV(t,
		"challenge_id,user_id,team_id,id,ip,provided,type,date\n"+
			"1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n"+
			"1,2,2,2,127.0.0.1,flag1,correct,2026-01-01\n"+
			"1,3,3,3,127.0.0.1,flag1,correct,2026-01-01\n")
	defer os.Remove(path)

	challs := Challs{
		"1": {Name: "C1", Genre: "Web", Score: 100, Solver: 0},
	}
	ignoredTeams := map[string]bool{"1": true}

	result, err := loadSolves(path, challs, ignoredTeams)
	if err != nil {
		t.Fatal(err)
	}

	if result["1"].Solver != 2 {
		t.Errorf("chall 1: solver = %d, want 2", result["1"].Solver)
	}
}

func TestLoadSolves_SkipUnknownChallenges(t *testing.T) {
	path := writeTempCSV(t,
		"challenge_id,user_id,team_id,id,ip,provided,type,date\n"+
			"1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n"+
			"999,2,2,2,127.0.0.1,flag999,correct,2026-01-01\n")
	defer os.Remove(path)

	challs := Challs{
		"1": {Name: "C1", Genre: "Web", Score: 100, Solver: 0},
	}

	result, err := loadSolves(path, challs, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result["1"].Solver != 1 {
		t.Errorf("chall 1: solver = %d, want 1", result["1"].Solver)
	}
	if result["999"] != nil {
		t.Error("chall 999 should not exist")
	}
}

func TestLoadSolves_DifferentColumnOrder(t *testing.T) {
	// team_id と challenge_id のカラム順序が通常と異なるCSV
	path := writeTempCSV(t,
		"id,team_id,user_id,challenge_id,ip,provided,type,date\n"+
			"1,1,1,1,127.0.0.1,flag1,correct,2026-01-01\n"+
			"2,2,2,1,127.0.0.1,flag1,correct,2026-01-01\n")
	defer os.Remove(path)

	challs := Challs{
		"1": {Name: "C1", Genre: "Web", Score: 100, Solver: 0},
	}

	result, err := loadSolves(path, challs, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result["1"].Solver != 2 {
		t.Errorf("chall 1: solver = %d, want 2", result["1"].Solver)
	}
}

func TestLoadSolves_EmptyFile(t *testing.T) {
	path := writeTempCSV(t, "")
	defer os.Remove(path)

	_, err := loadSolves(path, Challs{}, nil)
	if err == nil {
		t.Error("expected error for empty file (no header)")
	}
}
