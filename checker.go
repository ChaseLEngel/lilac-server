package main

import (
	"github.com/chaselengel/lilac/rss"
	"github.com/chaselengel/lilac/telegram"
	"github.com/chaselengel/lilac/transfer"
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
			matchHistory, err := request.history()
			if err != nil {
				log.Error("History error: ", err)
				continue
			}
			for _, mh := range matchHistory {
				re := regexp.MustCompile(request.Regex)
				if re.FindString(item.Title) == re.FindString(mh.File) {
					inmh = true
					break
				}
			}
			// Item is already in history
			if inmh {
				continue
			}

			// Request's download path takes priority over Group's
			var downloadPath string
			if request.DownloadPath != "" {
				downloadPath = request.DownloadPath
			} else if group.DownloadPath != "" {
				downloadPath = group.DownloadPath
			} else {
				log.Error("No download path")
				continue
			}

			filename, err := download(item, downloadPath)
			if err != nil {
				log.Error("Download error: ", err)
				continue
			}

			log.Infof("Downloaded %v to %v\n", item.Title, downloadPath)

			mh := MatchHistory{Timestamp: time.Now(), Regex: "", File: item.Title}
			if err := request.insertMatchHistory(&mh); err != nil {
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

			// If group's settings have auto transfer set
			// then transfer file to request's machines.
			if settings.AutoTransfer {
				if err := send(request, path.Join(downloadPath, filename)); err != nil {
					log.Error("Failed to transfer file:", err)
					continue
				}
			}
		}
	}
}

// Look up request's requestMachine and start transfer of source file to machines.
func send(request Request, source string) error {
	requestMachines, err := request.AllRequestMachines()
	if err != nil {
		log.Errorf("Failed to get %v match history: %v\n", request.Name, err)
		return err
	}
	for _, rm := range requestMachines {
		machine, err := findMachine(strconv.FormatUint(uint64(rm.MachineID), 10))
		if err != nil {
			log.Errorf("Failed to find machine for request %v: %v\n", request.Name, err)
			return err
		}
		if err := transfer.Transfer(source, rm.Destination, machine.Host, machine.Port, machine.User); err != nil {
			log.Errorf("Transfer failed for %v to %v: %v\n", request.Name, machine.Host, err)
			return err
		}
		log.Infof("Transfered %v to %v:%v\n", source, machine.Host, rm.Destination)
	}
	return nil
}

// Download RSS link to Group's Destination.
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
	ioutil.WriteFile(path.Join(destination, filename), body, 0644)
	return filename, nil
}

// Get filename from header
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
