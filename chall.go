package main

// Chall represents a single CTF challenge.
type Chall struct {
	Name   string `json:"name"`
	Genre  string `json:"genre"`
	Score  int    `json:"score"`
	Solver int    `json:"solver"`
}

// Challs maps challenge IDs to their details.
type Challs map[string]*Chall
