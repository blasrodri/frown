package lsof

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Process struct {
	Pid     int
	Name    string
	Sockets map[string]bool
}

func newProcess(pid int) (*Process, error) {
	statsBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return nil, err
	}
	statsLines := strings.Split(string(statsBytes), "\n")
	processName := strings.TrimPrefix(string(statsLines[0]), "Name:\t")
	if err != nil {
		return nil, err
	}
	return &Process{
		Pid:  pid,
		Name: processName,
	}, nil

}
