package stats

import (
	"github.com/blasrodri/frown/lsof"
	"github.com/blasrodri/frown/dns"
)

type ConnectionReport struct {
	DomainName       string
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
		DomainName:     domainName,
		SecurityLevel:  secLevel,
		AdditionalInfo: additionalInfo,
	}, err
}

func getSecurityLevel(connDeet *lsof.ConnectionDetails) (int, string) {
	return 0, ""
}
