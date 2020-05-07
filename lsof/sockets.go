package lsof

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func listOpenSockets(p *Process) ([]string, error) {
	fmt.Printf("Process: %+v\n", p)
	listFDNames := make([]string, 0)
	procPidPath := fmt.Sprintf("/proc/%d/fd/", p.Pid)
	files, err := ioutil.ReadDir(procPidPath)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		link, err := os.Readlink(procPidPath + f.Name())
		// TODO: Verify whether this makes any sense at all
		if link == "" {
			continue
		}
		if err != nil {
			return nil, err
		}
		listFDNames = append(listFDNames, link)

	}
	nameSockets := filterSocketsOnly(listFDNames[:])
	return nameSockets, nil
}

func filterSocketsOnly(listFDs []string) []string {
	listSockets := make([]string, 0)
	for _, fd := range listFDs {
		if strings.HasPrefix(fd, "socket:[") {
			element := strings.TrimSuffix(strings.TrimPrefix(fd, "socket:["), "]")
			listSockets = append(listSockets, element)
		}
	}
	return listSockets
}

func MonitorUserConnections() ([]*connectionDetails, error) {
	fileInfoBytes, err := ioutil.ReadFile("/proc/net/tcp")
	if err != nil {
		return nil, err
	}
	fileInfoStr := string(fileInfoBytes)
	connectionListStr := strings.Split(fileInfoStr, "\n")
	if len(connectionListStr) < 2 {
		return nil, errors.New("There are no open connections at the moment")
	}
	// the first line is a header that we do not care about
	openConnections := connectionListStr[1:len(connectionListStr) - 1]
	openConnectionsResult := make([]*connectionDetails, len(openConnections))
	for i, line := range openConnections {
		fields := strings.Fields(line)
		if len(fields) < 12 {
			return nil, errors.New("There are not enough attributes in line " + line)
		}
		connectionDetails, err := getConnectionDetails(fields)
		if err != nil {
			return nil, err
		}
		openConnectionsResult[i] = connectionDetails
	}
	return openConnectionsResult, nil
}

type socketId = string
type connectionDetails struct {
	SocketId       socketId
	LocalAddrIP    net.IP
	LocalAddrPort  string
	RemoteAddrIP   net.IP
	RemoteAddrPort string
}

func getConnectionDetails(connectionFields []string) (*connectionDetails, error) {
	// parse the /proc/{pid}/tcp line
	sid := connectionFields[9]
	localAddrIP := hexIpToDecimal(connectionFields[1][:8])
	localAddrPort := hexPortToDecimal(connectionFields[1][9:])
	remoteAddrIP := hexIpToDecimal(connectionFields[2][:8])
	remoteAddrPort := hexPortToDecimal(connectionFields[2][9:])
	return &connectionDetails{
		SocketId:       sid,
		LocalAddrIP:    net.ParseIP(localAddrIP),
		LocalAddrPort:  localAddrPort,
		RemoteAddrIP:   net.ParseIP(remoteAddrIP),
		RemoteAddrPort: remoteAddrPort,
	}, nil
}
