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

type File struct {
	File string
}

func transferRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Response

	groups, err := allGroups()
	if err != nil {
		res = Response{Status{400, err.Error()}, nil}
	}

	if r.Body == nil {
		res = Response{Status{400, "No body"}, nil}
	}

	var file File
	json.NewDecoder(r.Body).Decode(&file)

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
			basename := path.Base(file.File)
			fmt.Println(basename)
			var matched = false
			for i := 0; !matched && i < len(history); i++ {

				// Compare with torrent filename
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
				if err := send(request, file.File); err != nil {
					log.Errorf("Failed to send %v err: %v\n", file.File, err)
				}
			}()
		}
	}
	if err == nil {
		res = Response{Status{200, ""}, nil}
	}
	json.NewEncoder(w).Encode(res)
}

// Look up request's requestMachine and start transfer of source file to machines.
func send(request Request, source string) error {
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
