package dns

import (
	"fmt"
	"net"
	"strings"
)

func getIpAddr(ipAddr *net.IP) ([]string, error) {
	return net.LookupAddr(ipAddr.String())
}

var KnownDomains = map[string]string{
	"compute.amazonaws.com": "Amazon AWS",
	"1e100.net":             "Google Cloud",
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
	for _, domain := range dnsDomains {
		for knownDom, niceName := range KnownDomains {
			if strings.HasSuffix(strings.TrimSuffix(domain, "."), knownDom) {
				return fmt.Sprintf("%s (%s)", knownDom, niceName), nil
			}
		}
	}
	return dnsDomains[0], nil
}
