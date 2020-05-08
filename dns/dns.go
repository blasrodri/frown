package dns

import (
	"net"
	"strings"
)

func getIpAddr(ipAddr *net.IP) ([]string, error) {
	return net.LookupAddr(ipAddr.String())
}

func DomainName(ipAddr *net.IP) (string, error) {
	domains, err := getIpAddr(ipAddr)
	if err != nil {
		// Dmoain not found, then just yield the ip as a string
		return ipAddr.String(), nil
	}
	knownDomain, err := getKnownDomain(domains)
	if err != nil {
		return domains[0], err
	}
	return knownDomain, nil
}

func getKnownDomain(dnsDomains []string) (string, error) {
	knownDomains := map[string]string{
		"compute.amazonaws.com": "Amazon AWS",
		"1e100.net":             "Google Cloud",
	}
	for _, domain := range dnsDomains {
		for knownDom, niceName := range knownDomains {
			if strings.HasSuffix(strings.TrimSuffix(domain, "."), knownDom) {
				return niceName, nil
			}
		}
	}
	return dnsDomains[0], nil
}
