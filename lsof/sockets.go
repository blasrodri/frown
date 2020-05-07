package lsof

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func listOpenSockets(p *Process) ([]string, error) {
	listFDNames := make([]string, 0)
	procPidPath := fmt.Sprintf("/proc/%d/fd/", p.Pid)
	files, err := ioutil.ReadDir(procPidPath)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		link, err := os.Readlink(procPidPath + f.Name())
		if err !=nil {
			return nil, err
		}
		listFDNames = append(listFDNames, link)

	}
	nameSockets := filterSocketsOnly(listFDNames[:])
	return nameSockets, nil
}

func filterSocketsOnly(listFDs []string) []string {
	listSockets := make([]string, 0)
	for _, fd := range listFDs {
		if strings.HasPrefix(fd, "socket:[") {
			element := strings.TrimSuffix(strings.TrimPrefix(fd, "socket:["), "]")
			listSockets = append(listSockets, element)
		}
	}
	return listSockets
}
