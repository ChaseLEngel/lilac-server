package main

import (
	"time"
)

type Request struct {
	Regex        string
	Name         string
	DownloadPath string
	MatchCount   int
	Machines     []Machine
	History      []MatchHistory
}

type MatchHistory struct {
	Timestamp time.Time
	Regex     string
	File      string
}
