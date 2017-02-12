package main

type Group struct {
	Name          string `json:"name"`
	DownloadPath  string `json:"downloadpath"`
	Link          string `json:"link"`
	Request       []Request
	History       []MatchHistory
	Constraints   []Constraint
	Notifications []Notification
}
