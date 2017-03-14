package main

import (
	"fmt"
	"github.com/chaselengel/lilac/rss"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

// Creates cron jobs for all groups.
func InitChecker(groups []Group) *cron.Cron {
	c := cron.New()
	for _, group := range groups {
		settings, err := group.GroupSettings()
		if err != nil {
			fmt.Println("Failed to get group settings:", err)
			continue
		}
		formated := fmt.Sprintf("0 %v * * * *", settings.Interval)
		var checkgroup []Group
		checkgroup = append(checkgroup, group)
		err = c.AddFunc(formated, func() { check(checkgroup) })
		if err != nil {
			fmt.Printf("Failed to add cron for %v\n", group.Name)
		}
	}
	c.Start()
	return c
}

func check(groups []Group) {
	for _, group := range groups {
		err := group.updateLastChecked()
		if err != nil {
			fmt.Println("Failed to update checked:", err)
			continue
		}
		channel, err := rss.Get(group.Link)
		if err != nil {
			fmt.Println("RSS Get error:", err)
			continue
		}
		requests, err := group.allRequests()
		if err != nil {
			fmt.Println("allRequests error:", err)
			continue
		}
		group.search(channel.Items, requests)
	}
}

// Search RSS items for requests.
func (group Group) search(items []*rss.Item, requests []Request) {
	for _, request := range requests {
		for _, item := range items {
			// Check that RSS item matches Request regular expression
			if matched, err := regexp.MatchString(request.Regex, item.Title); !matched || err != nil {
				continue
			}

			// Check that we haven't downloaded this item before.
			inmh := false
			matchHistory, err := request.history()
			if err != nil {
				fmt.Println("History error: ", err)
				continue
			}
			for _, mh := range matchHistory {
				re := regexp.MustCompile(request.Regex)
				if re.FindString(item.Title) == re.FindString(mh.File) {
					inmh = true
					break
				}
			}
			if inmh {
				fmt.Println("Item already in history")
				continue
			}

			// Request's download path takes priority over Group's
			var downloadPath string
			if request.DownloadPath != "" {
				downloadPath = request.DownloadPath
			} else if group.DownloadPath != "" {
				downloadPath = group.DownloadPath
			} else {
				fmt.Println("No download path")
				continue
			}

			if err := download(item, downloadPath); err != nil {
				fmt.Println("Download error: ", err)
				continue
			}

			mh := MatchHistory{Timestamp: time.Now(), Regex: "", File: item.Title}
			if err := request.insertMatchHistory(&mh); err != nil {
				fmt.Println("Match History Insert error: ", err)
				continue
			}

		}
	}
}

// Download RSS link to Group's Destination.
func download(item *rss.Item, destination string) error {
	resp, err := http.Get(item.Link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	filename := filename(resp.Header)
	// No filename was found so just use item's Title.
	if filename == "" {
		filename = item.Title
	}
	ioutil.WriteFile(path.Join(destination, filename), body, 0644)
	return nil
}

// Get filename from header
func filename(header http.Header) string {
	content := header.Get("Content-Disposition")
	if content == "" {
		return ""
	}
	re := regexp.MustCompile("filename=\".*\"")
	match := re.FindString(content)
	filename := strings.Replace(match, "filename=", "", -1)
	// Return filename with surrounding quotes removed.
	return filename[1 : len(filename)-1]
}
