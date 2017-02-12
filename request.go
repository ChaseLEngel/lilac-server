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
}

type MatchHistory struct {
	Timestamp time.Time
	Regex     string
	File      string
}
