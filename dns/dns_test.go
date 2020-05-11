package dns

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestgetKnownDomain(t *testing.T) {
	knownDomain, _ := getKnownDomain([]string{"ec2-52-38-207-226.us-west-2.compute.amazonaws.com"})
	assert.Equal(t, "Amazon AWS", knownDomain)
	knownDomain, _ = getKnownDomain([]string{"arn09s10-in-f4.1e100.net"})
	assert.Equal(t, "1e100.net (Google Cloud)", knownDomain)
}

func TestDomainName(t *testing.T) {
	ip := net.ParseIP("172.217.22.174")
	domName, _ := DomainName(&ip)
	assert.Equal(t, "1e100.net (Google Cloud)", domName)

}
