package main

import (
	"github.com/blasrodri/frown/lsof"
	"github.com/blasrodri/frown/stats"
	"github.com/blasrodri/frown/ui"
	"log"
	"sync"
	"time"
)

type connectionsState struct {
	connDeets       map[int]map[lsof.SocketId]*lsof.ConnectionDetails
	processes       sync.Map // map[int]*lsof.Process
	listOpenSockets map[int]map[lsof.SocketId]bool
	socketIdToPid   sync.Map // map[lsof.SocketId]int
	mux             sync.Mutex
}

func newConnectionState() *connectionsState {
	return &connectionsState{
		connDeets:       make(map[int]map[lsof.SocketId]*lsof.ConnectionDetails),
		processes:       sync.Map{}, // make(map[int]*lsof.Process),
		listOpenSockets: make(map[int]map[lsof.SocketId]bool),
		socketIdToPid:   sync.Map{},
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
		pid, ok := c.socketIdToPid.Load(connDeet.SocketId)
		if ok {
			_, ok := c.connDeets[pid.(int)]
			if ok {
				c.connDeets[pid.(int)][connDeet.SocketId] = connDeet
				c.listOpenSockets[pid.(int)][connDeet.SocketId] = true
			}
		}
	}
	// TODO: Calculate hash
}

func (c *connectionsState) setProcesses(processes []*lsof.Process) {
	for _, process := range processes {
		c.processes.Store(process.Pid, process)
		// If we have not seen this pid before, then create the map
		// to store its open sockets
		if c.listOpenSockets[process.Pid] == nil {
			c.listOpenSockets[process.Pid] = make(map[lsof.SocketId]bool)

		}
		if c.connDeets[process.Pid] == nil {
			c.connDeets[process.Pid] = make(map[lsof.SocketId]*lsof.ConnectionDetails)

		}
	}
	// TODO: Calculate hash
}

func (c *connectionsState) getAllPIDs() []int {
	numProcessesStores := 0
	c.processes.Range(func(key, value interface{}) bool {
		numProcessesStores += 1
		return true
	})
	result := make([]int, numProcessesStores)

	idx := 0

	c.processes.Range(func(key, value interface{}) bool {
		result[idx] = key.(int)
		idx++
		return true
	})
	return result
}

func (c *connectionsState) setOpenSockets(pid int, listOpenSockets []lsof.SocketId) {
	mOpSock := make(map[lsof.SocketId]bool, len(c.listOpenSockets))
	for _, v := range listOpenSockets {
		mOpSock[v] = true
	}

	for k, _ := range c.listOpenSockets[pid] {
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
				state.processes.Range(func(pid, value interface{}) bool {
					p := &lsof.Process{
						Pid: pid.(int),
					}
					openSocketsPid, err := lsof.ListOpenSockets(p)
					if err != nil {
						// Assume that the pid is dead. Remove it from the state
						delete(state.connDeets, pid.(int))
						state.processes.Delete(pid)
						openSocketsForPid, ok := state.listOpenSockets[pid.(int)]
						if ok {
							for openSock, _ := range openSocketsForPid {
								state.socketIdToPid.Delete(openSock)
							}
						}
						delete(state.listOpenSockets, pid.(int))
					}
					for _, socketId := range openSocketsPid {
						state.socketIdToPid.Store(socketId, pid)
					}
					return true
				})
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
		for pid, sockIdToConnDeets := range c.connDeets {
			if _, present := c.processes.Load(pid); !present {
				continue
			}
			procInfo, _ := c.processes.Load(pid)
			processName := procInfo.(*lsof.Process).Name
			for socketId, connDeets := range sockIdToConnDeets {
				connectionReport, err := stats.AnalyzeSecurity(connDeets)
				if err != nil {
					log.Fatal(err)
				}
				report.AddConnectionReport(processName, pid, socketId, connectionReport)
			}
		}
		c.mux.Unlock()
		if len(report.ProcessInfo) > 0 {
			reportChan <- report
		}

	}
}
