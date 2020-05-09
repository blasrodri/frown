package stats

import (
	"github.com/blasrodri/frown/lsof"
)

type Report struct {
	ProcessInfo map[string]map[int]map[lsof.SocketId]*ConnectionReport
}

func NewReport() *Report {
	return &Report{
		ProcessInfo: make(map[string]map[int]map[lsof.SocketId]*ConnectionReport),
	}
}

func (r *Report) AddConnectionReport(processName string, pid int, socketId lsof.SocketId, connRep *ConnectionReport) error {
	if r.ProcessInfo[processName] == nil {
		r.ProcessInfo[processName] = make(map[int]map[lsof.SocketId]*ConnectionReport)
	}
	if r.ProcessInfo[processName][pid] == nil {
		r.ProcessInfo[processName][pid] = make(map[lsof.SocketId]*ConnectionReport)
	}
	r.ProcessInfo[processName][pid][socketId] = connRep
	return nil
}
