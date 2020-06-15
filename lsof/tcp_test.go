package lsof

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHexIpToDecimal(t *testing.T) {
	z := "0A010100"
	result := hexIpToDecimal(z)
	assert.Equal(t, "0.1.1.10", result)
}

func TestHexPortToDecimal(t *testing.T) {
	z := "01BB"
	result := hexPortToDecimal(z)
	assert.Equal(t, "443", result)
}
