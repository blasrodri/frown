package main

import (
	"log"
	"sync"
	"time"

	"github.com/blasrodri/frown/lsof"
	"github.com/blasrodri/frown/stats"
	"github.com/blasrodri/frown/ui"
)

type connectionsState struct {
	connDeets       map[int]map[lsof.SocketID]*lsof.ConnectionDetails
	processes       map[int]*lsof.Process
	listOpenSockets map[int]map[lsof.SocketID]bool
	SocketIDToPid   map[lsof.SocketID]int
	mux             sync.Mutex
}

func newConnectionState() *connectionsState {
	return &connectionsState{
		connDeets:       make(map[int]map[lsof.SocketID]*lsof.ConnectionDetails),
		processes:       make(map[int]*lsof.Process),
		listOpenSockets: make(map[int]map[lsof.SocketID]bool),
		SocketIDToPid:   make(map[lsof.SocketID]int),
	}
}
func (c *connectionsState) getConnDetails(pid int) map[string]*lsof.ConnectionDetails {
	result, ok := c.connDeets[pid]
	if !ok {
		return nil
	}
	return result
}

func (c *connectionsState) setConnDetails(deets []*lsof.ConnectionDetails) {
	for _, connDeet := range deets {
		// Check that the socket id has been mapped to a pid
		// and then add the connection details
		pid, ok := c.SocketIDToPid[connDeet.SocketID]
		if ok {
			_, ok := c.connDeets[pid]
			if ok {
				c.connDeets[pid][connDeet.SocketID] = connDeet
				c.listOpenSockets[pid][connDeet.SocketID] = true
			}
		}
	}
	// TODO: Calculate hash
}

func (c *connectionsState) setProcesses(processes []*lsof.Process) {
	for _, process := range processes {
		if c.processes == nil {
			c.processes = make(map[int]*lsof.Process)
		}
		c.processes[process.Pid] = process
		// If we have not seen this pid before, then create the map
		// to store its open sockets
		if c.listOpenSockets[process.Pid] == nil {
			c.listOpenSockets[process.Pid] = make(map[lsof.SocketID]bool)

		}
		if c.connDeets[process.Pid] == nil {
			c.connDeets[process.Pid] = make(map[lsof.SocketID]*lsof.ConnectionDetails)

		}
	}
	// TODO: Calculate hash
}

func (c *connectionsState) getAllPIDs() []int {
	result := make([]int, len(c.processes))
	idx := 0
	for k := range c.processes {
		result[idx] = k

	}
	return result
}

func (c *connectionsState) setOpenSockets(pid int, listOpenSockets []lsof.SocketID) {
	mOpSock := make(map[lsof.SocketID]bool, len(c.listOpenSockets))
	for _, v := range listOpenSockets {
		mOpSock[v] = true
	}

	for k := range c.listOpenSockets[pid] {
		_, ok := mOpSock[k]
		if !ok {
			delete(c.connDeets[pid], k)
		}
	}
}

func manageState(config *ui.UIConfig, uiFunc func(*ui.UIConfig, <-chan *stats.Report, chan<- bool)) {
	state := newConnectionState()
	processesChan := make(chan []*lsof.Process)
	connectionsChan := make(chan []*lsof.ConnectionDetails)
	reportChan := make(chan *stats.Report)
	closeChan := make(chan bool)
	go manageProcceses(processesChan)
	go manageConnections(state, connectionsChan)
	go reportSats(state, reportChan)
	go uiFunc(config, reportChan, closeChan)

	var shouldStop = false

	go func() {
		shouldStopTemp := <-closeChan
		state.mux.Lock()
		shouldStop = shouldStopTemp
		state.mux.Unlock()
	}()
	var keepRunning = true
	for keepRunning {
		time.Sleep(100 * time.Duration(time.Millisecond))
		select {
		case listProcesses := <-processesChan:
			state.mux.Lock()
			state.setProcesses(listProcesses)
			state.mux.Unlock()
			// remove state associated to dead processes
			go func() {
				state.mux.Lock()
				for pid := range state.processes {
					p := &lsof.Process{
						Pid: pid,
					}
					openSocketsPid, err := lsof.ListOpenSockets(p)
					if err != nil {
						// Assume that the pid is dead. Remove it from the state
						delete(state.connDeets, pid)
						delete(state.processes, pid)
						openSocketsForPid, ok := state.listOpenSockets[pid]
						if ok {
							for openSock := range openSocketsForPid {
								delete(state.SocketIDToPid, openSock)
							}
						}
						delete(state.listOpenSockets, pid)
					}
					for _, SocketID := range openSocketsPid {
						state.SocketIDToPid[SocketID] = pid
					}
				}
				state.mux.Unlock()
			}()
		case connDeets := <-connectionsChan:
			state.mux.Lock()
			state.setConnDetails(connDeets)
			state.mux.Unlock()
		default:
			// Not much to do
		}
		state.mux.Lock()
		keepRunning = !shouldStop
		state.mux.Unlock()
	}
}

func manageProcceses(processChan chan<- []*lsof.Process) {
	for {
		time.Sleep(200 * time.Duration(time.Millisecond))
		userPids, err := lsof.GetUserProcessList()
		if err != nil {
			log.Fatal(err)
		}
		processChan <- userPids
	}
}

func manageConnections(c *connectionsState, connectionsChan chan<- []*lsof.ConnectionDetails) {
	for {
		time.Sleep(200 * time.Duration(time.Millisecond))
		connDeets, err := lsof.MonitorUserConnections()
		if err != nil {
			log.Fatal(err)
		}
		connectionsChan <- connDeets
	}
}

func reportSats(c *connectionsState, reportChan chan<- *stats.Report) {
	for {
		report := stats.NewReport()
		time.Sleep(200 * time.Duration(time.Millisecond))
		c.mux.Lock()
		for pid, sockIDToConnDeets := range c.connDeets {
			if c.processes[pid] == nil {
				continue
			}
			processName := c.processes[pid].Name
			for SocketID, connDeets := range sockIDToConnDeets {
				connectionReport, err := stats.AnalyzeSecurity(connDeets)
				if err != nil {
					log.Fatal(err)
				}
				report.AddConnectionReport(processName, pid, SocketID, connectionReport)
			}
		}
		c.mux.Unlock()
		if len(report.ProcessInfo) > 0 {
			reportChan <- report
		}

	}
}
