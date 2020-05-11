package stats

import (
	"github.com/blasrodri/frown/dns"
	"github.com/blasrodri/frown/lsof"
)

type ConnectionReport struct {
	SocketId       string
	DomainName     string
	SecurityLevel  int
	AdditionalInfo string
}

func AnalyzeSecurity(connDeet *lsof.ConnectionDetails) (*ConnectionReport, error) {
	secLevel, additionalInfo := getSecurityLevel(connDeet)
	domainName, err := dns.DomainName(&connDeet.RemoteAddrIP)
	if err != nil {
		return nil, err
	}
	return &ConnectionReport{
		SocketId:       connDeet.SocketId,
		DomainName:     domainName,
		SecurityLevel:  secLevel,
		AdditionalInfo: additionalInfo,
	}, err
}

func getSecurityLevel(connDeet *lsof.ConnectionDetails) (int, string) {
	return 0, ""
}
