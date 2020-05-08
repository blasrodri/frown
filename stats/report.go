package stats

import (
)

type Report struct {
	ProcessInfo map[string]map[int]*ConnectionReport
}

func NewReport() *Report {
	return &Report{
		ProcessInfo: make(map[string]map[int]*ConnectionReport),
	}
}

func (r *Report) AddConnectionReport(processName string, pid int, connRep *ConnectionReport) error {
	if r.ProcessInfo[processName] == nil {
		r.ProcessInfo[processName] = make(map[int]*ConnectionReport)
	}
	r.ProcessInfo[processName][pid] = connRep
	return nil
}
