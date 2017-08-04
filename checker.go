package main

import (
	"github.com/chaselengel/lilac/rss"
	"github.com/chaselengel/lilac/telegram"
	"github.com/chaselengel/lilac/torrent"
	"github.com/chaselengel/lilac/worker"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var master *worker.Master

// Creates cron jobs for all groups.
func InitChecker(groups []Group) {
	log.Info("Initializing checker")
	master = worker.Init()
	for _, group := range groups {
		settings, err := group.GroupSettings()
		if err != nil {
			log.Error("Failed to get group settings:", err)
			continue
		}
		// Some kind of scoping issue here. If we don't define
		// a new group variable and pass it into func than all cron jobs
		// will execute with the same group.
		localGroupId := group.ID

		log.Infof("Adding slave for group %v(%v), interval: %v\n", group.Name, localGroupId, settings.Interval)
		err = master.AddSlave(int(group.ID), settings.Interval, func() { check(localGroupId) })
		if err != nil {
			log.Error("Failed to add cron for", group.Name)
		}
	}
	master.Start()
}

func check(groupId uint) {
	group, err := findGroup(strconv.FormatUint(uint64(groupId), 10))
	if err != nil {
		log.Error("Failed to find group for id", groupId)
	}

	log.Infof("Running check for %v\n", group.Name)
	err = group.updateLastChecked()
	if err != nil {
		log.Error("Failed to update checked:", err)
		return
	}

	channel, err := rss.Get(group.Link)
	if err != nil {
		log.Error("RSS Get error:", err)
		return
	}

	requests, err := group.allRequests()
	if err != nil {
		log.Error("allRequests error:", err)
		return
	}

	group.search(channel.Items, requests)
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

			// Get all match history for specific request.
			matchHistory, err := request.history()
			if err != nil {
				log.Error("History error: ", err)
				continue
			}

			// Compare match history with RSS item title using regex.
			for _, mh := range matchHistory {
				re := regexp.MustCompile(request.Regex)
				if re.FindString(item.Title) == re.FindString(mh.Title) {
					inmh = true
					break
				}
			}

			// Item is already in history
			if inmh {
				continue
			}

			// Download RSS item to group download path
			filepath, err := download(item, group.DownloadPath)
			if err != nil {
				log.Errorf("Download error: ", err)
				continue
			}

			log.Infof("Downloaded %v to %v\n", item.Title, filepath)

			// Get torrent file information
			tor, err := torrent.Parse(filepath)
			if err != nil {
				log.Errorf("Failed to parse torrent %v: %v", filepath, err)
			}

			var match MatchHistory

			match.Title = item.Title
			match.Name = tor.Info.Name
			match.Timestamp = time.Now()

			// Generate comma seperated list of torrent files.
			// Keep track of file sizes for later
			var size int
			for _, file := range tor.Info.Files {
				size += file.Length
				for _, path := range file.Path {
					if len(match.Files) > 0 {
						path = "," + path
					}
					match.Files += path
				}
			}

			// Should only be true if only one file is in torrent.
			if tor.Info.Length > 0 {
				match.Size = tor.Info.Length
			} else {
				match.Size = size
			}

			if err := request.insertMatchHistory(&match); err != nil {
				log.Error("Match History Insert error: ", err)
				continue
			}

			settings, err := group.GroupSettings()
			if err != nil {
				log.Error("Failed to get group's settings:", err)
				continue
			}

			if settings.TelegramApiToken != "" {
				tele := telegram.New(settings.TelegramApiToken)
				message := strings.Replace(settings.TelegramMessage, "%title%", item.Title, -1)
				if err := tele.SendMessage(settings.TelegramChatId, message); err != nil {
					log.Error("Failed to send Telegram:", err)
				}
			}
		}
	}
}

// Download RSS link to Group's Destination.
// Return path to file that was downloaded.
func download(item *rss.Item, destination string) (string, error) {
	resp, err := http.Get(item.Link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	filename := extractFilename(resp.Header)

	// No filename was found so just use item's Title.
	if filename == "" {
		filename = item.Title
	}

	pathToFile := path.Join(destination, filename)

	ioutil.WriteFile(pathToFile, body, 0644)

	return pathToFile, nil
}

// Attempt to get filename from header
func extractFilename(header http.Header) string {
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
