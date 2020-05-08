package main

import (
	"fmt"
	"github.com/blasrodri/frown/lsof"
	"log"
	"time"
)

type connectionsState struct {
	connDeets       map[int]map[lsof.SocketId]*lsof.ConnectionDetails
	processes       map[int]*lsof.Process
	listOpenSockets map[int]map[lsof.SocketId]bool
	socketIdToPid   map[lsof.SocketId]int
}

func newConnectionState() *connectionsState {
	return &connectionsState{
		connDeets:       make(map[int]map[lsof.SocketId]*lsof.ConnectionDetails),
		processes:       make(map[int]*lsof.Process),
		listOpenSockets: make(map[int]map[lsof.SocketId]bool),
		socketIdToPid:   make(map[lsof.SocketId]int),
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
		pid, ok := c.socketIdToPid[connDeet.SocketId]
		if ok {
			c.connDeets[pid][connDeet.SocketId] = connDeet
			c.listOpenSockets[pid][connDeet.SocketId] = true
		}
	}
	// TODO: Calculate hash
}

func (c *connectionsState) setProcesses(processes []*lsof.Process) {
	for _, process := range processes {
		c.processes[process.Pid] = process
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
	result := make([]int, len(c.processes))
	idx := 0
	for k, _ := range c.processes {
		result[idx] = k

	}
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

func manageState() {
	state := newConnectionState()
	processesChan := make(chan []*lsof.Process)
	connectionsChan := make(chan []*lsof.ConnectionDetails)
	go manageProcceses(processesChan)
	go manageConnections(connectionsChan)
	go reportSats(state)
	for {
		time.Sleep(100 * time.Duration(time.Millisecond))
		select {
		case listProcesses := <-processesChan:
			state.setProcesses(listProcesses)
			// remove state associated to dead processes
			go func() {
				for pid, _ := range state.processes {
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
							for openSock, _ := range openSocketsForPid {
								delete(state.socketIdToPid, openSock)
							}
						}
						delete(state.listOpenSockets, pid)
					}
					for _, socketId := range openSocketsPid {
						state.socketIdToPid[socketId] = pid
					}
				}
			}()
		case connDeets := <-connectionsChan:
			state.setConnDetails(connDeets)
		default:
			fmt.Println("Nothing goin on here...")
		}
	}
}

func manageProcceses(processChan chan<- []*lsof.Process) {
	for {
		time.Sleep(500 * time.Duration(time.Millisecond))
		userPids, err := lsof.GetUserProcessList()
		if err != nil {
			log.Fatal(err)
		}
		processChan <- userPids
	}
}

func manageConnections(connectionsChan chan<- []*lsof.ConnectionDetails) {
	for {
		time.Sleep(500 * time.Duration(time.Millisecond))
		connDeets, err := lsof.MonitorUserConnections()
		if err != nil {
			log.Fatal(err)
		}
		connectionsChan <- connDeets
	}
}

func reportSats(c *connectionsState) {
	// TODO: Do something with the data :)
	for {
		time.Sleep(500 * time.Duration(time.Millisecond))
		fmt.Printf("%+v\n", c.connDeets)
	}
}
