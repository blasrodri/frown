package lsof

import (
	"io/ioutil"
	"fmt"
)

type Process struct {
	Pid      int
	Name     string
	Children []Process
}

func newProcess(pid int) (*Process, error) {
	cmdLineBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, err
	}
	return &Process{
		Pid: pid,
		Name: string(cmdLineBytes),
		//TODO complete the rest
	}, nil

}
