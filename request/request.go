package request

import (
	"github.com/chaselengel/lilac/machine"
	"time"
)

type Request struct {
	Regex        string
	Name         string
	DownloadPath string
	MatchCount   int
	Machines     []machine.Machine
}

type MatchHistory struct {
	Timestamp time.Time
	Regex     string
	File      string
}
