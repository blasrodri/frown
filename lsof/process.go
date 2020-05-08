package lsof

import (
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
)

type Process struct {
	Pid      int
	Name     string
	Parent   *Process
	Sockets  map[string]bool
}

func newProcess(pid int) (*Process, error) {
	statsBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	//cmdLineBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, err
	}
	statsLines := strings.Split(string(statsBytes), "\n")
	parentPid, err := strconv.Atoi(strings.TrimPrefix(string(statsLines[6]), "PPid:\t"))
	processName := strings.TrimPrefix(string(statsLines[0]), "Name:\t")
	if err != nil {
		return nil, err
	}
	var parentProcess *Process
	// this guy has a parent, but we cannot access its data :)
	if pid == 1 {
		parentProcess = nil
	} else {
		parentProcess, err = newProcess(parentPid)
	}
	if err != nil {
		return nil, err
	}
	return &Process{
		Pid: pid,
		Name: processName,
		Parent: parentProcess,
	}, nil

}
