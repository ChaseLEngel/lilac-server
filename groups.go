package main

type Group struct {
	DownloadPath  string
	Link          string
	Request       []Request
	History       []MatchHistory
	Constraints   []Constraints
	Notifications []Notification
}
