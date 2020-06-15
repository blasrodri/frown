package stats

import (
	"github.com/blasrodri/frown/lsof"
)

type Report struct {
	ProcessInfo map[string]map[int]map[lsof.SocketID]*ConnectionReport
}

func NewReport() *Report {
	return &Report{
		ProcessInfo: make(map[string]map[int]map[lsof.SocketID]*ConnectionReport),
	}
}

func (r *Report) AddConnectionReport(processName string, pid int, SocketID lsof.SocketID, connRep *ConnectionReport) error {
	if r.ProcessInfo[processName] == nil {
		r.ProcessInfo[processName] = make(map[int]map[lsof.SocketID]*ConnectionReport)
	}
	if r.ProcessInfo[processName][pid] == nil {
		r.ProcessInfo[processName][pid] = make(map[lsof.SocketID]*ConnectionReport)
	}
	r.ProcessInfo[processName][pid][SocketID] = connRep
	return nil
}
