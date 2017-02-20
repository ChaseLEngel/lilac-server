package main

import (
	"fmt"
	"github.com/chaselengel/lilac/rss"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

func check() {
	groups, err := allGroups()
	if err != nil {
		return
	}
	for _, group := range groups {
		channel, err := rss.Get(group.Link)
		if err != nil {
			continue
		}
		requests, err := group.allRequests()
		if err != nil {
			continue
		}
		group.search(channel.Items, requests)
	}
}

func (group Group) search(items []*rss.Item, requests []Request) {
	for _, request := range requests {
		for _, item := range items {
			// Check that RSS item matches Request regular expression
			if matched, err := regexp.MatchString(request.Regex, item.Title); !matched || err != nil {
				continue
			}
			fmt.Printf("Matched %v\n", item.Title)

			// Check that we haven't downloaded this item before.
			inmh := false
			matchHistory, err := request.history()
			if err != nil {
				fmt.Println("History error: ", err)
				continue
			}
			for _, mh := range matchHistory {
				re := regexp.MustCompile(request.Regex)
				fmt.Println("item: ", re.FindString(item.Title))
				fmt.Println("file: ", re.FindString(mh.File))
				if re.FindString(item.Title) == re.FindString(mh.File) {
					inmh = true
					break
				}
			}
			if inmh {
				continue
			}
			fmt.Println("Not in history")

			// Request's download path takes priority over Group's
			var downloadPath string
			if request.DownloadPath != "" {
				downloadPath = request.DownloadPath
			} else if downloadPath != "" {
				downloadPath = group.DownloadPath
			} else {
				continue
			}

			if err := download(item.Link, downloadPath); err != nil {
				fmt.Println("Download error: ", err)
				continue
			}
			fmt.Println("Download")

			mh := MatchHistory{Timestamp: time.Now(), Regex: "", File: item.Title}
			if err := request.insertMatchHistory(&mh); err != nil {
				fmt.Println("Match History Insert error: ", err)
				continue
			}
			fmt.Printf("Inserted into match history: %#v\n", mh)
		}
	}
}

// Download RSS link to Group's Destination.
func download(url, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	filename := filename(resp.Header)
	ioutil.WriteFile(path.Join(destination, filename), body, 0644)
	return nil
}

// Try to get filename from header.
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
