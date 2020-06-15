package stats

import (
	"fmt"
	"net"
	"strconv"

	"github.com/blasrodri/frown/dns"
	"github.com/blasrodri/frown/lsof"
)

type ConnectionReport struct {
	SocketID       string
	DomainName     string
	SecurityLevel  int
	AdditionalInfo string
}

func AnalyzeSecurity(connDeet *lsof.ConnectionDetails) (*ConnectionReport, error) {
	domainName, err := dns.DomainName(&connDeet.RemoteAddrIP)
	if err != nil {
		return nil, err
	}
	secLevel, additionalInfo := getSecurityLevel(domainName, connDeet)
	return &ConnectionReport{
		SocketID:       connDeet.SocketID,
		DomainName:     domainName,
		SecurityLevel:  secLevel,
		AdditionalInfo: additionalInfo,
	}, err
}

func getSecurityLevel(domainName string, connDeet *lsof.ConnectionDetails) (int, string) {
	worseScore := 0
	message := ""
	portScore, msg := portHeuristic(connDeet)
	if portScore > worseScore {
		worseScore = portScore
		message = msg
	}

	hostScore, msg := hostHeuristic(domainName, connDeet)
	if hostScore > worseScore {
		worseScore = hostScore
		message = msg
	}
	return worseScore, message
}

func hostHeuristic(domainName string, connDeet *lsof.ConnectionDetails) (int, string) {
	if _, ok := dns.KnownDomains[domainName]; ok {
		return 0, ""
	}

	if connDeet.LocalAddrIP.Equal(net.ParseIP("0.0.0.0")) {
		return 2, fmt.Sprintf("Service is listening on all interfaces!")
	}

	if connDeet.RemoteAddrIP.Equal(net.ParseIP("0.0.0.0")) {
		return 2, fmt.Sprintf("Service is listening on all interfaces!")
	}

	return 1, fmt.Sprintf("Unknown host: %s", connDeet.RemoteAddrIP)
}

func portHeuristic(connDeet *lsof.ConnectionDetails) (int, string) {
	localPort, err := strconv.Atoi(connDeet.LocalAddrPort)
	remotePort, err := strconv.Atoi(connDeet.RemoteAddrPort)
	if err != nil {
		return 0, ""
	}
	unEncryptedRemotePorts := []int{80}
	safeRemotePorts := []int{443}
	databasePorts := []int{3306, 5432}
	sshPort := []int{22}

	if contains(sshPort, localPort) {
		return 3, "SSH Server open, and being accesed"
	}

	if contains(sshPort, remotePort) {
		return 3, "SSH Server open, and being accesed"
	}

	if contains(safeRemotePorts, remotePort) {
		return 0, ""
	}

	if contains(databasePorts, remotePort) {
		return 0, ""
	}

	if contains(databasePorts, localPort) {
		return 0, ""
	}

	if contains(unEncryptedRemotePorts, remotePort) {
		return 2, "Connection not encrypted"
	}

	if contains(sshPort, remotePort) {
		return 0, ""
	}

	return 1, fmt.Sprintf("Unknown ports: Remote %d, Local %d", remotePort, localPort)
}

func contains(arr []int, el int) bool {
	for _, value := range arr {
		if el == value {
			return true
		}
	}
	return false
}
