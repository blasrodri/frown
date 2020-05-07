package lsof

import ()

type Process struct {
	Pid      int
	Name     string
	Children []Process
}

func NewProcess(pid int) *Process {
	return &Process{
		Pid: pid,
		//TODO complete the rest
	}

}
