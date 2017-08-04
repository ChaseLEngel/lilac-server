package main

import (
	"encoding/json"
	"fmt"
	"github.com/chaselengel/lilac/transfer"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
)

type transferFile struct {
	File string
}

func transferRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	if r.Body == nil {
		res = NewResponse(400, fmt.Errorf("No body"), nil)
	} else {
		res = NewResponse(200, nil, nil)
	}

	var file transferFile
	json.NewDecoder(r.Body).Decode(&file)

	go searchForMatch(file.File)

	json.NewEncoder(w).Encode(res)
}

// Searches all requests across all groups for a match history torrent that matches file.
func searchForMatch(file string) {
	groups, err := allGroups()
	if err != nil {
		log.Errorf("Failed to get all groups: %v", err)
	}

	for _, group := range groups {
		requests, err := group.allRequests()
		if err != nil {
			log.Errorf("Failed to get allRequests for group %v err: %v\n", group.ID, err)
			continue
		}

		for _, request := range requests {
			history, err := request.history()
			if err != nil {
				log.Errorf("Failed to get match history for %v: %v", request.Name, err)
				continue
			}

			// Compare file to request's match history torrent files.
			basename := path.Base(file)
			var matched = false
			for i := 0; !matched && i < len(history); i++ {

				// Compare with torrent file
				if basename == history[i].Name {
					matched = true
					break
				}

				// Search torrent files
				for _, f := range strings.Split(history[i].Files, ",") {
					if f == basename {
						matched = true
					}
				}
			}

			// No match found, continue to next request.
			if !matched {
				continue
			}

			go func() {
				if err := sendToMachines(request, file); err != nil {
					log.Errorf("Failed to send %v err: %v\n", file, err)
				}
			}()
		}
	}
}

// Look up request's requestMachine and start transfer of source file to machines.
func sendToMachines(request Request, source string) error {
	requestMachines, err := request.AllRequestMachines()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, rm := range requestMachines {
		machine, err := findMachine(strconv.FormatUint(uint64(rm.MachineID), 10))
		if err != nil {
			return err
		}
		go func() {
			defer wg.Done()
			if err := transfer.Transfer(source, rm.Destination, machine.Host, machine.Port, machine.User); err != nil {
				log.Errorf("Transfer failed for %v to %v: %v\n", request.Name, machine.Host, err)
			} else {
				log.Infof("Transfer successful for %v to %v:%v\n", source, machine.Host, rm.Destination)
			}
		}()
		wg.Add(1)
	}

	wg.Wait()
	return nil
}
